package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"cs.utexas.edu/zjia/faas-retwis/utils"

	"cs.utexas.edu/zjia/faas/slib/statestore"
	"cs.utexas.edu/zjia/faas/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileInput struct {
	UserId string `json:"userId"`
}

type ProfileOutput struct {
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	UserName     string `json:"username,omitempty"`
	NumFollowers int    `json:"numFollowers"`
	NumFollowees int    `json:"numFollowees"`
	NumPosts     int    `json:"numPosts"`
}

type profileHandler struct {
	kind   string
	env    types.Environment
	client *mongo.Client
}

func NewSlibProfileHandler(env types.Environment) types.FuncHandler {
	return &profileHandler{
		kind: "slib",
		env:  env,
	}
}

func NewMongoProfileHandler(env types.Environment) types.FuncHandler {
	return &profileHandler{
		kind:   "mongo",
		env:    env,
		client: utils.CreateMongoClientOrDie(context.TODO()),
	}
}

func profileSlib(ctx context.Context, env types.Environment, input *ProfileInput) (*ProfileOutput, error) {
	output := &ProfileOutput{Success: true}

	store := statestore.CreateEnv(ctx, env)
	userObj := store.Object(fmt.Sprintf("userid:%s", input.UserId))
	if value, _ := userObj.Get("username"); !value.IsNull() {
		output.UserName = value.AsString()
	} else {
		return &ProfileOutput{
			Success: false,
			Message: fmt.Sprintf("Cannot find user with ID %s", input.UserId),
		}, nil
	}
	if value, _ := userObj.Get("followers"); !value.IsNull() {
		output.NumFollowers = value.Size()
	}
	if value, _ := userObj.Get("followees"); !value.IsNull() {
		output.NumFollowees = value.Size()
	}
	if value, _ := userObj.Get("posts"); !value.IsNull() {
		output.NumPosts = value.Size()
	}

	return output, nil
}

func profileMongo(ctx context.Context, client *mongo.Client, input *ProfileInput) (*ProfileOutput, error) {
	db := client.Database("retwis")

	var user bson.M
	if err := db.Collection("users").FindOne(ctx, bson.D{{"userId", input.UserId}}).Decode(&user); err != nil {
		return &ProfileOutput{
			Success: false,
			Message: fmt.Sprintf("Mongo failed: %v", err),
		}, nil
	}

	output := &ProfileOutput{Success: true}
	if value, ok := user["username"].(string); ok {
		output.UserName = value
	}
	if value, ok := user["followers"].(bson.M); ok {
		output.NumFollowers = len(value)
	}
	if value, ok := user["followees"].(bson.M); ok {
		output.NumFollowees = len(value)
	}
	if value, ok := user["posts"].(bson.A); ok {
		output.NumPosts = len(value)
	}

	return output, nil
}

func (h *profileHandler) onRequest(ctx context.Context, input *ProfileInput) (*ProfileOutput, error) {
	switch h.kind {
	case "slib":
		return profileSlib(ctx, h.env, input)
	case "mongo":
		return profileMongo(ctx, h.client, input)
	default:
		panic(fmt.Sprintf("Unknown kind: %s", h.kind))
	}
}

func (h *profileHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &ProfileInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := h.onRequest(ctx, parsedInput)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}
