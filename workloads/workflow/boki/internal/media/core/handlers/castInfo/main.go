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
	case "WriteCastInfo":
		var info core.CastInfo
		cayonlib.CHECK(mapstructure.Decode(rpcInput.Input, &info))
		core.WriteCastInfo(env, info)
		return 0
	case "ReadCastInfo":
		var castInfos []string
		cayonlib.CHECK(mapstructure.Decode(rpcInput.Input, &castInfos))
		return core.ReadCastInfo(env, castInfos)
	}
	panic("no such function")
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
