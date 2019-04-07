package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OpenConnection open connection to spesific mongo database & collection
func OpenConnection(timout int, dburl string, dbname string, dbcoll string) (context.Context, context.CancelFunc, *mongo.Client, *mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timout)*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	collection := client.Database(dbname).Collection(dbcoll)

	if err != nil {
		log.Println(err)
	}

	return ctx, cancel, client, collection, err
}
