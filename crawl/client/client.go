package client

import (
    "crawler/config"
    "crawler/consulsupport"
    "crawler/crawl"
    "crawler/engine"
    "crawler/rpcsupport"
    "errors"
    "fmt"
    "log"
    "net/rpc"
    "time"

    consulapi "github.com/hashicorp/consul/api"
)

type rpcClient struct {
    host      string
    client    *rpc.Client
    canBeUsed bool
}

var clientPool = make(map[string]*rpcClient)

// CreateProcessor 用于创建每个request对应的processor, processor把
// 每个URL的request进行序列化, 从clientChan获得crawl service,通过rpc
// 调用相应的爬虫接口, crawl service会对request的URL进行爬虫，并把相应
// 的结果返回给processor
func CreateProcessor(clientChan chan *rpcClient) engine.Processor {

    return func(r engine.Request) (engine.ParseResult, error) {
        var sResult crawl.ParseResult
        sRequest := crawl.SerializedRequest(r)

        client := <-clientChan

        if !clientPool[client.host].canBeUsed {
            return engine.ParseResult{}, fmt.Errorf("crawl service rpc can't be used: %s, request: %v", client.host, sRequest)
        }

        err := client.client.Call(config.CrawlServiceRPC, sRequest, &sResult)
        if err != nil {
            clientPool[client.host].canBeUsed = false // 标记client为不可用
            return engine.ParseResult{}, fmt.Errorf("Error: %v, crawl service rpc: %s, request: %v", err, client.host, sRequest)
        }

        return crawl.DeserializeParseResult(sResult), nil
    }
}

func addClientToPool(hosts []string) {
    for _, host := range hosts {
        client, err := rpcsupport.NewClient(host)
        if err != nil {
            log.Printf("Error connecting to: %s", host)
        } else {
            log.Printf("connected to: %s", host)
            var c rpcClient
            c.host = host
            c.client = client
            c.canBeUsed = true
            clientPool[host] = &c
        }
    }
}

func removeClientFromPool(host string) {
    // 关闭客户端连接, 并从连接池中删除
    clientPool[host].client.Close()
    delete(clientPool, host)
}

func differentHosts(hosts []string) ([]string, map[string]struct{}) {
    // 获得要添加或删除的host
    var newHosts []string
    deleteHosts := make(map[string]struct{}, len(clientPool))

    for host := range clientPool {
        deleteHosts[host] = struct{}{}
    }

    for _, host := range hosts {
        if _, ok := clientPool[host]; !ok {
            newHosts = append(newHosts, host)
        } else {
            delete(deleteHosts, host)
        }
    }
    return newHosts, deleteHosts
}

func updateClientPool(hosts []string) {
    newHosts, deleteHosts := differentHosts(hosts)

    // add new host to client pool
    if len(newHosts) != 0 {
        log.Printf("add host: %v\n", newHosts)
        addClientToPool(newHosts)
    }

    // delete host
    for host := range deleteHosts {
        log.Printf("delete host: %v\n", host)
        removeClientFromPool(host)
    }
}

// CreateClientPool 用于创建crawlServiceName对应服务的连接池
func CreateClientPool(consulClient *consulapi.Client, serviceName string) (chan *rpcClient, error) {
    // 从consul中获得crawlServiceName对应服务的 host:port
    hosts, err := consulsupport.DiscoveryService(consulClient, serviceName)
    if err != nil {
        return nil, err
    }

    // 创建连接池
    addClientToPool(hosts)
    if len(clientPool) == 0 {
        return nil, errors.New("There are no clients available")
    }

    clientChan := make(chan *rpcClient)

    go func() {
        t := time.NewTicker(time.Second * config.UpdateClientPoolSecond)

        for {
            select {
            case <-t.C:
                // 定时更新连接池
                if hosts, err := consulsupport.DiscoveryService(consulClient, serviceName); err != nil {
                    log.Println(err)
                } else {
                    updateClientPool(hosts)
                }
            default:
                // 轮询连接池中的每个crawl service
                // 如果crawl service可以使用，则通过clientChan发送给请求者
                for _, client := range clientPool {
                    if client.canBeUsed {
                        clientChan <- client
                    }
                }
            }
        }
    }()
    return clientChan, nil
}
