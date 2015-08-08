package main

import (
	"fmt"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

func main() {
	Auth := aws.Auth{AccessKey: `-----------`, SecretKey: `-----------`}
	AWSConnection := s3.New(Auth, aws.USEast)

	Bucket := AWSConnection.Bucket("social-images")

	bucketList, err := Bucket.List("", "", "", 100)
	fmt.Println(AWSConnection, Bucket, bucketList, err)

	if err != nil {
		fmt.Println(err.Error())
	}
	for _, item := range bucketList.Contents {
		fmt.Println(item.Key)
	}

}
