package core

import "github.com/eniac/Beldi/pkg/cayonlib"

type User struct {
	UserId    string
	FirstName string
	LastName  string
	Username  string
	Password  string
	Salt      string
}

type Review struct {
	ReviewId  string
	UserId    string
	ReqId     string
	Text      string
	MovieId   string
	Rating    int32
	Timestamp string
}

type CastInfo struct {
	CastInfoId string
	Name       string
	Gender     bool
	Intro      string
}

type Cast struct {
	CastId     string
	Character  string
	CastInfoId string
}

type MovieInfo struct {
	MovieId      string
	Title        string
	Casts        []Cast
	PlotId       string
	ThumbnailIds []string
	PhotoIds     []string
	VideoIds     []string
	AvgRating    float64
	NumRating    int32
}

type ReviewInfo struct {
	ReviewId  string
	Timestamp string
}

type Page struct {
	MovieInfo MovieInfo
	Reviews   []Review
	CastInfos []CastInfo
	Plot      string
}

type RPCInput struct {
	Function string
	Input    interface{}
}

func TCastInfo() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bCastInfo"
	} else {
		return "CastInfo"
	}
}

func TComposeReview() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bComposeReview"
	} else {
		return "ComposeReview"
	}
}

func TFrontend() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bFrontend"
	} else {
		return "Frontend"
	}
}

func TMovieId() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bMovieId"
	} else {
		return "MovieId"
	}
}

func TMovieInfo() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bMovieInfo"
	} else {
		return "MovieInfo"
	}
}

func TMovieReview() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bMovieReview"
	} else {
		return "MovieReview"
	}
}

func TPage() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bPage"
	} else {
		return "Page"
	}
}

func TPlot() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bPlot"
	} else {
		return "Plot"
	}
}

func TRating() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bRating"
	} else {
		return "Rating"
	}
}

func TReviewStorage() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bReviewStorage"
	} else {
		return "ReviewStorage"
	}
}

func TText() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bText"
	} else {
		return "Text"
	}
}

func TUniqueId() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bUniqueId"
	} else {
		return "UniqueId"
	}
}

func TUser() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bUser"
	} else {
		return "User"
	}
}

func TUserReview() string {
	if cayonlib.TYPE == "BASELINE" {
		return "bUserReview"
	} else {
		return "UserReview"
	}
}
