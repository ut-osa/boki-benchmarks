package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/geo"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *cayonlib.Env) interface{} {
	req := geo.Request{}
	err := mapstructure.Decode(env.Input, &req)
	cayonlib.CHECK(err)
	return geo.Nearby(env, req)
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
