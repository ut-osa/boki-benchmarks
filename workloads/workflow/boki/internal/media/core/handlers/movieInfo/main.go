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
	req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "WriteMovieInfo":
		var info core.MovieInfo
		cayonlib.CHECK(mapstructure.Decode(req["info"], &info))
		core.WriteMovieInfo(env, info)
		return 0
	case "ReadMovieInfo":
		return core.ReadMovieInfo(env, req["movieId"].(string))
	case "UpdateRating":
		core.UpdateRating(env, req["movieId"].(string), int32(req["sumUncommittedRating"].(float64)),
			int32(req["numUncommittedRating"].(float64)))
		return 0
	}
	panic("no such function")
}

func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
