package main

import (
    "crawler/config"
    "crawler/persist"
    "crawler/rpcsupport"
    "fmt"
    "log"
    "os"

    "github.com/olivere/elastic"
)

// var consulAddress = flag.String("consul_address", "localhost:8500", "consul address")

// func main() {
//     flag.Parse()

//     // 创建consul client
//     consulClient, err := consulsupport.NewClient(consulapi.DefaultNonPooledConfig(), *consulAddress)
//     if err != nil {
//         log.Fatal(err)
//     }

//     // 从consul获得elastic 对应的URL
//     elasticURL, err := consulsupport.GetConfig(consulClient, "elastic/url")
//     if err != nil {
//         log.Fatal(err)
//     }

//     // 启动rpc server
//     log.Fatal(serverRPC(fmt.Sprintf(":%d", config.ItemServerRPCPort), elasticURL, config.ElasticIndex))
// }

func main() {
    elasticServerName := os.Getenv("ELASTIC_SERVICE_NAME")
    if elasticServerName == "" {
        log.Fatal("must be specity elastic service name in environment variable.")
    }

    // 启动rpc server
    log.Fatal(serverRPC(fmt.Sprintf(":%d", config.ItemServerRPCPort),
        fmt.Sprintf("http://%s:9200/", elasticServerName),
        config.ElasticIndex))
}

func serverRPC(host, elasticURL, index string) error {
    client, err := elastic.NewClient(
        elastic.SetURL(elasticURL),
        elastic.SetSniff(false))
    if err != nil {
        return err
    }

    rpcsupport.ServerRPC(host, &persist.ItemSaverService{
        Client: client,
        Index:  index,
    })
    return nil
}
