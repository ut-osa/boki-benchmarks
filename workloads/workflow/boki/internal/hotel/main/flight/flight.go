package flight

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

type Flight struct {
	FlightId  string
	Cap       int32
	Customers []string
}

func BaseReserveFlight(env *cayonlib.Env, flightId string, userId string) bool {
	item := cayonlib.Read(env, data.Tflight(), flightId)
	var flight Flight
	cayonlib.CHECK(mapstructure.Decode(item, &flight))
	if flight.Cap == 0 {
		return false
	}
	cayonlib.Write(env, data.Tflight(), flightId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.cap"): expression.Value(flight.Cap),
	})
	return true
}

func ReserveFlight(env *cayonlib.Env, flightId string, userId string) bool {
	ok, item := cayonlib.TPLRead(env, data.Tflight(), flightId)
	if !ok {
		return false
	}
	var flight Flight
	cayonlib.CHECK(mapstructure.Decode(item, &flight))
	if flight.Cap == 0 {
		return false
	}
	ok = cayonlib.TPLWrite(env, data.Tflight(), flightId,
		aws.JSONValue{"V.Cap": flight.Cap})
	return ok
}

func AddFlight(env *cayonlib.Env, flightId string, cap int32) {
	cayonlib.Write(env, data.Tflight(), flightId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(Flight{
			FlightId:  flightId,
			Cap:       cap,
			Customers: []string{},
		}),
	})
}
