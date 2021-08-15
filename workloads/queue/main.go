package main

import (
	"fmt"

	"cs.utexas.edu/zjia/faas-queue/handlers"

	"cs.utexas.edu/zjia/faas"
	"cs.utexas.edu/zjia/faas/types"
)

type funcHandlerFactory struct {
}

func (f *funcHandlerFactory) New(env types.Environment, funcName string) (types.FuncHandler, error) {
	switch funcName {
	case "slibQueueProducer":
		return handlers.NewSlibProducerHandler(env), nil
	case "slibQueueConsumer":
		return handlers.NewSlibConsumerHandler(env), nil
	case "sqsInitQueue":
		return handlers.NewSqsInitHandler(env), nil
	case "sqsQueueProducer":
		return handlers.NewSqsProducerHandler(env), nil
	case "sqsQueueConsumer":
		return handlers.NewSqsConsumerHandler(env), nil
	case "pulsarQueueProducer":
		return handlers.NewPulsarProducerHandler(env), nil
	case "pulsarQueueConsumer":
		return handlers.NewPulsarConsumerHandler(env), nil
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
