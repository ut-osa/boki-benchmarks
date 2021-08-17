package cayonlib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"time"
)

func CreateMainTable(lambdaId string) {
	_, _ = DBClient.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("K"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("K"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName: aws.String(kTablePrefix + lambdaId),
	})
}

func CreateLogTable(lambdaId string) {
	panic("Not implemented")
}

func CreateCollectorTable(lambdaId string) {
	panic("Not implemented")
}

func CreateBaselineTable(lambdaId string) {
	panic("Not implemented")
}

func CreateLambdaTables(lambdaId string) {
	CreateMainTable(lambdaId)
	// CreateLogTable(lambdaId)
	// CreateCollectorTable(lambdaId)
}

func CreateTxnTables(lambdaId string) {
	CreateBaselineTable(lambdaId)
	CreateLogTable(lambdaId)
	CreateCollectorTable(lambdaId)
}

func DeleteTable(tablename string) {
	_, _ = DBClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(kTablePrefix + tablename)})
}

func DeleteLambdaTables(lambdaId string) {
	DeleteTable(lambdaId)
	// DeleteTable(fmt.Sprintf("%s-log", lambdaId))
	// DeleteTable(fmt.Sprintf("%s-collector", lambdaId))
}

func WaitUntilDeleted(tablename string) {
	for ; ; {
		res, err := DBClient.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(kTablePrefix + tablename)})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeResourceNotFoundException:
					return
				}
			}
		} else if *res.Table.TableStatus != "DELETING" {
			DeleteTable(tablename)
		}
		time.Sleep(3 * time.Second)
	}
}

func WaitUntilAllDeleted(tablenames []string) {
	for _, tablename := range tablenames {
		WaitUntilDeleted(tablename)
	}
}

func WaitUntilActive(tablename string) bool {
	counter := 0
	for ; ; {
		res, err := DBClient.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(kTablePrefix + tablename)})
		if err != nil {
			counter += 1
			fmt.Printf("%s DescribeTable error: %v\n", tablename, err)
		} else {
			if *res.Table.TableStatus == "ACTIVE" {
				return true
			}
			fmt.Printf("%s status: %s\n", tablename, *res.Table.TableStatus)
			if *res.Table.TableStatus != "CREATING" && counter > 6 {
				return false
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func WaitUntilAllActive(tablenames []string) bool {
	for _, tablename := range tablenames {
		res := WaitUntilActive(tablename)
		if !res {
			return false
		}
	}
	return true
}

func WriteHead(tablename string, key string) {
	panic("Not implemented")
}

func WriteTail(tablename string, key string, row string) {
	panic("Not implemented")
}

func WriteNRows(tablename string, key string, n int) {
	panic("Not implemented")
}

func Populate(tablename string, key string, value interface{}, baseline bool) {
	LibWrite(tablename, aws.JSONValue{"K": key},
		map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("VERSION"): expression.Value(0),
			expression.Name("V"):       expression.Value(value),
		})
}
