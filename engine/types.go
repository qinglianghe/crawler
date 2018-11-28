package engine

// ParseFunc 解析器对应的函数
type ParseFunc func(string, []byte) ParseResult

// Parser 解析器接口
type Parser interface {
    // Parse 解析器
    Parse(string, []byte) ParseResult

    // Serialize rpc传输时的序列化接口
    Serialize() (string, interface{})
}

// Request 每个URL请求对应一个Request
// 请求所得的结果, 通过调用Parser进行解析
type Request struct {
    // URL 请求的URL
    URL string

    // Parser 请求所得的结果的解析器
    Parser
}

// ParseResult 解析器解析完的结果
type ParseResult struct {
    // 每个URL对应的页面, 解析完后可能有多个结果
    Requests []Request

    // Items 解析器解析完对应的结构化数据
    Items []Item
}

// Item 解析器解析完对应的结构化数据
type Item struct {
    // 请求的URL
    URL string

    // 存储到elastic对应的type
    Type string

    // Id 爬取的用户id
    ID string

    // Payload 用户的信息
    Payload interface{}
}

// NilParser nil的解析器
type NilParser struct{}

// Parse NilParser对应的解析器
func (p *NilParser) Parse(_ string, _ []byte) ParseResult {
    return ParseResult{}
}

// Serialize NilParser对应的序列化接口
func (p *NilParser) Serialize() (name string, args interface{}) {
    return "NilParser", nil
}

// FuncParser 解析器转换为对应的解析函数
type FuncParser struct {
    // name 函数名
    name string

    // parser 解析器对应的解析函数
    parser ParseFunc
}

// Parse 调用ParseFunc中的parser解析函数对数据进行解析
func (p *FuncParser) Parse(url string, contents []byte) ParseResult {
    return p.parser(url, contents)
}

// Serialize 解析器序列化的结果
func (p *FuncParser) Serialize() (name string, args interface{}) {
    return p.name, nil
}

// NewParser 创建一个新的FuncParser
func NewParser(name string, parser ParseFunc) *FuncParser {
    return &FuncParser{
        name:   name,
        parser: parser,
    }
}
