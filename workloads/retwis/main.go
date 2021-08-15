package main

import (
	"fmt"

	"cs.utexas.edu/zjia/faas-retwis/handlers"

	"cs.utexas.edu/zjia/faas"
	"cs.utexas.edu/zjia/faas/types"
)

type funcHandlerFactory struct {
}

func (f *funcHandlerFactory) New(env types.Environment, funcName string) (types.FuncHandler, error) {
	switch funcName {
	case "RetwisInit":
		return handlers.NewSlibInitHandler(env), nil
	case "RetwisRegister":
		return handlers.NewSlibRegisterHandler(env), nil
	case "RetwisLogin":
		return handlers.NewSlibLoginHandler(env), nil
	case "RetwisProfile":
		return handlers.NewSlibProfileHandler(env), nil
	case "RetwisFollow":
		return handlers.NewSlibFollowHandler(env), nil
	case "RetwisPost":
		return handlers.NewSlibPostHandler(env), nil
	case "RetwisPostList":
		return handlers.NewSlibPostListHandler(env), nil
	case "mongoRetwisInit":
		return handlers.NewMongoInitHandler(env), nil
	case "mongoRetwisRegister":
		return handlers.NewMongoRegisterHandler(env), nil
	case "mongoRetwisLogin":
		return handlers.NewMongoLoginHandler(env), nil
	case "mongoRetwisProfile":
		return handlers.NewMongoProfileHandler(env), nil
	case "mongoRetwisFollow":
		return handlers.NewMongoFollowHandler(env), nil
	case "mongoRetwisPost":
		return handlers.NewMongoPostHandler(env), nil
	case "mongoRetwisPostList":
		return handlers.NewMongoPostListHandler(env), nil
	default:
		return nil, fmt.Errorf("Unknown function name: %s", funcName)
	}
}

func (f *funcHandlerFactory) GrpcNew(env types.Environment, service string) (types.GrpcFuncHandler, error) {
	return nil, fmt.Errorf("Not implemented")
}

func main() {
	faas.Serve(&funcHandlerFactory{})
}
