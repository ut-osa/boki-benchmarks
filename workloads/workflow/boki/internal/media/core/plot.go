package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/cayonlib"
)

func WritePlot(env *cayonlib.Env, plotId string, plot string) {
	cayonlib.Write(env, TPlot(), plotId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(aws.JSONValue{"plotId": plotId, "plot": plot}),
	})
}

func ReadPlot(env *cayonlib.Env, plotId string) string {
	item := cayonlib.Read(env, TPlot(), plotId)
	return item.(map[string]interface{})["plot"].(string)
}
