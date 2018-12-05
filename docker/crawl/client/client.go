package client

import (
    "crawler/config"
    "crawler/crawl"
    "crawler/engine"
    "crawler/rpcsupport"
    "fmt"
)

// CreateProcessor 用于创建每个request对应的processor, processor把
// 每个URL的request进行序列化, 从clientChan获得crawl service,通过rpc
// 调用相应的爬虫接口, crawl service会对request的URL进行爬虫，并把相应
// 的结果返回给processor
func CreateProcessor(crawlServiceName string) engine.Processor {

    return func(r engine.Request) (engine.ParseResult, error) {
        var sResult crawl.ParseResult
        sRequest := crawl.SerializedRequest(r)

        // RPC client
        client, err := rpcsupport.NewClient(fmt.Sprintf("%s:%d", crawlServiceName, config.CrawlServiceRPCPort))
        if err != nil {
            return engine.ParseResult{}, fmt.Errorf("Error connecting to crawl server: %v", err)
        }

        err = client.Call(config.CrawlServiceRPC, sRequest, &sResult)
        if err != nil {
            return engine.ParseResult{}, fmt.Errorf("Error: %v, request: %v", err, sRequest)
        }

        client.Close()
        return crawl.DeserializeParseResult(sResult), nil
    }
}
