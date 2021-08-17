package core

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

func WriteCastInfo(env *cayonlib.Env, info CastInfo) {
	cayonlib.Write(env, TCastInfo(), info.CastInfoId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(info),
	})
}

func ReadCastInfo(env *cayonlib.Env, castIds []string) []CastInfo {
	var res []CastInfo
	for _, id := range castIds {
		var castInfo CastInfo
		item := cayonlib.Read(env, TCastInfo(), id)
		cayonlib.CHECK(mapstructure.Decode(item, &castInfo))
		res = append(res, castInfo)
	}
	return res
}
