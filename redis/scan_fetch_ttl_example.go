package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var rHost string
var rPort int

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func worker(id int, jobs <-chan string, results chan<- string, redisClient *redis.Client) {
	for j := range jobs {
		l, err := redisClient.TTL(j).Result()
		errHndlr(err)
		results <- fmt.Sprintf("%s %d", j, l/time.Second)
	}
}

func scanKeys(redisClient *redis.Client, pattern string, batch int64) []string {
	var cursor uint64
	var allkeys []string
	for {
		var keys []string
		var err error
		keys, cursor, err = redisClient.Scan(cursor, pattern, batch).Result()
		if err != nil {
			panic(err)
		}
		for _, k := range keys {
			allkeys = append(allkeys, k)
		}

		if cursor == 0 {
			break
		}
	}
	return (allkeys)
}

func main() {
	redisHost := flag.String("host", "localhost", "Redis Host")
	redisPort := flag.Int("port", 6379, "Redis Port")
	redisPassword := flag.String("password", "", "RedisPassword")
	pattern := flag.String("pattern", "*", "Search Expression")
	threadCount := flag.Int("threadcount", 10, "run this man threads")
	batchSize := flag.Int64("batch-size", 10, "Scan for this many at a time")
	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", *redisHost, *redisPort),
		Password:        *redisPassword,
		DB:              0,
		MinIdleConns:    1,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 512 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      10, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	k := scanKeys(client, *pattern, *batchSize)
	keys := make(chan string, len(k))
	results := make(chan string, len(k))
	for _, j := range k {
		keys <- j
	}

	for w := 0; w <= *threadCount; w++ {
		go worker(w, keys, results, client)
	}
	close(keys)

	// Finally we collect all the results of the work.
	for a := 0; a < len(k); a++ {
		w := <-results
		fmt.Println(w)
	}

	os.Exit(0)

}
