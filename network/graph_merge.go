package main

// MAKE SURE TO SET export GOMAXPROCS=16  in your environment before running

import (
	"flag"
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	rg "github.com/maguec/redisgraph-go"
)

func graphinit(hostame string, port int, keyname string) {
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	defer conn.Close()
	graph := rg.Graph{}.New(keyname, conn)
	p := &rg.Node{
		Label: "Entry",
		Properties: map[string]interface{}{
			"EntryID": -1,
		},
	}
	graph.AddNode(p)
	graph.Flush()
	query := ("CREATE INDEX ON :Entry(EntryID)")
	graph.Query(query)
}

func worker(id int, jobs <-chan int, results chan<- string, hostame string, port int, keyname string) {
	conn, _ := redis.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	defer conn.Close()
	graph := rg.Graph{}.New(keyname, conn)

	for j := range jobs {
		query := fmt.Sprintf("MERGE (e:Entry{EntryID: '%d'}) SET e.FirstName='Chris', e.LastName='Mague', e.Title='Troublemaker', e.Location='San Francisco'", j)
		graph.Query(query)
		fmt.Println(j)
		results <- "OK"
	}
}

var hostname string
var keyname string
var port int
var concurrent int
var count int

func init() {
	flag.StringVar(&keyname, "keyname", "merge_test", "key name of the graphdb")
	flag.StringVar(&hostname, "hostname", "localhost", "hostname or ip to scan")
	flag.IntVar(&port, "port", 6379, "port to try to connect to")
	flag.IntVar(&concurrent, "concurrent", 10, "number of workers to run")
	flag.IntVar(&count, "count", 1000, "number of records")
	flag.Parse()
}

func main() {

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan int, count)
	results := make(chan string, count)

	graphinit(hostname, port, keyname)

	for w := 0; w < concurrent; w++ {
		go worker(w, jobs, results, hostname, port, keyname)
	}

	for j := 0; j <= count-1; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a <= count-1; a++ {
		<-results
	}
	os.Exit(0)
}
