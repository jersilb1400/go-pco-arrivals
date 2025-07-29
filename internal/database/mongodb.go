package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB represents the MongoDB database connection
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB() (*MongoDB, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is required")
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "pco_arrivals"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(databaseName)

	return &MongoDB{
		client:   client,
		database: database,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// Migrate creates the necessary collections and indexes
func (m *MongoDB) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create collections and indexes
	collections := []string{"users", "sessions", "events", "locations", "check_ins", "notifications", "billboard_states", "security_codes"}

	for _, collectionName := range collections {
		collection := m.database.Collection(collectionName)

		// Create collection if it doesn't exist
		if err := m.database.CreateCollection(ctx, collectionName); err != nil {
			// Collection might already exist, which is fine
			log.Printf("Collection %s might already exist: %v", collectionName, err)
		}

		// Create indexes based on collection
		switch collectionName {
		case "users":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "pco_user_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Printf("Failed to create index on users.pco_user_id: %v", err)
			}

		case "sessions":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "token", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Printf("Failed to create index on sessions.token: %v", err)
			}

		case "events":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "pco_event_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Printf("Failed to create index on events.pco_event_id: %v", err)
			}

		case "locations":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "pco_location_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Printf("Failed to create index on locations.pco_location_id: %v", err)
			}

		case "check_ins":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "pco_check_in_id", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Printf("Failed to create index on check_ins.pco_check_in_id: %v", err)
			}

		case "notifications":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "pco_check_in_id", Value: 1}},
			})
			if err != nil {
				log.Printf("Failed to create index on notifications.pco_check_in_id: %v", err)
			}

		case "security_codes":
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "code", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			if err != nil {
				log.Printf("Failed to create index on security_codes.code: %v", err)
			}
		}
	}

	log.Println("MongoDB migration completed successfully")
	return nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// InsertOne inserts a single document
func (m *MongoDB) InsertOne(collectionName string, document interface{}) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// FindOne finds a single document
func (m *MongoDB) FindOne(collectionName string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	return collection.FindOne(ctx, filter).Decode(result)
}

// Find finds multiple documents
func (m *MongoDB) Find(collectionName string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

// UpdateOne updates a single document
func (m *MongoDB) UpdateOne(collectionName string, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteOne deletes a single document
func (m *MongoDB) DeleteOne(collectionName string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// CountDocuments counts documents matching a filter
func (m *MongoDB) CountDocuments(collectionName string, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	return collection.CountDocuments(ctx, filter)
}

// Aggregate performs an aggregation pipeline
func (m *MongoDB) Aggregate(collectionName string, pipeline interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.database.Collection(collectionName)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}
