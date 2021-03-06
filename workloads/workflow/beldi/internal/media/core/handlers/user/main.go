package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/media/core"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *beldilib.Env) interface{} {
	var rpcInput core.RPCInput
	beldilib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "RegisterUser":
		core.RegisterUser(env, req["firstName"].(string), req["lastName"].(string),
			req["username"].(string), req["password"].(string))
		return 0
	case "Login":
		return core.Login(env, req["username"].(string), req["password"].(string))
	case "UploadUser":
		core.UploadUser(env, req["reqId"].(string), req["username"].(string))
		return 0
	}
	panic("no such function")
}

func main() {
	// lambda.Start(beldilib.Wrapper(Handler))
	faas.Serve(beldilib.CreateFuncHandlerFactory(Handler))
}
