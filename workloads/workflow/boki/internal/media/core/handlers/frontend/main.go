package main

import (
	"fmt"
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
	case "Compose":
		var input core.ComposeInput
		cayonlib.CHECK(mapstructure.Decode(rpcInput.Input, &input))
		core.Compose(env, input)
		return 0
	}
	fmt.Println("ERROR: no such function")
	panic(rpcInput)
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
