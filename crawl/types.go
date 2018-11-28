package crawl

import (
    "crawler/config"
    "crawler/engine"
    "crawler/zhenai/parser"
    "errors"
    "fmt"
    "log"
)

// SerializedParser 用于序列化解析器
type SerializedParser struct {
    // Name 解析器的名字
    Name string

    // 解析器所需要的参数
    Args interface{}
}

// Request 用于序列化requst
type Request struct {
    // 爬虫页面对应的URL
    URL string

    // 解析器序列化接口
    Parser SerializedParser
}

// ParseResult 用于序列化爬虫所得结果
type ParseResult struct {
    Requests []Request
    Items    []engine.Item
}

// SerializedRequest 将engine.Request序列化为Request
func SerializedRequest(r engine.Request) Request {
    name, args := r.Serialize()
    return Request{
        URL: r.URL,
        Parser: SerializedParser{
            Name: name,
            Args: args,
        },
    }
}

// SerializedParseResult 将engine.ParseResult序列化为ParseResult
func SerializedParseResult(p engine.ParseResult) ParseResult {
    result := ParseResult{
        Items: p.Items,
    }
    for _, r := range p.Requests {
        result.Requests = append(result.Requests, SerializedRequest(r))
    }
    return result
}

// DeserializeRequest 将Request反序列化为engine.Request
func DeserializeRequest(r Request) (engine.Request, error) {
    parser, err := deserializeParser(r.Parser)
    if err != nil {
        return engine.Request{}, err
    }

    request := engine.Request{
        URL:    r.URL,
        Parser: parser,
    }

    return request, nil
}

func deserializeParser(p SerializedParser) (engine.Parser, error) {
    switch p.Name {
    case config.ParseCityList:
        return engine.NewParser(config.ParseCityList, parser.ParseCityList), nil
    case config.ParseCity:
        return engine.NewParser(config.ParseCity, parser.ParseCity), nil
    case config.ProfileParser:
        if userName, ok := p.Args.(string); ok {
            return parser.NewProfileParser(userName), nil
        } else {
            return nil, fmt.Errorf("invaild arg: %v", p.Args)
        }

    case config.NilParser:
        return &engine.NilParser{}, nil
    default:
        return nil, errors.New("unknow parser name")
    }
}

// DeserializeParseResult 将ParseResult反序列化为engine.ParseResult
func DeserializeParseResult(p ParseResult) engine.ParseResult {
    result := engine.ParseResult{
        Items: p.Items,
    }
    for _, r := range p.Requests {
        engineRequest, err := DeserializeRequest(r)
        if err != nil {
            log.Printf("error deserializing request: %v", engineRequest)
            continue
        }
        result.Requests = append(result.Requests, engineRequest)
    }
    return result
}
