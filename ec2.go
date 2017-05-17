package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
	InstanceID string
}

func newSession() *ec2.EC2 {
	creds := credentials.NewEnvCredentials()

	_, err := creds.Get()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	svc := ec2.New(session.New(&aws.Config{Region: aws.String("us-east-2"), Credentials: creds}))

	return svc
}
func (i *Instance) launchInstance(svc *ec2.EC2) {

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// A coreos image in us-east-2 region
		ImageId:      aws.String("ami-102f0875"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		log.Println("Could not create instance", err)
		return
	}

	log.Println("Created instance", *runResult.Instances[0].InstanceId)
	i.InstanceID = *runResult.Instances[0].InstanceId

	// Add tags to the created instance
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("MyFirstInstance"),
			},
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String("turnstile"),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}

	log.Println("Successfully tagged instance")

}

func (i *Instance) state(svc *ec2.EC2) {
	machineInfo := descirbe(svc, i.InstanceID)
	fmt.Println(*machineInfo.Reservations[0].Instances[0].State.Name)
}

func descirbe(svc *ec2.EC2, instanceID string) *ec2.DescribeInstancesOutput {
	params := &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{ // Required
			// Name: aws.String("String"),
			// Values: []*string{
			// 	aws.String("String"), // Required
			// 	// More values...
			// },
			},
			// More values...
		},
		InstanceIds: []*string{
			aws.String(instanceID), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}

	// Pretty-print the response data.
	// fmt.Println(resp)

	return resp
}
