package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cs.utexas.edu/zjia/faas-queue/common"
)

func JsonPostRequest(client *http.Client, url string, request interface{}, response interface{}) error {
	encoded, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("[FATAL] Failed to encode JSON request: %v", err)
	}
	resp, err := client.Post(url, "application/json", bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-OK response: %d", resp.StatusCode)
	}
	reader, err := common.DecompressFromReader(resp.Body)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(reader).Decode(response); err != nil {
		log.Fatalf("[FATAL] Failed to decode JSON response: %v", err)
	}
	return nil
}

func BuildFunctionUrl(gatewayAddr string, fnName string) string {
	return fmt.Sprintf("http://%s/function/%s", gatewayAddr, fnName)
}

func BuildQueueName(prefix string, index int, fifo bool) string {
	var queueName string
	if prefix == "" {
		queueName = fmt.Sprintf("queue-%d", index)
	} else {
		queueName = fmt.Sprintf("%s-queue-%d", prefix, index)
	}
	if fifo {
		queueName += ".fifo"
	}
	return queueName
}
