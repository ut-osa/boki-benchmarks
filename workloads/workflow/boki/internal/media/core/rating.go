package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/cayonlib"
)

func UploadRating2(env *cayonlib.Env, reqId string, rating int32) {
	cayonlib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadRating",
		Input: aws.JSONValue{
			"reqId":  reqId,
			"rating": rating,
		},
	})
}
