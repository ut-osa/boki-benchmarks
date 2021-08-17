package main

import (
	"sync"
	"time"
	"log"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func main() {
	services := []string{"flight", "hotel", "order"}
	statics := []string{"user", "search", "recommendation", "rate", "profile", "geo", "gateway", "frontend"}

	for {
		var wg sync.WaitGroup
		for _, service := range services {
			wg.Add(1)
			go func(service string) {
				defer wg.Done()
				log.Printf("[INFO] Start GC: %s", service)
				beldilib.GC(service)
			}(service)
		}
		for _, service := range statics {
			wg.Add(1)
			go func(service string) {
				defer wg.Done()
				log.Printf("[INFO] Start static GC: %s", service)
				beldilib.StaticGC(service)
			}(service)
		}
		wg.Wait()
		time.Sleep(100 * time.Millisecond)
	}
}
