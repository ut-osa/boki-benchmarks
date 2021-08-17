package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/internal/hotel/main/search"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *cayonlib.Env) interface{} {
	req := search.Request{}
	err := mapstructure.Decode(env.Input, &req)
	cayonlib.CHECK(err)
	return aws.JSONValue{"search": search.Nearby(env, req)}
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
