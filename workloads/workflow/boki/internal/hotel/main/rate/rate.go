package rate

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
	"sort"
)

func GetRates(env *cayonlib.Env, req Request) Result {
	var plans RatePlans
	for _, i := range req.HotelIds {
		plan := data.RatePlan{}
		res := cayonlib.Read(env, data.Trate(), i)
		err := mapstructure.Decode(res, &plan)
		cayonlib.CHECK(err)
		if plan.HotelId != "" {
			plans = append(plans, plan)
		}
	}
	sort.Sort(plans)
	return Result{RatePlans: plans}
}
