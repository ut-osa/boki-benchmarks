package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"os"
	// "time"
)

func ClearAll() {
	cayonlib.DeleteLambdaTables("singleop")
	cayonlib.DeleteLambdaTables("nop")
	// beldilib.DeleteTable("bsingleop")
	// beldilib.DeleteTable("bnop")
	// beldilib.DeleteLambdaTables("tsingleop")
	// beldilib.DeleteLambdaTables("tnop")
}

func main() {
	if len(os.Args) >= 2 {
		option := os.Args[1]
		if option == "clean" {
			ClearAll()
			return
		}
	}
	ClearAll()
	cayonlib.WaitUntilAllDeleted([]string{"singleop", "nop"})

	cayonlib.CreateLambdaTables("singleop")
	cayonlib.CreateLambdaTables("nop")

	cayonlib.WaitUntilAllActive([]string{"singleop", "nop"})

	// beldilib.CreateBaselineTable("bsingleop")
	// beldilib.CreateBaselineTable("bnop")

	// beldilib.WaitUntilAllActive([]string{
	// 	"singleop", "singleop-log", "singleop-collector",
	// 	"nop", "nop-log", "nop-collector",
	// 	"bsingleop", "bnop",
	// })

	// beldilib.CreateTxnTables("tsingleop")
	// beldilib.CreateTxnTables("tnop")

	// time.Sleep(60 * time.Second)
	// beldilib.WriteNRows("singleop", "K", 20)

	cayonlib.LibWrite("singleop", aws.JSONValue{"K": "K"}, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(1),
	})

	// beldilib.LibWrite("tsingleop", aws.JSONValue{"K": "K"}, map[expression.NameBuilder]expression.OperandBuilder{
	// 	expression.Name("V"): expression.Value(1),
	// })
}
