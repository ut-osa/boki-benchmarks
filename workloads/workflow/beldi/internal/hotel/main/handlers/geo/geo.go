package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/geo"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *beldilib.Env) interface{} {
	req := geo.Request{}
	err := mapstructure.Decode(env.Input, &req)
	beldilib.CHECK(err)
	return geo.Nearby(env, req)
}

func main() {
	// lambda.Start(beldilib.Wrapper(Handler))
	faas.Serve(beldilib.CreateFuncHandlerFactory(Handler))
}
