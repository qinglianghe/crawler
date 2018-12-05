package parser

import (
    "crawler/engine"
    "regexp"
)

var (
    profileRe = regexp.MustCompile(`<th><a href="(http://album.zhenai.com/u/[0-9]+)" target="_blank">([^<]+)</a></th>`)
    cityURLRe = regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
)

// ParseCity 用于解析从ParseCityList获得城市的信息
// 这个页面中包含用户信息，和城市列表信息
// 通过profileRe正则表达式获得用户信息，对应的解析器为ProfileParser
// 通过cityURLRe正则表达式获得用户信息，对应的解析器为ParseCity
func ParseCity(_ string, contents []byte) engine.ParseResult {
    matches := profileRe.FindAllSubmatch(contents, -1)

    result := engine.ParseResult{}

    for _, m := range matches {
        result.Requests = append(
            result.Requests, engine.Request{
                URL:    string(m[1]),
                Parser: NewProfileParser(string(m[2])),
            })
    }

    matches = cityURLRe.FindAllSubmatch(contents, -1)

    for _, m := range matches {
        result.Requests = append(
            result.Requests, engine.Request{
                URL:    string(m[1]),
                Parser: engine.NewParser("ParseCity", ParseCity),
            })
    }
    return result
}
