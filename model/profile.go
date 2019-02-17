package model

import (
    "encoding/json"
)

// Profile 爬虫最终获得的用户信息
type Profile struct {
    // Name 名字
    Name string

    // Gender 性别
    Gender string

    // Age 年龄
    Age int

    // Height 身高
    Height int

    // Weight 体重
    Weight int

    // Income 收入
    Income string

    // Marriage 婚姻状况
    Marriage string

    // Education 学历
    Education string

    // Occupation 职业
    Occupation string

    // WorkPlace 工作地
    WorkPlace string

    // Xinzuo 星座
    Xinzuo string

    // House 是否已购房
    House string

    // Car 是否已购车
    Car string
}

// FromJSONObj Profile的转化接口
func FromJSONObj(o interface{}) (Profile, error) {
    var profile Profile
    s, err := json.Marshal(o)
    if err != nil {
        return profile, err
    }

    err = json.Unmarshal(s, &profile)
    return profile, err
}
