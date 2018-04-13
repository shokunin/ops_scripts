package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/araddon/dateparse"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type CatIndices []struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

var (
	esHost = flag.String("host", "localhost", "Elasticsearch host defaults to localhost")
	esPort = flag.Int("port", 80, "Elasticsearch port defaults to 80")
	days   = flag.Int("days", 100, "Number of days to keep")
	isDry  = flag.Bool("dry-run", false, "Output dry run information")
)

func delIndex(esHost string, esPort int, index string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s:%d/%s", esHost, esPort, index), nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error deleting index:", fmt.Sprintf("http://%s:%d/%s", esHost, esPort, index))
	} else {
		fmt.Println("Sucessfully deleted index:", fmt.Sprintf("http://%s:%d/%s", esHost, esPort, index), "Response:", resp.Status)
	}
}

func main() {

	flag.Parse()

	url := fmt.Sprintf("http://%s:%d/_cat/indices?format=json", *esHost, *esPort)
	r, _ := regexp.Compile(`^logstash-(\d{4}\.\d{2}\.\d{2}$)`)

	esClient := http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error fetching from:", url, ":", err)
		os.Exit(1)
	}
	res, getErr := esClient.Do(req)
	if getErr != nil {
		fmt.Println("Error fetching from:", url, ":", getErr)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		fmt.Println("Error fetching from:", url)
		os.Exit(1)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println("Error decoding:", readErr)
		os.Exit(1)
	}
	indices := CatIndices{}
	jsonErr := json.Unmarshal(body, &indices)
	if jsonErr != nil {
		fmt.Println("Error unmarshal:", jsonErr)
		os.Exit(1)
	}

	if *isDry {
		fmt.Println("Dry run results:")
		fmt.Println("older than", *days, "days")
	}

	for _, i := range indices {
		if r.MatchString(i.Index) == true {
			m := r.FindStringSubmatch(i.Index)[1]
			d, err := dateparse.ParseLocal(strings.Replace(m, ".", "-", -1))
			if err != nil {
				fmt.Println("Unable to parse date:", m)
			}
			delta := time.Now().Sub(d)
			if int(delta.Hours()/24) >= *days {
				if *isDry {
					fmt.Println("delete:", i.Index)
				} else {
					delIndex(*esHost, *esPort, i.Index)
				}
			}

		}

	}

}
