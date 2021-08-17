package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/cayonlib"
)

func UploadText2(env *cayonlib.Env, reqId string, text string) {
	cayonlib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadText",
		Input:    aws.JSONValue{"reqId": reqId, "text": text},
	})
}
