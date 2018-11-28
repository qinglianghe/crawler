package crawl

import "crawler/engine"

// CrawlService 爬虫模块
type CrawlService struct{}

// Process 爬虫的处理函数
// 反序列化request
// 爬取相应页面的数据
// 序列化爬取的结果
func (s *CrawlService) Process(request Request, result *ParseResult) error {
    engineRequest, err := DeserializeRequest(request)
    if err != nil {
        return err
    }

    engineResult, err := engine.Worker(engineRequest)
    if err != nil {
        return err
    }

    *result = SerializedParseResult(engineResult)
    return nil
}
