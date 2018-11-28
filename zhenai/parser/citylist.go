package parser

import (
	"crawler/engine"
	"regexp"
)

var cityListRe = regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)

// ParseCityList 通过cityListRe正则表达式解析城市列表
// 每个城市对应的解析器为ParseCity
func ParseCityList(url string, contents []byte) engine.ParseResult {
	matches := cityListRe.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	for _, m := range matches {
		result.Requests = append(
			result.Requests, engine.Request{
				URL:    string(m[1]),
				Parser: engine.NewParser("ParseCity", ParseCity),
			})
	}
	return result
}
