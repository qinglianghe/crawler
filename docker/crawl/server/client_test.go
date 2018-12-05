package main

import (
    "crawler/config"
    "crawler/crawl"
    "crawler/rpcsupport"
    "fmt"
    "testing"
    "time"
)

func TestCrawlService(t *testing.T) {
    host := fmt.Sprintf(":%d", config.CrawlServiceRPCPort)

    go rpcsupport.ServerRPC(host, &crawl.CrawlService{})

    time.Sleep(1 * time.Second)

    client, err := rpcsupport.NewClient(host)
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
    client.Close()
}
