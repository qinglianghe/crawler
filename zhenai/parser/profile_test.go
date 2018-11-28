package parser

import (
	"crawler/engine"
	"crawler/model"
	"io/ioutil"
	"testing"
)

func TestParseProfile(t *testing.T) {
	contents, err := ioutil.ReadFile("profile_test_data_01.html")
	if err != nil {
		panic(err)
	}

	result := parseProfile("http://album.zhenai.com/u/1439637023", contents, "一切随缘")

	if len(result.Items) != 1 {
		t.Errorf("Items should contain 1 element; but was %v", result.Items)
	}

	actual := result.Items[0]

	expected := engine.Item{
		URL:  "http://album.zhenai.com/u/1439637023",
		Type: "zhenai",
		ID:   "1439637023",
		Payload: model.Profile{
			Name:       "一切随缘",
			Gender:     "女",
			Age:        23,
			Height:     160,
			Weight:     48,
			Income:     "8千-1.2万",
			Marriage:   "未婚",
			Education:  "中专",
			Occupation: "其他职业",
			// NativePlace: "广东广州",
			WorkPlace: "广州黄埔区",
			Xinzuo:    "狮子座",
			// House:       "打算婚后购房",
			// Car:         "未买车",
		},
	}

	if expected != actual {
		t.Errorf("expected profile %v; but was %v", expected, actual)
	}
}
