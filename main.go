package main

import (
	"context"
	"fmt"
	"go-backend-rinha/datastore"
	"go-backend-rinha/handlers"
	Util "go-backend-rinha/util"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	mongoDBClient := setupMongoDBClient(ctx)

	mongoPessoaClient := datastore.NewPessoaClient(mongoDBClient)
	mongoPessoaClient.InitPessoas(ctx)

	router := gin.Default()

	SERVER_PORT := Util.GetEnvVariable("SERVER_PORT")

	if err := handlers.Setup(ctx, router, mongoPessoaClient).Run(":" + SERVER_PORT); err != nil {
		return
	}

}

func setupMongoDBClient(ctx context.Context) *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/test?authSource=admin",
		Util.GetEnvVariable("DB_USER"),
		Util.GetEnvVariable("DB_PWD"),
		Util.GetEnvVariable("DB_HOST"))

	log.Default().Println(uri)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(fmt.Errorf("could not connect to databse: %w", err))
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(fmt.Errorf("could not ping databse: %w", err))
	}

	log.Println("Connected to MongoDB")

	return client
}
