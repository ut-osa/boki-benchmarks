package main

import (
	// "github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/internal/hotel/main/hotel"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"

	"cs.utexas.edu/zjia/faas"
)

func Handler(env *cayonlib.Env) interface{} {
	var rpcInput data.RPCInput
	cayonlib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "ReserveHotel":
		return hotel.ReserveHotel(env, req["hotelId"].(string), req["userId"].(string))
	case "BaseReserveHotel":
		return hotel.BaseReserveHotel(env, req["hotelId"].(string), req["userId"].(string))
	case "AddHotel":
		hotel.AddHotel(env, req["hotelId"].(string), int32(req["cap"].(float64)))
		return 0
	}
	panic("no such function")
}
func main() {
	// lambda.Start(cayonlib.Wrapper(Handler))
	faas.Serve(cayonlib.CreateFuncHandlerFactory(Handler))
}
