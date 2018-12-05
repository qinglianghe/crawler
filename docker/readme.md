# 微服务化部署

以下通过Docker Swarm对项目进行微服务部署。

## 搭建 Swarm 集群

搭建如下Swarm集群，包括三个Node，对应的主机名和IP地址分别如下所示，这三台机器是可以互相通信的，`192.168.33.10`节点为manager节点、其他两个节点为worker节点：

| Node | IP |
|-|:-|
| manager | 192.168.33.10 |
| worker1 | 192.168.33.11 |
| worker2 | 192.168.33.12 |

### 启动 manager 节点

    ubuntu@swarm-manager:~/go/src/crawler/docker$ docker swarm init --advertise-addr=192.168.33.10

### 启动 worker 节点

    ubuntu@swarm-worker1:~/go/src/crawler/docker$ docker swarm join \
         --token SWMTKN-1-07hop0ri280sugqz79zoffh8wmkbf2sxey3k9s3d1ywhaexp6a-by2h0yygptlm0wfkdgoohs7e6 \
         192.168.33.10:2377


    ubuntu@swarm-worker2:~$ docker swarm join \
         --token SWMTKN-1-07hop0ri280sugqz79zoffh8wmkbf2sxey3k9s3d1ywhaexp6a-by2h0yygptlm0wfkdgoohs7e6 \
         192.168.33.10:2377

### 查看节点状态

    ubuntu@swarm-manager:~/go/src/crawler/docker$ docker node ls
    ID                           HOSTNAME       STATUS  AVAILABILITY  MANAGER STATUS
    rkur719xhnk0zi1b700rjgw8p    swarm-worker1  Ready   Active
    xal3c9dh09pogk7dzoyxr2xtg    swarm-worker2  Ready   Active
    y271lmw7frng7p6btfwic0pfs *  swarm-manager  Ready   Active        Leader

## 生成 Dcoker 镜像

### engine 镜像

进入到`crawler/docker`目录，编译生成 engine service：

    go build -o engine main.go

使用本目录下的[Dockerfile](https://github.com/qinglianghe/crawler/blob/master/docker/Dockerfile)，编译生成 engine image:

    docker build -t heqingliang/engine .

### itemsaver 镜像

进入到`crawler/docker/persist/server`目录，编译生成 itemsaver service：

    go build -o itemsaver server.go

使用本目录下的[Dockerfile](https://github.com/qinglianghe/crawler/blob/master/docker/persist/server/Dockerfile)，编译生成 itemsaver image:

    docker build -t heqingliang/itemsaver .

### crawl 镜像

进入到`crawler/docker/crawl/server/`目录，编译生成 crawl service：

    go build -o crawl server.go

使用本目录下的[Dockerfile](https://github.com/qinglianghe/crawler/blob/master/docker/crawl/server/Dockerfile)，编译生成 crawl image:

    docker build -t heqingliang/crawl .

需要把上面生成的镜像push到自己的Docker Hub，因为部署的每个镜像部署到哪个节点是不定的。

### 部署

在manage节点，使用`crawler/docker/`目录下的[dcoker-compose.yml](https://github.com/qinglianghe/crawler/blob/master/docker/dcoker-compose.yml)文件，用`docker satck`进行部署，`dcoker-compose.yml`定义了项目中需要启动的service、每个service需要的image、重启策略、启动的顺序等：

    ubuntu@swarm-manager:~/go/src/crawler/docker$ docker stack deploy --compose-file docker-compose.yml crawler
    Creating network crawler_overlay_network
    Creating service crawler_itemsaver
    Creating service crawler_crawl
    Creating service crawler_engine
    Creating service crawler_elastic
    Creating service crawler_redis

查看服务的状态：

    ubuntu@swarm-manager:~/go/src/crawler/docker$ docker service ls
    ID            NAME               MODE        REPLICAS  IMAGE
    0mx9nau0l052  crawler_elastic    global      1/1       elasticsearch:6.4.1
    95vldgo1l115  crawler_redis      global      1/1       redis
    ort8ukftix0u  crawler_crawl      replicated  3/3       heqingliang/crawl
    p9myrfzi6g3x  crawler_engine     global      3/3       heqingliang/engine
    zqxknuufbc20  crawler_itemsaver  global      3/3       heqingliang/itemsaver

查看某个服务(crawl)的状态：

    ubuntu@swarm-manager:~/go/src/crawler/docker$ docker service ps crawler_crawl
    ID            NAME             IMAGE              NODE           DESIRED STATE  CURRENT STATE          ERROR  PORTS
    czpr9pjacx7j  crawler_crawl.1  heqingliang/crawl  swarm-manager  Running        Running 2 minutes ago
    jeqwjb0bja3d  crawler_crawl.2  heqingliang/crawl  swarm-worker1  Running        Running 2 minutes ago
    sjgsvyncrc5o  crawler_crawl.3  heqingliang/crawl  swarm-worker2  Running        Running 3 minutes ago