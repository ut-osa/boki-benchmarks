package frontend

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
)

func SendRequest(env *cayonlib.Env, userId string, flightId string, hotelId string) string {
	if cayonlib.TYPE == "BASELINE" {
		input := map[string]string{
			"hotelId": hotelId,
			"userId":  userId,
		}
		cayonlib.SyncInvoke(env, data.Thotel(), data.RPCInput{
			Function: "BaseReserveHotel",
			Input:    input,
		})
		input = map[string]string{
			"flightId": flightId,
			"userId":   userId,
		}
		cayonlib.SyncInvoke(env, data.Tflight(), data.RPCInput{
			Function: "BaseReserveFlight",
			Input:    input,
		})
		input = map[string]string{
			"flightId": flightId,
			"hotelId":  hotelId,
			"userId":   userId,
		}
		cayonlib.AsyncInvoke(env, data.Torder(), data.RPCInput{
			Function: "PlaceOrder",
			Input:    input,
		})
		return ""
	}
	cayonlib.BeginTxn(env)
	input := map[string]string{
		"hotelId": hotelId,
		"userId":  userId,
	}
	res, _ := cayonlib.SyncInvoke(env, data.Thotel(), data.RPCInput{
		Function: "ReserveHotel",
		Input:    input,
	})
	if !res.(bool) {
		cayonlib.AbortTxn(env)
		return "Place Order Fails"
	}
	input = map[string]string{
		"flightId": flightId,
		"userId":   userId,
	}
	res, _ = cayonlib.SyncInvoke(env, data.Tflight(), data.RPCInput{
		Function: "ReserveFlight",
		Input:    input,
	})
	if !res.(bool) {
		cayonlib.AbortTxn(env)
		return "Place Order Fails"
	}
	input = map[string]string{
		"flightId": flightId,
		"hotelId":  hotelId,
		"userId":   userId,
	}
	cayonlib.CommitTxn(env)
	cayonlib.AsyncInvoke(env, data.Torder(), data.RPCInput{
		Function: "PlaceOrder",
		Input:    input,
	})
	return "Place Order Success"
}
