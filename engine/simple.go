package engine

import (
    "log"
)

type SimpleEngine struct{}

func (SimpleEngine) Run(seeds ...Request) {
    var requests []Request
    for _, s := range seeds {
        requests = append(requests, s)
    }

    itemCount := 0
    for len(requests) > 0 {
        r := requests[0]
        requests = requests[1:]
        if isDuplicate(r.URL) {
            continue
        }
        parseResult, err := Worker(r)
        if err != nil {
            continue
        }
        requests = append(requests, parseResult.Requests...)

        for _, item := range parseResult.Items {
            log.Printf("Got Item #%d: %v", itemCount, item)
            itemCount++
        }
    }
}
