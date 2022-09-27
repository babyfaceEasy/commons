package mongo

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	localMongoURL = "mongodb://localhost:27017"
)

// DBConfig configures mongoDB
type DBConfig struct {
	DBURL  string
	DBName string
}

// DBProviderFunc provides the functionality of retuning a mongoDB database
type DBProviderFunc func() *mongo.Database

// CloseMongoFunc provides the functionality of closing mongoDB connection
type CloseMongoFunc func()

// DBProvider retuens a new mongoDB database
func DBProvider(c *mongo.Client, dbname string) DBProviderFunc {
	return func() *mongo.Database {
		return c.Database(dbname)
	}
}

// ToProvider returns a mongoDB provider from the config
func (c DBConfig) ToProvider() (DBProviderFunc, CloseMongoFunc, error) {

	timeout := getTimeout()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	mClient, err := mongo.Connect(ctx, options.Client().ApplyURI(c.DBURL))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Unable to connect to mongo using config=%+v", c)
	}
	return DBProvider(mClient, c.DBName), DisconnectMongo(ctx, mClient), nil
}

// NewConfigFromEnvVar returns mongo configuration from environment variables
func NewConfigFromEnvVar() DBConfig {
	return DBConfig{
		DBURL:  os.Getenv("MONGO_URL"),
		DBName: os.Getenv("MONGO_DB_NAME"),
	}
}

// NewLocalConfig returns mongo configuration for local testing
func NewLocalConfig(dbName string) DBConfig {
	return DBConfig{
		DBURL:  localMongoURL,
		DBName: dbName,
	}
}

// DisconnectMongo disconnects the mongo after a specific timeout embedded in the context
func DisconnectMongo(ctx context.Context, client *mongo.Client) CloseMongoFunc {
	return func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect mongo with error: ", err)
			return
		}
		log.Println("Disconnected...")
	}
}

func getTimeout() int64 {

	t := os.Getenv("MONGO_TIMEOUT")
	if t != "" {
		newTimeout, err := strconv.Atoi(t)
		if err != nil {
			log.Println("Could not convert mongo timeout with error: ", err)
			return int64(10)
		}
		return int64(newTimeout)
	}
	return int64(10)
}
