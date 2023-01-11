package mongodb

import (
	"context"
	"fmt"
	rkmongo "github.com/rookie-ninja/rk-db/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	onceBucket sync.Once
	bucket     *gridfs.Bucket
)

var (
	onceDb sync.Once
	db     *mongo.Database
)

func GetDb() *mongo.Database {
	//TODO : Externalize
	db := rkmongo.GetMongoDB("bee-mongo", "bee")

	if db == nil {
		fmt.Println("Mongo DB not definied in configuration")

		panic("Mongo DB not definied in configuration")

		return nil
	}

	return db
}

func GetDbBucket(name string) *gridfs.Bucket {
	onceBucket.Do(func() { // <-- atomic, does not allow repeating
		db := GetDb()
		bucket = createBucket(db, name)
	})

	return bucket
}

func createBucket(db *mongo.Database, name string) *gridfs.Bucket {
	opts := &options.BucketOptions{
		Name: &name,
	}
	bucket, err := gridfs.NewBucket(db, opts)
	if err != nil {
		fmt.Println("bucket exists may be, continue")

		panic(err)

		return nil
	}

	return bucket
}

func GetDbCollection(ctx context.Context, name string) *mongo.Collection {
	onceDb.Do(func() { // <-- atomic, does not allow repeating
		db := GetDb()
		createCollection(ctx, db, name)
	})

	return db.Collection(name)
}

func createCollection(ctx context.Context, db *mongo.Database, name string) {
	opts := options.CreateCollection()
	err := db.CreateCollection(ctx, name, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
}
