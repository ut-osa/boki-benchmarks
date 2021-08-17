package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/pkg/beldilib"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *beldilib.Env) interface{} {
	return 0
}

func main() {
	// lambda.Start(beldilib.Wrapper(Handler))
	faas.Serve(beldilib.CreateFuncHandlerFactory(Handler))
}
