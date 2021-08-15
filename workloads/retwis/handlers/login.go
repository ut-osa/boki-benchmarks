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

type LoginInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	UserId  string `json:"userId"`
	Auth    string `json:"auth"`
}

type loginHandler struct {
	kind   string
	env    types.Environment
	client *mongo.Client
}

func NewSlibLoginHandler(env types.Environment) types.FuncHandler {
	return &loginHandler{
		kind: "slib",
		env:  env,
	}
}

func NewMongoLoginHandler(env types.Environment) types.FuncHandler {
	return &loginHandler{
		kind:   "mongo",
		env:    env,
		client: utils.CreateMongoClientOrDie(context.TODO()),
	}
}

func loginSlib(ctx context.Context, env types.Environment, input *LoginInput) (*LoginOutput, error) {
	txn, err := statestore.CreateReadOnlyTxnEnv(ctx, env)
	if err != nil {
		return nil, err
	}

	userId := ""

	userNameObj := txn.Object(fmt.Sprintf("username:%s", input.UserName))
	if value, _ := userNameObj.Get("id"); !value.IsNull() {
		userId = value.AsString()
	} else {
		return &LoginOutput{
			Success: false,
			Message: fmt.Sprintf("User name \"%s\" does not exists", input.UserName),
		}, nil
	}

	userObj := txn.Object(fmt.Sprintf("userid:%s", userId))
	if value, _ := userObj.Get("password"); !value.IsNull() {
		if input.Password != value.AsString() {
			return &LoginOutput{
				Success: false,
				Message: "Incorrect password",
			}, nil
		}
	} else {
		return &LoginOutput{
			Success: false,
			Message: fmt.Sprintf("Cannot find user with ID %s", userId),
		}, nil
	}

	output := &LoginOutput{Success: true, UserId: userId}
	if value, _ := userObj.Get("auth"); !value.IsNull() {
		output.Auth = value.AsString()
	}
	return output, nil
}

func loginMongo(ctx context.Context, client *mongo.Client, input *LoginInput) (*LoginOutput, error) {
	db := client.Database("retwis")

	var user bson.M
	if err := db.Collection("users").FindOne(ctx, bson.D{{"username", input.UserName}}).Decode(&user); err != nil {
		return &LoginOutput{
			Success: false,
			Message: fmt.Sprintf("Mongo failed: %v", err),
		}, nil
	}

	if input.Password != user["password"].(string) {
		return &LoginOutput{
			Success: false,
			Message: "Incorrect password",
		}, nil
	}

	return &LoginOutput{
		Success: true,
		UserId:  user["userId"].(string),
		Auth:    user["auth"].(string),
	}, nil
}

func (h *loginHandler) onRequest(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	switch h.kind {
	case "slib":
		return loginSlib(ctx, h.env, input)
	case "mongo":
		return loginMongo(ctx, h.client, input)
	default:
		panic(fmt.Sprintf("Unknown kind: %s", h.kind))
	}
}

func (h *loginHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &LoginInput{}
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
