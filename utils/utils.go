package utils

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Paginator struct {
	Previous int    `json:"previous"`
	Next     string `json:"next"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
}

func Printer(queue <-chan string) {
	for msg := range queue {
		fmt.Print(msg)
	}
	fmt.Println("PRINTER DONE")
}

type SafeMap struct {
	Map map[string]int
	sync.Mutex
}
type SafeWaitGroup struct {
	WaitGroup sync.WaitGroup
	sync.Mutex
}

func (swg SafeWaitGroup) Add(n int) {
	swg.Lock()
	swg.WaitGroup.Add(n)
	swg.Unlock()
}

func (swg SafeWaitGroup) Done() {
	swg.Lock()
	swg.WaitGroup.Done()
	swg.Unlock()
}

func (swg SafeWaitGroup) Wait() {
	swg.WaitGroup.Wait()
}

func WaitAndClose(wg sync.WaitGroup, grossChannel chan []interface{}, channels ...chan interface{}) {
	wg.Wait()
	for _, c := range channels {
		close(c)
	}
}

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
