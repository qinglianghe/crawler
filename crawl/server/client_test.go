package main

import (
    "crawler/config"
    "crawler/consulsupport"
    "crawler/crawl"
    "crawler/rpcsupport"
    "fmt"
    "testing"
    "time"

    consulapi "github.com/hashicorp/consul/api"
)

func TestCrawlService(t *testing.T) {
    consulClient, err := consulsupport.NewClient(consulapi.DefaultNonPooledConfig(), "localhost:8500")
    if err != nil {
        panic(err)
    }

    hosts, err := consulsupport.DiscoveryService(consulClient, "crawl")
    if err != nil {
        panic(err)
    }

    if len(hosts) == 0 {
        panic("There are no service available")
    }

    go rpcsupport.ServerRPC(hosts[0], &crawl.CrawlService{})

    time.Sleep(1 * time.Second)

    client, err := rpcsupport.NewClient(hosts[0])
    if err != nil {
        panic(err)
    }

    request := crawl.Request{
        URL: "http://album.zhenai.com/u/1439637023",
        Parser: crawl.SerializedParser{
            Name: config.ProfileParser,
            Args: "一切随缘",
        },
    }

    var result crawl.ParseResult

    err = client.Call(config.CrawlServiceRPC, request, &result)
    if err != nil {
        t.Errorf("result: %v; err: %s", result, err)
    } else {
        fmt.Println(result)
    }
}
