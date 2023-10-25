package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"american-state-backend/graph/model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CNX = Connection()
func Connection() *mongo.Client {
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)

    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)

    if err != nil {
	    log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")

    return client
}

// QueryValidStates is the resolver for the queryValidStates field.
func (r *queryResolver) QueryValidStates(ctx context.Context, keyword *string) ([]string, error) {
	collection := CNX.Database("mydb").Collection("states")
	var results []string 

    // filter out all names that contains keyword
	filter := bson.M{
        "name": bson.M{
            "$regex":   *keyword,
            "$options": "i",
        },
    }

    // query names from database
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    // handle the data to reformat to graphQL response format
    for cursor.Next(ctx) {
        var state model.State
        if err := cursor.Decode(&state); err != nil {
            return nil, err
        }
        results = append(results, state.Name)
    }
    fmt.Println(results)
	if err := cursor.Err(); err != nil {
        return nil, err
    }

	return results, nil
}

// GetStateInfo is the resolver for the getStateInfo field.
func (r *queryResolver) GetStateInfo(ctx context.Context, keyword string) (*model.State, error) {
	collection := CNX.Database("mydb").Collection("states")

    // filter to query the detailed information of the specific state
	filter := bson.M{
        "name": keyword,
    }

    var state model.State
    err := collection.FindOne(ctx, filter).Decode(&state)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }

    return &state, nil


}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
