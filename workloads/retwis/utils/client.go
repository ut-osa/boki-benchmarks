package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpResult struct {
	Err        error
	Success    bool
	StatusCode int
	Message    string
	Duration   time.Duration
}

type JSONValue = map[string]interface{}

func JsonPostRequest(client *http.Client, url string, request JSONValue) *HttpResult {
	encoded, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("[FATAL] Failed to encode JSON request: %v", err)
	}
	start := time.Now()
	resp, err := client.Post(url, "application/json", bytes.NewReader(encoded))
	if err != nil {
		log.Printf("[ERROR] HTTP Post failed: %v", err)
		return &HttpResult{Err: err, Success: false}
	}
	elapsed := time.Since(start)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[ERROR] Non-OK response: %d", resp.StatusCode)
		return &HttpResult{Success: false, StatusCode: resp.StatusCode}
	}
	var response JSONValue
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalf("[FATAL] Failed to decode JSON response: %v", err)
	}
	success := response["success"].(bool)
	if !success {
		message := response["message"].(string)
		log.Printf("[WARN] Request failed: %v", message)
		return &HttpResult{
			Success:    false,
			StatusCode: 200,
			Message:    message,
			Duration:   elapsed,
		}
	}
	return &HttpResult{
		Success:    true,
		StatusCode: 200,
		Duration:   elapsed,
	}
}

func BuildFunctionUrl(gatewayAddr string, fnName string) string {
	return fmt.Sprintf("http://%s/function/%s", gatewayAddr, fnName)
}

type FaasCall struct {
	FnName string
	Input  JSONValue
	Result *HttpResult
}

type faasWorker struct {
	gateway string
	client  *http.Client
	reqChan chan *FaasCall
	wg      *sync.WaitGroup
	results []*FaasCall
}

func (w *faasWorker) start() {
	defer w.wg.Done()
	for {
		call, more := <-w.reqChan
		if !more {
			break
		}
		url := BuildFunctionUrl(w.gateway, call.FnName)
		call.Result = JsonPostRequest(w.client, url, call.Input)
		call.Input = nil
		w.results = append(w.results, call)
	}
}

type FaasClient struct {
	reqChan chan *FaasCall
	workers []*faasWorker
	wg      *sync.WaitGroup
}

func NewFaasClient(faasGateway string, concurrency int) *FaasClient {
	reqChan := make(chan *FaasCall, concurrency)
	workers := make([]*faasWorker, concurrency)
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		worker := &faasWorker{
			gateway: faasGateway,
			client: &http.Client{
				Transport: &http.Transport{
					MaxConnsPerHost: 1,
					MaxIdleConns:    1,
					IdleConnTimeout: 30 * time.Second,
				},
				Timeout: 4 * time.Second,
			},
			reqChan: reqChan,
			wg:      wg,
			results: make([]*FaasCall, 0, 128),
		}
		go worker.start()
		workers[i] = worker
	}
	return &FaasClient{
		reqChan: reqChan,
		workers: workers,
		wg:      wg,
	}
}

func (c *FaasClient) AddJsonFnCall(fnName string, input JSONValue) {
	call := &FaasCall{
		FnName: fnName,
		Input:  input,
		Result: nil,
	}
	c.reqChan <- call
}

func (c *FaasClient) WaitForResults() []*FaasCall {
	close(c.reqChan)
	c.wg.Wait()
	results := make([]*FaasCall, 0)
	for _, worker := range c.workers {
		results = append(results, worker.results...)
	}
	return results
}
