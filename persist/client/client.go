package client

import (
    "crawler/config"
    "crawler/engine"
    "crawler/rpcsupport"
    "log"
)

// ItemSaver 用于创建itemSaver Service的rpc连接, 并启动一个goroutine从
// itemChan接收数据, 当从itemChan接收到数据后, 通过rpc调用itemSaver Service
// 的save接口, 把数据存储到elastic
func ItemSaver(host string) (chan engine.Item, error) {
    client, err := rpcsupport.NewClient(host)
    if err != nil {
        return nil, err
    }

    out := make(chan engine.Item)
    itemCount := 0

    go func() {
        for {
            item := <-out
            itemCount++
            log.Printf("Item Saver: got item #%d, %v", itemCount, item)

            // Call RPC to save item
            result := false
            err := client.Call(config.ItemServerRPC, item, &result)
            if err != nil {
                log.Printf("Item Saver: error saving item %v %v", item, err)
            }
        }
    }()
    return out, nil
}
