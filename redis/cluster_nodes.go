package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var clusterNodes []clusterNode

type clusterNode struct {
	id      string
	ip      string
	port    int
	cmdport int
	role    string
	slaves  []string
}

func parse_nodes(nodes *redis.StringCmd) {
	for _, line := range strings.Split(nodes.Val(), "\n") {
		ln := strings.Split(line, " ")
		if len(ln) > 1 {
			role := "slave"
			r := regexp.MustCompile(`(\S+):(\d+)@(\d+)`)
			res := r.FindStringSubmatch(ln[1])
			match, _ := regexp.MatchString("master", ln[2])
			if match {
				role = "master"
			}

			i, _ := strconv.Atoi(res[2])
			j, _ := strconv.Atoi(res[3])

			n := clusterNode{
				id:      ln[0],
				role:    role,
				ip:      res[1],
				port:    i,
				cmdport: j,
			}
			clusterNodes = append(clusterNodes, n)

			if role == "slave" {
				for i, v := range clusterNodes {
					if v.id == ln[3] {
						clusterNodes[i].slaves = append(clusterNodes[i].slaves, ln[0])
					}
				}
			}

		}
	}

}

func main() {

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":30001"},
	})
	j := rdb.ClusterNodes()
	parse_nodes(j)
	fmt.Println(clusterNodes)
	os.Exit(0)

}
