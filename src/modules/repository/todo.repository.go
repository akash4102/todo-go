package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo/src/modules/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(uri, dbName, collectionName string) (*MongoRepo, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoRepo{collection: collection}, nil
}

func (repo *MongoRepo) GetAll() ([]*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []*models.Todo
	for cursor.Next(ctx) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, cursor.Err()
}

func (repo *MongoRepo) Create(todo *models.Todo) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate a new ObjectID for the Todo
	todo.ID = primitive.NewObjectID()
	_, err := repo.collection.InsertOne(ctx, todo)
	return todo, err
}

func (repo *MongoRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete todo: %v", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("todo not found")
	}
	return nil
}
func (repo *MongoRepo) GetByID(id string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var todo models.Todo
	if err := repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("todo not found")
		}
		return nil, fmt.Errorf("failed to retrieve todo: %v", err)
	}
	return &todo, nil
}

func (repo *MongoRepo) Update(id string, updatedTodo *models.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	update := bson.M{"$set": bson.M{
		"title":   updatedTodo.Title,
		"content": updatedTodo.Content,
		"done":    updatedTodo.Done,
	}}

	result, err := repo.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return fmt.Errorf("failed to update todo: %v", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("todo not found")
	}
	if result.ModifiedCount == 0 {
		return errors.New("todo found but nothing was updated")
	}
	return nil
}
