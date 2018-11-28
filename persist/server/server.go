package main

import (
    "crawler/config"
    "crawler/consulsupport"
    "crawler/persist"
    "crawler/rpcsupport"
    "flag"
    "fmt"
    "log"

    "github.com/olivere/elastic"

    consulapi "github.com/hashicorp/consul/api"
)

var consulAddress = flag.String("consul_address", "localhost:8500", "consul address")

func main() {
    flag.Parse()

    // 创建consul client
    consulClient, err := consulsupport.NewClient(consulapi.DefaultNonPooledConfig(), *consulAddress)
    if err != nil {
        log.Fatal(err)
    }

    // 从consul获得elastic 对应的URL
    elasticURL, err := consulsupport.GetConfig(consulClient, "elastic/url")
    if err != nil {
        log.Fatal(err)
    }

    // 启动rpc server
    log.Fatal(serverRPC(fmt.Sprintf(":%d", config.ItemServerRPCPort), elasticURL, config.ElasticIndex))
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
