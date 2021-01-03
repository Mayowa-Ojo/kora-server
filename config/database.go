package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBConn -
type DBConn struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var conn DBConn

// InitDB - create a mongo connection
func InitDB(env *EnvConfig) (*DBConn, error) {
	mongoURI := fmt.Sprintf("%s", env.DBUri)
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.NewClient(clientOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(env.DBName)

	if err != nil {
		return nil, err
	}

	conn := &DBConn{
		Client: client,
		DB:     db,
	}

	log.Println("[x] - connected to database")

	return conn, nil
}
