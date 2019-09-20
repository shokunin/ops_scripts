package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nu7hatch/gouuid"
	"os"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func worker(id int, jobs <-chan int, results chan<- int, hostname string, port int, password string) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", hostname, port),
		Password: password, // no password set
	})
	for j := range jobs {
		r, _ := uuid.NewV4()
		u := r.String()
		_, err := client.HMSet(u, map[string]interface{}{
			"key1": u,
			"key2": u,
			"key3": u,
			"key4": u,
			"key5": u,
			"key6": u,
			"key7": u,
		}).Result()
		errHndlr(err)
		results <- j
	}
	client.Close()
}

func main() {
	redisHost := flag.String("host", "localhost", "Redis Host")
	redisPort := flag.Int("port", 6379, "Redis Port")
	redisPassword := flag.String("password", "", "RedisPassword")
	messageCount := flag.Int("message_count", 100000, "run this man times")
	threadCount := flag.Int("threadcount", 10, "run this man threads")
	flag.Parse()

	jobs := make(chan int, *messageCount)
	results := make(chan int, *messageCount)

	for w := 0; w <= *threadCount; w++ {
		go worker(w, jobs, results, *redisHost, *redisPort, *redisPassword)
	}

	for j := 0; j <= *messageCount-1; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a <= *messageCount-1; a++ {
		<-results
	}
	os.Exit(0)

}
