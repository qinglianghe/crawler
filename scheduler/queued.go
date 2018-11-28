package scheduler

import (
    "crawler/engine"
)

// QueuedScheduler 用于分发request到worker
type QueuedScheduler struct {
    // RequestChan 用于接收从engine发过来的request
    RequestChan chan engine.Request

    // WorkerChan 用于把RequestChan的request发送给worker
    WorkerChan chan chan engine.Request
}

// Submit 将Request发送给RequestChan
func (s *QueuedScheduler) Submit(r engine.Request) {
    s.RequestChan <- r
}

// WorkerReady 将worker的Request chan发送给WorkerChan
func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
    s.WorkerChan <- w
}

// WorkChan 用于创建一个engine.Request channel
func (s *QueuedScheduler) WorkChan() chan engine.Request {
    return make(chan engine.Request)
}

// Run 用于分发request到worker
// 从RequestChan接收engine发过来的request, 放入到requestQueue队列
// 从WorkerChan接收空闲的worker, 放入到requestQueue队列
// 当requestQueue队列有request时, 且WorkerChan有空闲的worker, 则把request发送给指定的worker
func (s *QueuedScheduler) Run() {
    s.RequestChan = make(chan engine.Request)
    s.WorkerChan = make(chan chan engine.Request)

    go func() {
        var requestQueue []engine.Request
        var workerQueue []chan engine.Request
        for {
            var activeRequest engine.Request
            var activeWorker chan engine.Request
            if len(requestQueue) > 0 && len(workerQueue) > 0 {
                activeRequest = requestQueue[0]
                activeWorker = workerQueue[0]
            }

            select {
            case r := <-s.RequestChan:
                requestQueue = append(requestQueue, r)
            case w := <-s.WorkerChan:
                workerQueue = append(workerQueue, w)
            case activeWorker <- activeRequest:
                requestQueue = requestQueue[1:]
                workerQueue = workerQueue[1:]
            }
        }
    }()
}
