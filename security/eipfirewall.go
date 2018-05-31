package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8140, "port to open up to EIPs")
	flag.Parse()
}

func getregions() []string {
	var regions = []string{}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	svc := ec2.New(sess)
	input := &ec2.DescribeRegionsInput{}

	result, err := svc.DescribeRegions(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return regions
	}

	//fmt.Println(result.Regions)
	for _, j := range result.Regions {
		regions = append(regions, *j.RegionName)
	}
	return regions

}

func geteips(region string) []string {
	var eips = []string{}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	svc := ec2.New(sess)
	input := &ec2.DescribeAddressesInput{}

	result, err := svc.DescribeAddresses(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return eips
	}

	for _, j := range result.Addresses {
		eips = append(eips, *j.PublicIp)
	}
	return eips
}

func main() {
	regions := getregions()
	for _, j := range regions {
		for _, q := range geteips(j) {
			cmd := exec.Command("/usr/sbin/ufw", "allow", "from", q, "to", "any", "port", strconv.Itoa(port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("error adding ", q)
			}
		}
	}
}
