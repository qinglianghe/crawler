package parser

import (
    "crawler/engine"
    "crawler/model"
    "regexp"
    "strconv"
)

var (
    // ageRe         = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)岁</td>`)
    // heightRe      = regexp.MustCompile(`<td><span class="label">身高：</span>(\d+)CM</td>`)
    // incomeRe      = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
    // weightRe      = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">(\d+)KG</span></td>`)
    // genderRe      = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
    // xinzuoRe      = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
    // marriageRe    = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
    // educationRe   = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
    // occupationRe  = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
    // nativePlaceRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
    // houseRe       = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
    // carRe         = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
    // idRe          = regexp.MustCompile(`http://album.zhenai.com/u/(\d+)`)

    personInfoRe = regexp.MustCompile(`<div class="m-btn purple[^>]*>([^<]+)</div>`)
    ageRe        = regexp.MustCompile(`(\d+)岁`)
    genderRe     = regexp.MustCompile(`"genderString":"([^士]+)士`)
    educationRe  = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
    xinzuoRe     = regexp.MustCompile(`([^\(]+)\(`)
    heightRe     = regexp.MustCompile(`(\d+)cm`)
    weightRe     = regexp.MustCompile(`(\d+)kg`)
    workPlaceRe  = regexp.MustCompile(`工作地:(.*)`)
    incomeRe     = regexp.MustCompile(`月收入:(.*)`)
    idRe         = regexp.MustCompile(`http://album.zhenai.com/u/(\d+)`)
)

// parseProfile 通过正则表达式获得每个字段对应的结果
func parseProfile(url string, contents []byte, name string) engine.ParseResult {
    profile := model.Profile{}

    profile.Name = name

    matches := personInfoRe.FindAllSubmatch(contents, -1)
    if len(matches) < 9 {
        return engine.ParseResult{}
    }

    profile.Marriage = string(matches[0][1])
    if age, err := strconv.Atoi(extractString(matches[1][1], ageRe)); err == nil {
        profile.Age = age
    }
    profile.Xinzuo = extractString(matches[2][1], xinzuoRe)

    if height, err := strconv.Atoi(extractString(matches[3][1], heightRe)); err == nil {
        profile.Height = height
    }

    if weight, err := strconv.Atoi(extractString(matches[4][1], weightRe)); err == nil {
        profile.Weight = weight
    }

    profile.Income = extractString(matches[6][1], incomeRe)
    profile.Occupation = string(matches[7][1])
    profile.Education = string(matches[8][1])
    profile.WorkPlace = extractString(matches[5][1], workPlaceRe)
    profile.Gender = extractString(contents, genderRe)

    // if age, err := strconv.Atoi(extractString(contents, ageRe)); err == nil {
    //     profile.Age = age
    // }
    // if height, err := strconv.Atoi(extractString(contents, heightRe)); err == nil {
    //     profile.Height = height
    // }
    // if weight, err := strconv.Atoi(extractString(contents, weightRe)); err == nil {
    //     profile.Weight = weight
    // }

    // profile.Gender = extractString(contents, genderRe)
    // profile.Income = extractString(contents, incomeRe)
    // profile.Marriage = extractString(contents, marriageRe)
    // profile.Education = extractString(contents, educationRe)
    // profile.Occupation = extractString(contents, occupationRe)
    // profile.NativePlace = extractString(contents, nativePlaceRe)
    // profile.Xinzuo = extractString(contents, xinzuoRe)
    // profile.House = extractString(contents, houseRe)
    // profile.Car = extractString(contents, carRe)

    result := engine.ParseResult{
        Items: []engine.Item{
            {
                URL:     url,
                Type:    "zhenai",
                ID:      extractString([]byte(url), idRe),
                Payload: profile,
            },
        },
    }
    return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
    match := re.FindSubmatch(contents)
    if len(match) >= 2 {
        return string(match[1])
    }
    return ""
}

// ProfileParser 用户信息对应的解析器
type ProfileParser struct {
    userName string
}

// Parse 调用parseProfile对用户信息页面进行解析
func (p *ProfileParser) Parse(url string, contents []byte) engine.ParseResult {
    return parseProfile(url, contents, p.userName)
}

// Serialize ProfileParser的序列化接口
func (p *ProfileParser) Serialize() (name string, args interface{}) {
    return "ProfileParser", p.userName
}

// NewProfileParser 用于创建ProfileParser
// name是由城市页面调用ParseCity所得的结果
func NewProfileParser(name string) *ProfileParser {
    return &ProfileParser{
        userName: name,
    }
}
