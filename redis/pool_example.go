package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
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

func worker(id int, jobs <-chan int, results chan<- int, redisClient *redis.Client) {
	for j := range jobs {
		pong, err := redisClient.Ping().Result()
		errHndlr(err)
		fmt.Println(pong, ",", id)
		results <- j
	}
}

//func randomDialer(redisHost string, redisPort int) (net.Conn, error) {
func randomDialer() (net.Conn, error) {
	ips, reserr := net.LookupIP(rHost)
	if reserr != nil {
		return nil, reserr
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	n := rand.Int() % len(ips)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ips[n], rPort))
	return conn, err
}

func main() {
	redisHost := flag.String("host", "localhost", "Redis Host")
	redisPort := flag.Int("port", 6379, "Redis Port")
	redisPassword := flag.String("password", "", "RedisPassword")
	messageCount := flag.Int("message_count", 100000, "run this man times")
	threadCount := flag.Int("threadcount", 10, "run this man threads")
	flag.Parse()
	rHost = *redisHost
	rPort = *redisPort

	client := redis.NewClient(&redis.Options{
		Dialer:          randomDialer, // Randomly pick an IP address from the list of ips retruned
		Password:        *redisPassword,
		DB:              0,
		MinIdleConns:    1,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 512 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      10, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	jobs := make(chan int, *messageCount)
	results := make(chan int, *messageCount)

	for w := 0; w <= *threadCount; w++ {
		go worker(w, jobs, results, client)
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
