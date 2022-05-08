package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"time"
)

func main() {
	ctx := context.TODO()
	//cfg, err := config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
	//	o.Region = "us-east-1"
	//	return nil
	//})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedCredentialsFiles(
			[]string{"test/credentials", "data/credentials"},
		),
		config.WithSharedConfigFiles(
			[]string{"test/config", "data/config"},
		),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "", SecretAccessKey: "",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	tn := "Users"
	out, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("UserId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("UserId"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(tn),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func waitForTable(ctx context.Context, db *dynamodb.Client, tn string) error {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tn),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return errors.Wrap(err, "timed out while waiting for table to become active")
	}

	return err
}
