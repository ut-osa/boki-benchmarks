package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/internal/hotel/main/recommendation"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *cayonlib.Env) interface{} {
	req := recommendation.Request{}
	err := mapstructure.Decode(env.Input, &req)
	cayonlib.CHECK(err)
	res := recommendation.GetRecommendations(env, req)
	return aws.JSONValue{"recommend": res}
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
