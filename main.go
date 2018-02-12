package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/davecgh/go-spew/spew"
)

// These values need to be exact or the layout formatter takes a nosedive.
const (
	AwsTimeFormat = "2006-01-02T15:04:05.000Z"
	DateTime      = "2006-01-02"
	TimeStamp     = "15:04:05.000"
)

// TODO: logical splitting of the file into multiple packages.

type AMIInfo struct {
	Name         *string
	ImageId      *string
	CreationDate *string
}

type amiInfoSlice []AMIInfo

func (a amiInfoSlice) Len() int {
	return len(a)
}

func (a amiInfoSlice) Less(i, j int) bool {
	x, _ := time.Parse(AwsTimeFormat, aws.StringValue(a[i].CreationDate))
	y, _ := time.Parse(AwsTimeFormat, aws.StringValue(a[j].CreationDate))
	return x.Before(y)
}

func (a amiInfoSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func main() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println("unable to load SDK config", err.Error())
		log.Fatal(err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = endpoints.EuWest2RegionID

	// Using the Config value, create the EC2 client
	svc := ec2.New(cfg)

	// Expand to all of the versions we want to support
	amiNames := []string{
		"Windows_Server-2012-R2_RTM-English-64Bit-Core*",
		"Windows_Server-2012-R2_RTM-English-64Bit-Base*",
	}

	// TODO: make dynamic based on arguments for the tags as well as ownership if we want to get privately owned amis.

	// Filter parameters
	params := &ec2.DescribeImagesInput{
		Filters: []ec2.Filter{
			{
				Name:   aws.String("architecture"),
				Values: []string{"x86_64"},
			}, {
				Name:   aws.String("platform"),
				Values: []string{"windows"},
			},
			{
				Name:   aws.String("name"),
				Values: amiNames,
			},
		},
		Owners: []string{"self", "amazon"},
	}

	req := svc.DescribeImagesRequest(params)
	resp, err := req.Send()
	if err != nil {
		fmt.Println("AMI DescribeImages failed, ", err.Error())
		log.Fatal(err.Error())
	}

	var amiFound *AMIInfo
	reDate := regexp.MustCompile("^([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3})Z$")
	// TODO: Can parallelize / partition the code here at the ami name because we already have the response from aws.
	// just use each of the parallelized lines to return a single ami for each name based on the dates.
	for _, name := range amiNames {
		amis := make(amiInfoSlice, 0)
		reAmiName := regexp.MustCompile(name)
		for _, ami := range resp.Images {
			// Must match ami name regex, and the string creation date regex
			if reAmiName.MatchString(aws.StringValue(ami.Name)) && reDate.MatchString(aws.StringValue(ami.CreationDate)) {
				amis = append(amis, AMIInfo{
					Name:         ami.Name,
					ImageId:      ami.ImageId,
					CreationDate: ami.CreationDate,
				})
			}
		}
		sort.Sort(amis)
		amiFound = &amis[len(amis)-1]
		// TODO: make this the return of the partitioned code
		spew.Dump(amiFound)
	}
}
