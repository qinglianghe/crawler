package main

import (
    "crawler/config"
    "crawler/crawl"
    "crawler/rpcsupport"
    "fmt"
    "log"
)

func main() {
    // 启动crawl service
    log.Fatal(rpcsupport.ServerRPC(
        fmt.Sprintf(":%d", config.CrawlServiceRPCPort),
        &crawl.CrawlService{}))
}
