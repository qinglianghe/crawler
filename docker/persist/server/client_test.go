package main

import (
    "crawler/config"
    "crawler/engine"
    "crawler/model"
    "crawler/rpcsupport"
    "fmt"
    "testing"
    "time"
)

func TestItemSaver(t *testing.T) {
    host := fmt.Sprintf(":%d", config.ItemServerRPCPort)

    // start ItemSaverServer
    go serverRPC(host, "http://localhost:9200", "test1")

    time.Sleep(1 * time.Second)

    // start ItemSaverClient
    client, err := rpcsupport.NewClient(host)
    if err != nil {
        panic(err)
    }

    result := false

    item := engine.Item{
        URL:  "http://album.zhenai.com/u/1439637023",
        Type: "zhenai",
        ID:   "1439637023",
        Payload: model.Profile{
            Name:       "一切随缘",
            Gender:     "女",
            Age:        23,
            Height:     160,
            Weight:     48,
            Income:     "8001-12000元",
            Marriage:   "未婚",
            Education:  "中专",
            Occupation: "其他职业",
            WorkPlace:  "广东广州",
            Xinzuo:     "狮子座",
            House:      "打算婚后购房",
            Car:        "未购车",
        },
    }

    // Call Save
    err = client.Call(config.ItemServerRPC, item, &result)
    if err != nil || result != true {
        t.Errorf("result: %v; err: %s", result, err)
    }
}
