package core

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/cayonlib"
	"github.com/mitchellh/mapstructure"
)

func WriteMovieInfo(env *cayonlib.Env, info MovieInfo) {
	cayonlib.Write(env, TMovieInfo(), info.MovieId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(info),
	})
}

func ReadMovieInfo(env *cayonlib.Env, movieId string) MovieInfo {
	var movieInfo MovieInfo
	item := cayonlib.Read(env, TMovieId(), movieId)
	cayonlib.CHECK(mapstructure.Decode(item, &movieInfo))
	return movieInfo
}

func UpdateRating(env *cayonlib.Env, movieId string, sumUncommittedRating int32, numUncommittedRating int32) {
	var movieInfo MovieInfo
	item := cayonlib.Read(env, TMovieId(), movieId)
	cayonlib.CHECK(mapstructure.Decode(item, &movieInfo))
	movieInfo.AvgRating = (movieInfo.AvgRating*float64(movieInfo.NumRating) + float64(sumUncommittedRating)) / float64(movieInfo.NumRating+numUncommittedRating)
	movieInfo.NumRating += numUncommittedRating
	cayonlib.Write(env, TMovieId(), movieId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(movieInfo),
	})
}
