package main

import (
    "crawler/config"
    "crawler/consulsupport"
    "crawler/engine"
    "crawler/model"
    "crawler/rpcsupport"
    "testing"
    "time"

    consulapi "github.com/hashicorp/consul/api"
)

func TestItemSaver(t *testing.T) {
    consulClient, err := consulsupport.NewClient(consulapi.DefaultNonPooledConfig(), "localhost:8500")
    if err != nil {
        panic(err)
    }

    elasticURL, err := consulsupport.GetConfig(consulClient, "elastic/url")
    if err != nil {
        panic(err)
    }

    hosts, err := consulsupport.DiscoveryService(consulClient, "itemsaver")
    if err != nil {
        panic(err)
    }

    if len(hosts) == 0 {
        panic("There are no service available")
    }

    // start ItemSaverServer
    go serverRPC(hosts[0], elasticURL, "test1")

    time.Sleep(1 * time.Second)

    // start ItemSaverClient
    client, err := rpcsupport.NewClient(hosts[0])
    if err != nil {
        panic(err)
    }

    result := false

    item := engine.Item{
        URL:  "http://album.zhenai.com/u/1439637023",
        Type: "zhenai",
        ID:   "1439637023",
        Payload: model.Profile{
            Name:        "一切随缘",
            Gender:      "女",
            Age:         23,
            Height:      160,
            Weight:      48,
            Income:      "8001-12000元",
            Marriage:    "未婚",
            Education:   "中专",
            Occupation:  "其他职业",
            WorkPlace:   "广东广州",
            Xinzuo:      "狮子座",
            House:       "打算婚后购房",
            Car:         "未购车",
        },
    }

    // Call Save
    err = client.Call(config.ItemServerRPC, item, &result)
    if err != nil || result != true {
        t.Errorf("result: %v; err: %s", result, err)
    }
}
