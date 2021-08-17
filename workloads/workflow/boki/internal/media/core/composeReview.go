package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
	"sync"
)

func UploadReq(env *cayonlib.Env, reqId string) {
	cayonlib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(aws.JSONValue{"reqId": reqId, "counter": 0}),
	})
}

func UploadUniqueId(env *cayonlib.Env, reqId string, reviewId string) {
	cayonlib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.reviewId"): expression.Value(reviewId),
		expression.Name("V.counter"):  expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadText(env *cayonlib.Env, reqId string, text string) {
	cayonlib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.text"):    expression.Value(text),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadRating(env *cayonlib.Env, reqId string, rating int32) {
	cayonlib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.rating"):  expression.Value(rating),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadUserId(env *cayonlib.Env, reqId string, userId string) {
	cayonlib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.userId"):  expression.Value(userId),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadMovieId(env *cayonlib.Env, reqId string, movieId string) {
	cayonlib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.movieId"): expression.Value(movieId),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func Cleanup(reqId string) {
	// Debugging
/*
	if cayonlib.TYPE == "BASELINE" {
		cayonlib.LibDelete(TComposeReview(), aws.JSONValue{"K": reqId})
		return
	}
	cayonlib.LibDelete(TComposeReview(), aws.JSONValue{
		"K":       reqId,
		"ROWHASH": "HEAD",
	})
*/
	//cond := expression.Key("K").Equal(expression.Value(reqId))
	//expr, err := expression.NewBuilder().
	//	WithProjection(cayonlib.BuildProjection([]string{"K", "ROWHASH"})).
	//	WithKeyCondition(cond).Build()
	//cayonlib.CHECK(err)
	//res, err := cayonlib.DBClient.Query(&dynamodb.QueryInput{
	//	TableName:                 aws.String(TComposeReview()),
	//	KeyConditionExpression:    expr.KeyCondition(),
	//	ProjectionExpression:      expr.Projection(),
	//	ExpressionAttributeNames:  expr.Names(),
	//	ExpressionAttributeValues: expr.Values(),
	//	ConsistentRead:            aws.Bool(true),
	//})
	//cayonlib.CHECK(err)
	//var items []aws.JSONValue
	//err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &items)
	//cayonlib.CHECK(err)
	//for _, item := range items {
	//	cayonlib.LibDelete(TComposeReview(), item)
	//}
}

func TryComposeAndUpload(env *cayonlib.Env, reqId string) {
	item := cayonlib.Read(env, TComposeReview(), reqId)
	if item == nil {
		return
	}
	res := item.(map[string]interface{})
	if counter, ok := res["counter"].(float64); ok {
		if int32(counter) == 5 {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				Cleanup(reqId)
			}()
			var review Review
			cayonlib.CHECK(mapstructure.Decode(res, &review))
			cayonlib.AsyncInvoke(env, TReviewStorage(), RPCInput{
				Function: "StoreReview",
				Input:    review,
			})
			cayonlib.AsyncInvoke(env, TUserReview(), RPCInput{
				Function: "UploadUserReview",
				Input: aws.JSONValue{
					"userId":    review.UserId,
					"reviewId":  review.ReviewId,
					"timestamp": review.Timestamp,
				},
			})
			cayonlib.AsyncInvoke(env, TMovieReview(), RPCInput{
				Function: "UploadMovieReview",
				Input: aws.JSONValue{
					"movieId":   review.MovieId,
					"reviewId":  review.ReviewId,
					"timestamp": review.Timestamp,
				},
			})
			wg.Wait()
		}
	} else {
		panic("counter not found")
	}
}
