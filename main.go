package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

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

	// Using the Config value, create the DynamoDB client
	svc := ec2.New(cfg)

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
				Name: aws.String("name"),
				Values: []string{
					"Windows_Server-2012-R2_RTM-English-64Bit-Core*",
					"Windows_Server-2012-R2_RTM-English-64Bit-Base*"},
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

	regexDate := regexp.MustCompile("^([0-9]{4}-[0-9]{2}-[0-9]{2})T([0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}Z)$")
	for _, ami := range resp.Images {
		// creationDate, _ := time.Parse(dateLayout, aws.StringValue(ami.CreationDate))
		// if creationDate.Year() >= minDate.Year() {
		// fmt.Println("Response", ami)
		// }
		if regexDate.MatchString(aws.StringValue(ami.CreationDate)) {
			// if regexDate.FindStringSubmatch(aws.StringValue(ami.CreationDate)) {

			// }
			fmt.Printf("%s: %s => %s\n", aws.StringValue(ami.Name), aws.StringValue(ami.ImageId), aws.StringValue(ami.CreationDate))
		}
	}
}
