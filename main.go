package main

import (
	"fmt"

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
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = endpoints.EuWest2RegionID

	// Using the Config value, create the DynamoDB client
	svc := ec2.New(cfg)

	params := &ec2.DescribeImagesInput{
		Filters: []ec2.Filter{
			{
				Name:   aws.String("image-id"),
				Values: []string{"ami-ff7d649b"},
			},
		},
		Owners: []string{"self", "amazon"},
	}

	req := svc.DescribeImagesRequest(params)
	resp, err := req.Send()
	if err != nil {
		panic("AMI DescribeImages failed, " + err.Error())
	}

	fmt.Println("Response", resp)
}
