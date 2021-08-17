package user

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

func CheckUser(env *cayonlib.Env, req Request) Result {
	var user data.User
	item := cayonlib.Read(env, data.Tuser(), req.Username)
	err := mapstructure.Decode(item, &user)
	cayonlib.CHECK(err)
	return Result{Correct: req.Password == user.Password}
}
