package engine

import (
    "crawler/fetcher"
    "log"
)

// Worker 爬取request对应的URL的数据
// 对爬取的数据使用解析器进行解析
func Worker(r Request) (ParseResult, error) {
    log.Printf("fetching: %s", r.URL)
    contents, err := fetcher.Fetch(r.URL)
    if err != nil {
        log.Printf("Fetcher: error fetching url %s: %v", r.URL, err)
        return ParseResult{}, err
    }
    return r.Parse(r.URL, contents), nil
}
