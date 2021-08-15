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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostListInput struct {
	UserId string `json:"userId,omitempty"`
	Skip   int    `json:"skip,omitempty"`
}

type PostListOutput struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Posts   []interface{} `json:"posts,omitempty"`
}

type postListHandler struct {
	kind   string
	env    types.Environment
	client *mongo.Client
}

func NewSlibPostListHandler(env types.Environment) types.FuncHandler {
	return &postListHandler{
		kind: "slib",
		env:  env,
	}
}

func NewMongoPostListHandler(env types.Environment) types.FuncHandler {
	return &postListHandler{
		kind:   "mongo",
		env:    env,
		client: utils.CreateMongoClientOrDie(context.TODO()),
	}
}

const kMaxReturnPosts = 8

func postListSlib(ctx context.Context, env types.Environment, input *PostListInput) (*PostListOutput, error) {
	txn, err := statestore.CreateReadOnlyTxnEnv(ctx, env)
	if err != nil {
		return nil, err
	}

	var postList []interface{}

	if input.UserId == "" {
		timelineObj := txn.Object("timeline")
		if value, _ := timelineObj.Get("posts"); !value.IsNull() {
			postList = value.AsArray()
		} else {
			postList = make([]interface{}, 0)
		}
	} else {
		userObj := txn.Object(fmt.Sprintf("userid:%s", input.UserId))
		if value, _ := userObj.Get("posts"); !value.IsNull() {
			postList = value.AsArray()
		} else {
			return &PostListOutput{
				Success: false,
				Message: fmt.Sprintf("Cannot find user with ID %s", input.UserId),
			}, nil
		}
	}

	output := &PostListOutput{
		Success: true,
		Posts:   make([]interface{}, 0),
	}

	if input.Skip >= len(postList) {
		return output, nil
	}
	postList = postList[0 : len(postList)-input.Skip]

	for i := len(postList) - 1; i >= 0; i-- {
		postId := postList[i].(string)
		postObj := txn.Object(fmt.Sprintf("post:%s", postId))
		post := make(map[string]string)
		if value, _ := postObj.Get("body"); !value.IsNull() {
			post["body"] = value.AsString()
		}
		if value, _ := postObj.Get("userName"); !value.IsNull() {
			post["user"] = value.AsString()
		}
		if len(post) > 0 {
			output.Posts = append(output.Posts, post)
			if len(output.Posts) == kMaxReturnPosts {
				break
			}
		}
	}
	return output, nil
}

func postListMongo(ctx context.Context, client *mongo.Client, input *PostListInput) (*PostListOutput, error) {
	sess, err := client.StartSession(options.Session())
	if err != nil {
		return nil, err
	}
	defer sess.EndSession(ctx)

	db := client.Database("retwis")

	posts, err := sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		postColl := db.Collection("posts")
		usersColl := db.Collection("users")
		posts := make([]interface{}, 0, kMaxReturnPosts)

		if input.UserId == "" {
			opts := options.Find()
			opts.SetSort(bson.D{{"_id", -1}})
			opts.SetSkip(int64(input.Skip))
			opts.SetLimit(kMaxReturnPosts)
			cursor, err := postColl.Find(sessCtx, bson.D{}, opts)
			if err != nil {
				return nil, err
			}
			var results []bson.M
			err = cursor.All(sessCtx, &results)
			if err != nil {
				return nil, err
			}
			for _, post := range results {
				posts = append(posts, map[string]string{
					"body": post["body"].(string),
					"user": post["userName"].(string),
				})
				if len(posts) == kMaxReturnPosts {
					break
				}
			}
		} else {
			var user bson.M
			if err := usersColl.FindOne(sessCtx, bson.D{{"userId", input.UserId}}).Decode(&user); err != nil {
				return nil, err
			}
			elements := user["posts"].(bson.A)
			if len(elements) > input.Skip {
				end := len(elements) - input.Skip
				for i := end - 1; i >= 0; i-- {
					postId := elements[i]
					var post bson.M
					err := postColl.FindOne(sessCtx, bson.D{{"_id", postId}}).Decode(&post)
					if err != nil {
						return nil, err
					}
					posts = append(posts, map[string]string{
						"body": post["body"].(string),
						"user": post["userName"].(string),
					})
					if len(posts) == kMaxReturnPosts {
						break
					}
				}
			}
		}

		return posts, nil
	}, utils.MongoTxnOptions())

	if err != nil {
		return &PostListOutput{
			Success: false,
			Message: fmt.Sprintf("Mongo failed: %v", err),
		}, nil
	}

	return &PostListOutput{
		Success: true,
		Posts:   posts.([]interface{}),
	}, nil
}

func (h *postListHandler) onRequest(ctx context.Context, input *PostListInput) (*PostListOutput, error) {
	switch h.kind {
	case "slib":
		return postListSlib(ctx, h.env, input)
	case "mongo":
		return postListMongo(ctx, h.client, input)
	default:
		panic(fmt.Sprintf("Unknown kind: %s", h.kind))
	}
}

func (h *postListHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &PostListInput{}
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
