package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/media/core"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *cayonlib.Env) interface{} {
	var rpcInput core.RPCInput
	cayonlib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	switch rpcInput.Function {
	case "StoreReview":
		var review core.Review
		cayonlib.CHECK(mapstructure.Decode(rpcInput.Input, &review))
		core.StoreReview(env, review)
		return 0
	case "ReadReviews":
		var reviewIds []string
		cayonlib.CHECK(mapstructure.Decode(rpcInput.Input, &reviewIds))
		return core.ReadReviews(env, reviewIds)
	}
	panic("no such function")
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
