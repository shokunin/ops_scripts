package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type clusterNode struct {
	id      string
	ip      string
	port    int
	cmdport int
	role    string
	slaves  []string
}

func parse_nodes(nodes *redis.StringCmd) map[string]clusterNode {
	clusterNodes := make(map[string]clusterNode)
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
			clusterNodes[ln[0]] = n

			if role == "slave" {
				fmt.Println(ln[3])
				clusterNodes[ln[3]].slaves = append(clusterNodes[ln[3]].slaves, ln[0])
			}

		}
	}
	return (clusterNodes)

}

func main() {

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":30001"},
	})
	j := rdb.ClusterNodes()
	fmt.Println(parse_nodes(j))
	os.Exit(0)

}
