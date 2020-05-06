package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var rHost string
var rPort int

func errHndlr(err error, s string) {
	if err != nil {
		fmt.Println("error:", s, " - ", err)
	}
}

func randomDialer() (net.Conn, error) {
	ips, reserr := net.LookupIP(rHost)
	if reserr != nil {
		return nil, reserr
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	n := rand.Int() % len(ips)
	fmt.Println("Dialing....")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ips[n], rPort))
	return conn, err
}

func main() {
	redisHost := flag.String("host", "localhost", "Redis Host")
	redisPort := flag.Int("port", 6379, "Redis Port")
	redisPassword := flag.String("password", "", "RedisPassword")
	messageCount := flag.Int("message_count", 100000, "run this man times")
	flag.Parse()
	rHost = *redisHost
	rPort = *redisPort

	client := redis.NewClient(&redis.Options{
		Dialer:          randomDialer, // Randomly pick an IP address from the list of ips retruned
		Password:        *redisPassword,
		DB:              0,
		MinIdleConns:    1,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 5000 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      50, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	err := client.Set("counter", 0, 0).Err()
	errHndlr(err, "initial set")
	fmt.Println("Key reset sleeping works")
	time.Sleep(10 * time.Second)
	fmt.Println("done sleeping")

	cnt := 0
	for j := 0; j <= *messageCount-1; j++ {
		_, e := client.Incr("counter").Result()
		errHndlr(e, "Incr it")
		if e != nil {
			for {
				fmt.Println("Retrying....")
				time.Sleep(1 * time.Second)
				c, perr := client.Get("counter").Result()
				cnt, _ = strconv.Atoi(c)
				if perr != nil {
					fmt.Println(perr)
				} else {
					fmt.Println("break:", c, "-", j)
					break
				}
			}
			if cnt < j {
				_, e2 := client.Incr("counter").Result()
				errHndlr(e2, "Incr it retry")
			}
		}
	}

	j := client.Get("counter")
	fmt.Println(j)

}
