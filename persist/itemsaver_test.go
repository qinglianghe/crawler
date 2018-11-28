package persist

import (
    "context"
    "crawler/engine"
    "crawler/model"
    "encoding/json"
    "testing"

    "github.com/olivere/elastic"
)

func TestSave(t *testing.T) {
    expected := engine.Item{
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
            NativePlace: "广东广州",
            Xinzuo:      "狮子座",
            House:       "打算婚后购房",
            Car:         "未购车",
        },
    }

    // TODO: try to start up elastic search
    // here using docker go client
    client, err := elastic.NewClient(
        //        elastic.SetURL("http://192.168.1.135:9200/"),
        elastic.SetSniff(false))
    if err != nil {
        panic(err)
    }

    // save expected item
    const index = "dating_test"
    err = Save(client, expected, index)
    if err != nil {
        panic(err)
    }

    resp, err := client.Get().
        Index(index).
        Type(expected.Type).
        Id(expected.ID).
        Do(context.Background())
    if err != nil {
        panic(err)
    }

    t.Logf("%s", resp.Source)

    var actual engine.Item
    err = json.Unmarshal(*resp.Source, &actual)
    if err != nil {
        panic(err)
    }

    actualProfile, err := model.FromJSONObj(actual.Payload)
    if err != nil {
        panic(err)
    }
    actual.Payload = actualProfile

    if actual != expected {
        t.Errorf("got %v; expected %v", actual, expected)
    }
}
