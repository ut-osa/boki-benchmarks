package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/internal/hotel/main/flight"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *beldilib.Env) interface{} {
	var rpcInput data.RPCInput
	beldilib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "ReserveFlight":
		return flight.ReserveFlight(env, req["flightId"].(string), req["userId"].(string))
	case "BaseReserveFlight":
		return flight.BaseReserveFlight(env, req["flightId"].(string), req["userId"].(string))
	case "AddFlight":
		flight.AddFlight(env, req["flightId"].(string), int32(req["cap"].(float64)))
		return 0
	}
	panic("no such function")
}

func main() {
	// lambda.Start(beldilib.Wrapper(Handler))
	faas.Serve(beldilib.CreateFuncHandlerFactory(Handler))
}
