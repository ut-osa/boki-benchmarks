package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"

	"cs.utexas.edu/zjia/faas-retwis/utils"
)

var FLAGS_faas_gateway string
var FLAGS_fn_prefix string
var FLAGS_num_users int
var FLAGS_followers_per_user int
var FLAGS_concurrency int
var FLAGS_rand_seed int

func init() {
	flag.StringVar(&FLAGS_faas_gateway, "faas_gateway", "127.0.0.1:8081", "")
	flag.StringVar(&FLAGS_fn_prefix, "fn_prefix", "", "")
	flag.IntVar(&FLAGS_num_users, "num_users", 1000, "")
	flag.IntVar(&FLAGS_followers_per_user, "followers_per_user", 8, "")
	flag.IntVar(&FLAGS_concurrency, "concurrency", 1, "")
	flag.IntVar(&FLAGS_rand_seed, "rand_seed", 23333, "")

	rand.Seed(int64(FLAGS_rand_seed))
}

func createUsers() {
	client := utils.NewFaasClient(FLAGS_faas_gateway, FLAGS_concurrency)
	for i := 0; i < FLAGS_num_users; i++ {
		client.AddJsonFnCall(FLAGS_fn_prefix+"RetwisRegister", utils.JSONValue{
			"username": fmt.Sprintf("testuser_%d", i),
			"password": fmt.Sprintf("password_%d", i),
		})
	}
	results := client.WaitForResults()

	numSuccess := 0
	for _, result := range results {
		if result.Result.Success {
			numSuccess++
		}
	}
	if numSuccess < FLAGS_num_users {
		log.Printf("[ERROR] %d UserRegister requests failed", FLAGS_num_users-numSuccess)
	}
}

func createFollowers() {
	userIds1 := make([]int, 0, 1024)
	userIds2 := make([]int, 0, 1024)
	for i := 0; i < FLAGS_num_users; i++ {
		for j := 0; j < FLAGS_followers_per_user; j++ {
			followeeId := 0
			for {
				followeeId = rand.Intn(FLAGS_num_users)
				if followeeId != i {
					break
				}
			}
			userIds1 = append(userIds1, i)
			userIds2 = append(userIds2, followeeId)
		}
	}
	totalRequests := len(userIds1)
	rand.Shuffle(totalRequests, func(i, j int) {
		userIds1[i], userIds1[j] = userIds1[j], userIds1[i]
		userIds2[i], userIds2[j] = userIds2[j], userIds2[i]
	})

	client := utils.NewFaasClient(FLAGS_faas_gateway, FLAGS_concurrency)
	for i := 0; i < totalRequests; i++ {
		client.AddJsonFnCall(FLAGS_fn_prefix+"RetwisFollow", utils.JSONValue{
			"userId":     fmt.Sprintf("%08x", userIds1[i]),
			"followeeId": fmt.Sprintf("%08x", userIds2[i]),
		})
	}
	results := client.WaitForResults()

	numSuccess := 0
	for _, result := range results {
		if result.Result.Success {
			numSuccess++
		}
	}
	if numSuccess < totalRequests {
		log.Printf("[ERROR] %d Follow requests failed", totalRequests-numSuccess)
	}
}

func main() {
	flag.Parse()
	createUsers()
	createFollowers()
}
