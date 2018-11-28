package engine

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// ConcurrentEngine 核心调度模块，用于将request发送给Scheduler
// 创建指定数量的worker，并接收从worker中取得的结果
// 把爬虫取得Item发送给存储模块
type ConcurrentEngine struct {
	// Scheduler 用于分发request到worker
	Scheduler

	// WorkerCount 要创建worker的数量
	WorkerCount int

	// ItemChan 用于爬虫取得Item发送给存储模块
	ItemChan chan Item

	// RequestProcessor 用于把request发送给crawl service
	RequestProcessor Processor

	// RedisHost redis的host, 用于URL的去重
	RedisHost string
}

// Processor 用于把request发送给crawl service
type Processor func(r Request) (ParseResult, error)

// Scheduler 用于分发request到worker
type Scheduler interface {
	// Submit 将Request发送给RequestChan
	Submit(Request)

	// WorkerChan 用于创建一个engine.Request channel
	WorkChan() chan Request

	// Run 调度模块
	Run()

	// 当worker空闲时, 将worker的Request chan发送给WorkerChan
	ReadyNotifier
}

// ReadyNotifier 表示当worker空闲时, 将worker的Request chan发送给WorkerChan
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

var seen = make(map[string]bool)

func isDuplicate(url string) bool {
	if seen[url] == true {
		return true
	}
	seen[url] = true
	return false
}

// Run 创建指定数量的worker
// 把request发送给Scheduler
// 把爬取的item发送给存储模块
func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.WorkChan(), out, e.Scheduler)
	}

	for _, r := range seeds {
		e.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			// save item
			go func(item Item) {
				e.ItemChan <- item
			}(item)
		}
		for _, request := range result.Requests {
			// if isDuplicate(request.Url) {
			// 	continue
			// }
			e.Submit(request)
		}
	}
}

// 从channel接收request请求
// 使用redis对URL进行去重
// 将request通过rpc发送给爬虫模块
// 接收爬虫模块爬取的数据, 并通过channel发送给engine
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) error {
	redisClient, err := redis.Dial("tcp", e.RedisHost)
	if err != nil {
		log.Printf("Error: connecting to redis host: %s", e.RedisHost)
		return err
	}

	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in

			count, err := redis.Int64(redisClient.Do("SADD", "url_set", request.URL))
			if err != nil || count == 1 {
				// Ignore redis error
				if err != nil {
					log.Printf("Error: send the URL to redis: %s", request.URL)
				}

				result, err := e.RequestProcessor(request) //Worker(request)
				if err != nil {
					log.Println(err)
					continue
				}
				out <- result
			} else {
				log.Printf("Duplicate request URL: %s", request.URL)
			}
		}
	}()
	return nil
}
