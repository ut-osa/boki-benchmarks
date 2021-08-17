package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *cayonlib.Env) interface{} {
	var rpcInput data.RPCInput
	cayonlib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	//req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "search":
		res, _ := cayonlib.SyncInvoke(env, data.Tsearch(), rpcInput.Input)
		return res
	case "recommend":
		res, _ := cayonlib.SyncInvoke(env, data.Trecommendation(), rpcInput.Input)
		return res
	case "user":
		res, _ := cayonlib.SyncInvoke(env, data.Tuser(), rpcInput.Input)
		return res
	case "reserve":
		res, _ := cayonlib.SyncInvoke(env, data.Tfrontend(), rpcInput.Input)
		return res
	}
	return 0
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
