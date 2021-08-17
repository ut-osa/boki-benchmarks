package hotel

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

type Hotel struct {
	HotelId   string
	Cap       int32
	Customers []string
}

func BaseReserveHotel(env *cayonlib.Env, hotelId string, userId string) bool {
	item := cayonlib.Read(env, data.Thotel(), hotelId)
	var hotel Hotel
	cayonlib.CHECK(mapstructure.Decode(item, &hotel))
	if hotel.Cap == 0 {
		return false
	}
	cayonlib.Write(env, data.Thotel(), hotelId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.Cap"): expression.Value(hotel.Cap),
	})
	return true
}

func ReserveHotel(env *cayonlib.Env, hotelId string, userId string) bool {
	ok, item := cayonlib.TPLRead(env, data.Thotel(), hotelId)
	if !ok {
		return false
	}
	var hotel Hotel
	cayonlib.CHECK(mapstructure.Decode(item, &hotel))
	if hotel.Cap == 0 {
		return false
	}
	ok = cayonlib.TPLWrite(env, data.Thotel(), hotelId,
		aws.JSONValue{"V.Cap": hotel.Cap})
	return ok
}

func AddHotel(env *cayonlib.Env, hotelId string, cap int32) {
	cayonlib.Write(env, data.Thotel(), hotelId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(Hotel{
			HotelId:   hotelId,
			Cap:       cap,
			Customers: []string{},
		}),
	})
}
