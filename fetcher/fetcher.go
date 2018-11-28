package fetcher

import (
    "bufio"
    "crawler/config"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"

    "golang.org/x/net/html/charset"
    "golang.org/x/text/encoding"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
)

var rateLimiter = time.Tick(time.Second / config.QPS)

// Fetch 使用http请求获得对应URL的数据
// 判断数据的编码, 最终把数据转换为utf-8编码
func Fetch(url string) ([]byte, error) {
    <-rateLimiter
    // resp, err := http.Get(url)

    request, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    request.Header.Add("User-Agent",
        "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36")

    resp, err := http.DefaultClient.Do(request)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
    }
    bodyReader := bufio.NewReader(resp.Body)
    e := determineEncoding(bodyReader)
    utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
    return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
    bytes, err := r.Peek(1024)
    if err != nil {
        log.Printf("Fetcher error: %v", err)
        return unicode.UTF8
    }
    e, _, _ := charset.DetermineEncoding(bytes, "")
    return e
}
