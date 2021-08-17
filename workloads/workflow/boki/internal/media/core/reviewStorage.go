package core

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

func StoreReview(env *cayonlib.Env, review Review) {
	cayonlib.Write(env, TReviewStorage(), review.ReviewId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(review),
	})
}

func ReadReviews(env *cayonlib.Env, ids []string) []Review {
	var reviews []Review
	for _, id := range ids {
		var review Review
		res := cayonlib.Read(env, TReviewStorage(), id)
		cayonlib.CHECK(mapstructure.Decode(res, &review))
		reviews = append(reviews, review)
	}
	return reviews
}
