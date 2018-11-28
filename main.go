package main

import (
    "crawler/config"
    "crawler/consulsupport"
    crawl "crawler/crawl/client"
    "crawler/engine"
    itemsaver "crawler/persist/client"
    "crawler/scheduler"
    "crawler/zhenai/parser"
    "flag"
    "log"

    consulapi "github.com/hashicorp/consul/api"
)

var (
    //    crawlHosts           = flag.String("crawl_hosts", "", "crawl host(comma split)")
    consulAddress        = flag.String("consul_address", "localhost:8500", "consul address")
    crawlServiceName     = flag.String("crawl_service_name", "", "crawl service name")
    itemSaverServiceName = flag.String("item_saver_service_name", "", "itemsaver service name")
)

// usage:
//  go run main.go -crawl_service_name=crawl -item_saver_service_name=itemsaver -consul_address=172.17.0.2:8500
//  go run main.go -crawl_service_name=crawl -item_saver_service_name=itemsaver
func main() {
    flag.Parse()

    if *itemSaverServiceName == "" {
        log.Fatal("must specity itemsaver service name")
    }

    if *crawlServiceName == "" {
        log.Fatal("must specity crawler service name")
    }

    // 创建consul client
    consulClient, err := consulsupport.NewClient(consulapi.DefaultConfig(), *consulAddress)
    if err != nil {
        log.Fatal(err)
    }

    // 从consul中获得itemSaverServiceName对应服务的 host:port
    itemSaverhost, err := getHostWithServiceName(consulClient, *itemSaverServiceName)
    if err != nil {
        log.Fatal(err)
    }

    // 创建itemSaver Service的client
    itemChan, err := itemsaver.ItemSaver(itemSaverhost)
    if err != nil {
        log.Fatal(err)
    }

    // pool := createClientPool(strings.Split(*crawlHosts, ","))

    // 创建连接池
    clientChan, err := crawl.CreateClientPool(consulClient, *crawlServiceName)
    if err != nil {
        log.Fatal(err)
    }

    // 从consul中获得redis的 host:port
    redisHost, err := consulsupport.GetConfig(consulClient, "redis/address")
    if err != nil {
        log.Fatal(err)
    }

    // 创建每个URL的request对应的processor
    processor := crawl.CreateProcessor(clientChan)

    e := engine.ConcurrentEngine{
        Scheduler:        &scheduler.QueuedScheduler{},
        WorkerCount:      100,
        ItemChan:         itemChan,
        RequestProcessor: processor,
        RedisHost:        redisHost,
    }

    // 启动爬虫，从起始页面 http://www.zhenai.com/zhenghun 开始爬取
    e.Run(engine.Request{
        URL:    "http://www.zhenai.com/zhenghun",
        Parser: engine.NewParser(config.ParseCityList, parser.ParseCityList),
    })
}

func getHostWithServiceName(consulClient *consulapi.Client, serviceName string) (string, error) {
    hosts, err := consulsupport.DiscoveryService(consulClient, serviceName)
    if err != nil {
        return "", err
    }

    if len(hosts) > 1 {
        log.Printf("There to many %s service, choice %s to connect\n", serviceName, hosts[0])
    }
    return hosts[0], nil
}

// func createClientPool(hosts []string) chan *rpc.Client {
//     var clients []*rpc.Client

//     for _, host := range hosts {
//         client, err := rpcsupport.NewClient(host)
//         if err != nil {
//             log.Printf("Error connecting to: %s", host)
//         } else {
//             clients = append(clients, client)
//             log.Printf("Connected to: %s", host)
//         }
//     }

//     clientChan := make(chan *rpc.Client)

//     go func() {
//         for {
//             for _, client := range clients {
//                 clientChan <- client
//             }
//         }
//     }()

//     return clientChan
// }
