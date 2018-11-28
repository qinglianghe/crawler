package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCity(t *testing.T) {
	contents, err := ioutil.ReadFile("city_test_data.html")
	if err != nil {
		panic(err)
	}

	result := ParseCity("", contents)
	const resultSize = 42
	expectedUrls := []string{
		"http://album.zhenai.com/u/105379158",
		"http://album.zhenai.com/u/1098986823",
		"http://album.zhenai.com/u/98571982",
	}

	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d requests; but had %d", resultSize, len(result.Requests))
	}
	for i, url := range expectedUrls {
		if url != result.Requests[i].URL {
			t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].URL)
		}
	}

	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d cities; but had %d", resultSize, len(result.Requests))
	}
}
