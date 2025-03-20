package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"diet_bot/internal/entity"
)

const (
	databaseName = "diet_bot"
)

type Client struct {
	client       *mongo.Client
	databaseName string

	db *mongo.Database
}

func NewClient() (*Client, error) {
	// Retrieve MongoDB URL from environment variables
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return nil, fmt.Errorf("MONGO_URL environment variable not set")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:       client,
		databaseName: databaseName,
		db:           client.Database(databaseName),
	}, nil
}

func (c *Client) Close() error {
	return c.client.Disconnect(context.TODO())
}

func (c *Client) UpsertUser(user *entity.User) error {
	collection := c.db.Collection(user.CollectionName())

	filter := bson.M{"_id": user.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "first_name", Value: user.FirstName},
			{Key: "last_name", Value: user.LastName},
			{Key: "username", Value: user.Username},
			{Key: "language_code", Value: user.LanguageCode},
			{Key: "is_bot", Value: user.IsBot},
			{Key: "created_at", Value: user.CreatedAt},
			{Key: "updated_at", Value: user.UpdatedAt},
		}},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	return err
}

func (c *Client) GetUser(id int64) (*entity.User, error) {
	var user *entity.User

	collection := c.db.Collection(user.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(user)

	return user, err
}

func (c *Client) UpsertChat(chat *entity.Chat) error {
	collection := c.db.Collection(chat.CollectionName())

	filter := bson.M{"_id": chat.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: chat.UserID},
			{Key: "state", Value: chat.State},
		}},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	return err
}

func (c *Client) GetChat(id int64) (*entity.Chat, error) {
	chat := &entity.Chat{}

	collection := c.db.Collection(chat.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(chat)

	return chat, err
}

func (c *Client) UpsertDietConfiguration(dietConfiguration *entity.DietConfiguration) error {
	collection := c.db.Collection(dietConfiguration.CollectionName())

	if dietConfiguration.ID == uuid.Nil {
		dietConfiguration.ID = uuid.New()
	}

	filter := bson.M{"user_id": dietConfiguration.UserID}
	update := bson.M{"$set": dietConfiguration}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	return err
}

func (c *Client) GetDietConfiguration(userID int64) (*entity.DietConfiguration, error) {
	dietConfiguration := &entity.DietConfiguration{}

	collection := c.db.Collection(dietConfiguration.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(dietConfiguration)

	return dietConfiguration, err
}

func (c *Client) CreateDiet(diet *entity.GeneratedDiet) error {
	collection := c.db.Collection(diet.CollectionName())

	_, err := collection.InsertOne(context.TODO(), diet)

	return err
}

func (c *Client) GetDiet(userID int64) (*entity.GeneratedDiet, error) {
	diet := &entity.GeneratedDiet{}

	collection := c.db.Collection(diet.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(diet)

	return diet, err
}
