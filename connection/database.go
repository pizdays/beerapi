package connection

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mongoCollection *mongo.Collection

func NewMongoDB() *mongo.Collection {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:password@localhost:27017"))
	if err != nil {
		fmt.Println(err)

	}

	db := client.Database("test")
	collection := db.Collection("logs")

	mongoCollection = collection
	return mongoCollection
}

var mariaDBInstance *gorm.DB

func NewMariaDB() *gorm.DB {
	dsn := "root:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mariaDBInstance = db
	return mariaDBInstance
}
