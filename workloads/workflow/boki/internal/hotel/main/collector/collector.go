package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/pkg/cayonlib"
)

func Handler() {
	cayonlib.RestartAll("gateway")
}

func main() {
	lambda.Start(Handler)
}
