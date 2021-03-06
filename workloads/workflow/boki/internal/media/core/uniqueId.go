package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/lithammer/shortuuid"
)

func UploadUniqueId2(env *cayonlib.Env, reqId string) {
	reviewId := shortuuid.New()
	cayonlib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadUniqueId",
		Input:    aws.JSONValue{"reqId": reqId, "reviewId": reviewId},
	})
}
