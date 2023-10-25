package main

import (
	"american-state-backend/graph"
	"american-state-backend/graph/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

type States struct {
    State []*State `json:"state"`
}

type Point struct {
	Lat string `json:"_lat"`
	Lng string `json:"_lng"`
}

type State struct {
	Name   string   `json:"_name"`
	Colour string   `json:"_colour"`
	Points []*Point `json:"point"`
}


// importStateData will import states data from state.json to mongoDB database,
// which will only be executed at the first time of project initialization.
func importStateData() {
	jsonFile, err := os.Open("state.json")
    if err != nil {
        fmt.Println("Error opening JSON file:", err)
        return
    }
    defer jsonFile.Close()

    // read the json data of states
    jsonData, err := io.ReadAll(jsonFile)
    if err != nil {
        fmt.Println("Error reading JSON data:", err)
        return
    }
    fmt.Println("JSON data:", string(jsonData))

	var items States
	err = json.Unmarshal(jsonData, &items)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
    fmt.Println(items)
    // connect with mongoDB database
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        fmt.Println("Error connecting to MongoDB:", err)
        return
    }

    collection := client.Database("mydb").Collection("states")

    // Insert formatted data into Database
	for _, state := range items.State {
        
		_, err := collection.InsertOne(context.TODO(), state)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Initilize Database Successfully")
}



var CNX = Connection()
func Connection() *mongo.Client {
    // Set client options
    // Here I have tried to connect Go backend server with dockerized MongoDB container. I pulled mongodb image and create several 
    // mongoDB containers. I tried to use container IP address/container name/container key to make connections
    // (like "mongodb://6b81121b8818add50c92f964f0c748e402fbd0fc6f2ec69a58763e399d897031:27017", "mongodb://container-name:27017" and so on)
    // but all failed. I think I have not found the correct way to connect to dockerized MongoDB container. Now I am trying to search for
    // the possible method.
    clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
    // clientOptions := options.Client().ApplyURI("mongodb://6b81121b8818add50c92f964f0c748e402fbd0fc6f2ec69a58763e399d897031:27017")

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


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}


	collection := CNX.Database("mydb").Collection("states")

    // use findOne to check whether the database is empty; if empty, we will call importStateData to initialize database;
    // otherwise, we will skip and execute normally
    var state model.State
    filter := bson.M{}
    err := collection.FindOne(context.TODO(), filter).Decode(&state)

    if err == mongo.ErrNoDocuments {
        importStateData()
        fmt.Println("Database is empty; initilizing database")
    } else if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println("Database is not empty")
    }


    // "http://localhost:8081" is not allowed to send request to ""http://localhost:8080" since the 
    // "Cross-Origin Resource Sharing" strategy. So we need to add these code to allow Cross-Origin HTTP Request.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"}, // allowed front end address
		AllowedMethods: []string{"POST"},                 // allowed http method
	})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))


	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
