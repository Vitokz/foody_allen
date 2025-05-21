package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"diet_bot/internal/entity"
	internalerrors "diet_bot/internal/entity/errors"
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
	if err != nil {
		return internalerrors.ErrorFailedToSaveUser
	}

	return err
}

func (c *Client) GetUser(id int64) (*entity.User, error) {
	user := &entity.User{}

	collection := c.db.Collection(user.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, internalerrors.ErrorUserNotFound
		}

		return nil, internalerrors.ErrorFailedToGetUser
	}

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
	if err != nil {
		return internalerrors.ErrorFailedToSaveChat
	}

	return err
}

func (c *Client) GetChat(id int64) (*entity.Chat, error) {
	chat := &entity.Chat{}

	collection := c.db.Collection(chat.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(chat)

	return chat, err
}

func (c *Client) CreateDiet(diet *entity.GeneratedDiet) error {
	collection := c.db.Collection(diet.CollectionName())

	_, err := collection.InsertOne(context.TODO(), diet)

	return err
}

func (c *Client) DeleteDiet(userID int64) error {
	diet := &entity.GeneratedDiet{}
	collection := c.db.Collection(diet.CollectionName())

	_, err := collection.DeleteOne(context.TODO(), bson.M{"user_id": userID})

	return err
}

func (c *Client) GetDiet(userID int64) (*entity.GeneratedDiet, error) {
	diet := &entity.GeneratedDiet{}

	collection := c.db.Collection(diet.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(diet)

	return diet, err
}

func (c *Client) GetLatestDiet(userID int64) (*entity.GeneratedDiet, error) {
	diet := &entity.GeneratedDiet{}
	collection := c.db.Collection(diet.CollectionName())

	opts := options.FindOne().SetSort(bson.D{{"created_at", -1}})

	err := collection.FindOne(
		context.TODO(),
		bson.M{"user_id": userID},
		opts,
	).Decode(diet)

	return diet, err
}

func (c *Client) CreateUserConfiguration(userConfiguration *entity.UserConfiguration) error {
	userConfiguration.CreatedAt = time.Now()
	userConfiguration.UpdatedAt = time.Now()
	collection := c.db.Collection(userConfiguration.CollectionName())

	_, err := collection.InsertOne(context.TODO(), userConfiguration)

	return err
}

func (c *Client) SaveUserConfiguration(userConfiguration *entity.UserConfiguration) error {
	userConfiguration.UpdatedAt = time.Now()

	collection := c.db.Collection(userConfiguration.CollectionName())

	filter := bson.M{"user_id": userConfiguration.UserID}
	update := bson.M{"$set": userConfiguration}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	return err
}

func (c *Client) GetUserConfiguration(userID int64) (*entity.UserConfiguration, error) {
	userConfiguration := &entity.UserConfiguration{}

	collection := c.db.Collection(userConfiguration.CollectionName())

	err := collection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(userConfiguration)

	return userConfiguration, err
}
