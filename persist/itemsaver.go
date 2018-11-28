package persist

import (
    "context"
    "crawler/engine"
    "errors"

    "github.com/olivere/elastic"
)

// Save 用于将爬取的item存入elastic中
func Save(client *elastic.Client, item engine.Item, index string) error {
    if item.Type == "" {
        return errors.New("must supple Type")
    }

    indexService := client.Index().
        Index(index).
        Type(item.Type).
        BodyJson(item)

    if item.ID != "" {
        indexService.Id(item.ID)
    }

    _, err := indexService.Do(context.Background())
    if err != nil {
        return err
    }
    return nil
}
