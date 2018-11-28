# 分布式爬虫

这是由go编写的分布式爬虫的项目，用于爬取真爱网的用户信息。

## 爬虫的过程

1. 通过起始页面[http://www.zhenai.com/zhenghun](http://www.zhenai.com/zhenghun)可以获取到城市列表的链接。

    ![zhenghun.png](/images/zhenghun.png)

2. 通过城市列表页面可以获取到用户列表的链接和这个城市区域的链接，例如：在上面的城市列表中，点击广州的连接[http://www.zhenai.com/zhenghun/guangzhou](http://www.zhenai.com/zhenghun/guangzhou)：

    ![guangzhou.png](/images/guangzhou.png)

3. 查看用户页面，可以得到用户信息：

    ![profile.png](/images/profile.png)

对应获得的结构化数据如下：

```go
    type Profile struct {
        // Name 名字
        Name string

        // Gender 性别
        Gender string

        // Age 年龄
        Age int

        // Height 身高
        Height int

        // Weight 体重
        Weight int

        // Income 收入
        Income string

        // Marriage 婚姻状况
        Marriage string

        // Education 学历
        Education string

        // Occupation 职业
        Occupation string

        // WorkPlace 工作地
        WorkPlace string

        // Xinzuo 星座
        Xinzuo string

        // 是否已购房
        House string

        // 是否已购车
        Car string
    }
```

## 架构设计

如下图是程序对应的架构：

![crawler_arch.jpg](/images/crawler_arch.jpg)

### engine

- 通过`createWorker`函数，创建指定数量的worker(默认为100个)

- 通过`Submit`函数，把Request提交给Scheduler

- 从channel接收worker返回的解析结果

- 如果爬取到了Item，则启动goroutine发送到存储模块的channel

### Scheduler

- 通过channel接收从engine发送过来的Request，并放入Request队列

- 通过channel接收从worker发送过来的空闲的worker，并放入worker队列

- 当Request队列有Request，并且worker队列有空闲的worker时，则把Request发送给空闲的worker

### worker

- 通过Redis对URL进行去重

- 序列化Request

- 从连接池中获得已连接crawl service的client，通过RPC发送给crawl service

- 通过RPC接收crawl service返回的结果，反序列化后，通过channel发送给engine

### crawl service 连接池

- 通过服务的名字从consul获得所有crawl service的host:port，并建立连接，放入连接池

- 启动一个goroutine轮询连接池，实现负载均衡

- goroutine会每隔一定的时间，从consul发现服务，更新连接池。

### crawl service

- 由`crawler/crawl/server/consul.d/crawl_service.json`的配置文件，注册服务`crawl`到Consul集群中

- 可以注册多个多个服务，但是IP地址必须不同，端口号必须和config配置的`CrawlServiceRPCPort`一样(默认为9987)

- crawl service启动后，会启动一个RPC服务`:9987`，用于接收Request请求

- 当crawl service接收到Request请求后，首先会反序列化Request，通过http请求，获得指定URL的页面，并通过Request指定的页面解析器对页面数据进行解析，并将结果序列后，通过RPC发送给客户端

- 爬虫的速度由config配置的`QPS`进行控制

### Consul

- 提供服务注册和服务发现功能

- 使用kv存储redis的地址和elasticsearch的URL

### itemsaver client

- 启动goroutine从channel中接收由engine发送的Item

- 通过RPC把Item发送给itemsaver service

### itemsaver service

- 由`crawler/persist/server/consul.d/itemsaver_service.json`的配置文件，注册服务`itemsaver`到Consul集群中

- itemsaver service中并没有多个服务创建连接池，如果有多个服务，只使用一个服务。默认使用的是从`getHostWithServiceName`返回的地址

- itemsaver_service.json配置的端口号必须和config配置的`ItemServerRPCPort`一样(默认为1234)

- itemsaver service启动时，会从Consul中获取elasticsearch的URL，并与elasticsearch建立连接

- 当crawl service接收到客户端发送过来的Item时，会调用`Save`函数存储到elasticsearch中

### 业务逻辑 (解析器 Parse)

在`crawler/zhenai`目录下，有相应页面的解析代码。每个页面都有一个对应的解析器，用于解析对应的页面的数据。如：

1. [http://www.zhenai.com/zhenghun](http://www.zhenai.com/zhenghun)：起始页面对应解析器为：ParseCityList，用于获取城市列表的URL，每个城市的URL对应的解析器为ParseCity，对应的单元测试函数为TestParseCityList。

2. [http://www.zhenai.com/zhenghun/guangzhou](http://www.zhenai.com/zhenghun/guangzhou)：城市页面对应的解析器为：ParseCity，用于获取城市区域的URL，每个城市区域的URL对应的解析器为ParseCity和用户信息的URL对应的解析器为ProfileParser，对应的单元测试函数为TestParseCity。

3. 用户信息的解析器为：ProfileParser，用于获取用户信息，对应的单元测试函数为TestParseProfile。

### Parse 序列化和反序列化

在客户端和crawl service通信时，使用的是RPC。所以要对Parse进行序列化和反序列化。在`crawler/crawl/server/types.go`定义有序列化和反序列化接口函数：`SerializedRequest`、`SerializedParseResult`、`DeserializeRequest`、`DeserializeParseResult`。

### 运行

需要在本机中安装Docker和Consul。

#### 安装依赖库

    go get -v -u github.com/gpmgo/gopm
    gopm get -g -v golang.org/x/text
    gopm get -g -v golang.org/x/net/html
    go get -v github.com/hashicorp/consul
    go get -v github.com/garyburd/redigo/redis
    go get -v github.com/olivere/elastic

#### 启动 Consul集群

下面使用启动Docker启动4个Consul集群：

    # 启动第1个Server节点，集群要求要有3个Server，将容器8500端口映射到主机8900端口，同时开启管理界面
    docker run -d --name=consul1 -p 8900:8500 -e CONSUL_BIND_INTERFACE=eth0 consul agent --server=true --bootstrap-expect=3 --client=0.0.0.0 -ui

    # 启动第2个Server节点，并加入集群
    docker run -d --name=consul2 -e CONSUL_BIND_INTERFACE=eth0 consul agent --server=true --client=0.0.0.0 --join 172.17.0.2

    # 启动第3个Server节点，并加入集群
    docker run -d --name=consul3 -e CONSUL_BIND_INTERFACE=eth0 consul agent --server=true --client=0.0.0.0 --join 172.17.0.2

    # 启动第4个Client节点，并加入集群
    docker run -d --name=consul4 -e CONSUL_BIND_INTERFACE=eth0 consul agent --server=false --client=0.0.0.0 --join 172.17.0.2

因本机和Consul容器是通过bridge网络连接的，所以本机和各个Consul容器间是可以通信的，在本机浏览器中访问`http://127.0.0.1:8900/ui/dc1/nodes`，可以看到Consul集群是否启动成功：

![consul_nodes.png](/images/consul_nodes.png)

#### 启动 elasticsearch

    docker run -p 9200:9200 -d --name=elastic elasticsearch:6.4.1

在本机浏览器中访问`http://127.0.0.1:9200/`，查看elasticsearch是否启动成功：

![elasticsearch.png](/images/elasticsearch.png)

#### 启动 Redis

    docker run -d -p 6379:6379 --name=redis redis

#### 设置 elasticsearch 和 Reids 的配置

    docker exec -it consul1 sh

    # 192.168.1.138为本机的IP
    consul kv put elastic/url http://192.168.1.138:9200/
    consul kv put redis/address 192.168.1.138:6379

#### 启动 itemsaver service

    cd crawler/persist/server

    # 编译 itemsaver service
    CGO_ENABLED=0 go build -o itemsaver -a -installsuffix cgo 

    docker cp itemsaver consul2:/
    docker cp consul.d/itemsaver_service.json consul2:/consul/config/

    # 重新加载consul配置，注册 itemsaver service
    docker exec -it consul2 sh
    consul reload

    # 启动 itemsaver service
    ./itemsaver

在浏览器中输入`http://127.0.0.1:8900/ui/dc1/services/itemsaver`，可以看到itemsaver service 启动成功：

![itemsaver.png](/images/itemsaver.png)

#### 启动 crawl service

    cd crawler/crawl/server

    # 编译 crawl service
    CGO_ENABLED=0 go build -o crawl -a -installsuffix cgo 

    docker cp crawl consul3:/
    docker cp consul.d/crawl_service.json consul3:/consul/config/

    # 重新加载consul配置，注册 crawl service
    docker exec -it consul3 sh
    consul reload

    # 启动 crawl service
    ./crawl

    # 在本机中再启动一个 crawl service，需要更改crawl_service.json的IP地址为本机的IP
    consul agent --server=false -config-dir=crawler/crawl/server/consul.d/ -data-dir=/tmp --bind=192.168.1.138 --join 172.17.0.2 &
    go run crawl/server/server.go

在浏览器中输入`http://127.0.0.1:8900/ui/dc1/services/crawl`，可以看到启动了 crawl service：

![crawl.png](/images/crawl.png)

#### 启动主程序

    go run crawler/main.go -crawl_service_name=crawl -item_saver_service_name=itemsaver -consul_address=localhost:8500

### 责任声明

如果使用次项目，触犯了任何商业利益，本人不承担任何责任。本人承诺，此项目本人未使用于任何的商业用途，也为将数据分享过任何人。如果此项目触犯了贵公司的相关的利益，本人愿意立即删除此项目。