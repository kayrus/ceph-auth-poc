package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	region   string
	endpoint string
	access   string
	secret   string
	requests int
)

func init() {
	flag.StringVar(&region, "region", "", "The `region` to list buckets from.")
	flag.StringVar(&endpoint, "endpoint", "", "The `endpoint` to list buckets from.")
	flag.StringVar(&access, "access", "", "The `access key` to list buckets from.")
	flag.StringVar(&secret, "secret", "", "The `secret key` to list buckets from.")
	flag.IntVar(&requests, "requests", 1, "The number of simultaneous requests.")
}

func main() {
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access, secret, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			})),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	input := &s3.ListBucketsInput{}

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := client.ListBuckets(context.TODO(), input)
			if err != nil {
				log.Fatalf("unable to list buckets %v", err)
			}
			log.Printf("Listed buckets #%d successful\n", i)
		}(i)
	}

	wg.Wait()
}
