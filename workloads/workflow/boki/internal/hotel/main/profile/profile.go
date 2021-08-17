package profile

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

func GetProfiles(env *cayonlib.Env, req Request) Result {
	var hotels []data.Hotel
	for _, i := range req.HotelIds {
		hotel := data.Hotel{}
		res := cayonlib.Read(env, data.Tprofile(), i)
		err := mapstructure.Decode(res, &hotel)
		cayonlib.CHECK(err)
		hotels = append(hotels, hotel)
	}
	return Result{Hotels: hotels}
}
