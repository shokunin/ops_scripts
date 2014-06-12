package main

import "fmt"
import "flag"
import "time"
import "bytes"
import "os"
import "math/rand"
import "github.com/fzzy/radix/redis"

// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.

func randInt(min int, max int) int {
    rand.Seed(time.Now().UTC().UnixNano())
    return min + rand.Intn(max-min)
}

func randomString(l int) string {
    var result bytes.Buffer
    var temp string
    for i := 0; i < l; {
        if string(randInt(65, 90)) != temp {
            temp = string(randInt(65, 90))
            result.WriteString(temp)
            i++
        }
    }
    return result.String()
}

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		c, err := redis.DialTimeout("tcp", "172.19.48.11:6380", time.Duration(10)*time.Second)
		errHndlr(err)
		r := c.Cmd("select", 1)
		errHndlr(r.Err)
		s := randomString(20)
		now := time.Now()
		r = c.Cmd("set", s, s)
		r = c.Cmd("get", s)
		later := time.Now()
		fmt.Println("request:", j, "time:", later.Sub(now))
		c.Close()
		results <- 1
	}
}

func init() {
	flag.StringVar(&hostname, "hostname", "localhost", "hostname or ip to scan")
	flag.IntVar(&port, "port", 6379, "port to try to connect to")
	flag.IntVar(&concurrent, "concurrent", 10, "number of workers to run")
	flag.IntVar(&requests, "requests", 10000, "number of requests to run")
}

func main() {

    // In order to use our pool of workers we need to send
    // them work and collect their results. We make 2
    // channels for this.
    jobs := make(chan int, 1000000)
    results := make(chan int, 1000000)

    // This starts up 3 workers, initially blocked
    // because there are no jobs yet.
    for w := 1; w <= 128; w++ {
        go worker(w, jobs, results)
    }

    // Here we send 9999999 `jobs` and then `close` that
    // channel to indicate that's all the work we have.
    for j := 1; j <= 10000; j++ {
        jobs <- j
    }
    close(jobs)

    // Finally we collect all the results of the work.
    for a := 1; a <= 1000000; a++ {
        <-results
    }
}

