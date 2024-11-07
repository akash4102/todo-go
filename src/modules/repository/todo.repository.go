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

func (r *MongoRepo) GetTodoMetrics(ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "date", Value: "$created"},
				{Key: "type", Value: "$type"},
				{Key: "completed", Value: "$done"},
			}},
			{Key: "taskCount", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "effortSum", Value: bson.D{{Key: "$sum", Value: "$effortHr"}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.type"},
			{Key: "totalTasks", Value: bson.D{{Key: "$sum", Value: "$taskCount"}}},
			{Key: "completedTasks", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{"$_id.completed", "$taskCount", 0}}}}}},
			{Key: "notCompletedTasks", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{"$_id.completed", 0, "$taskCount"}}}}}},
			{Key: "totalEffort", Value: bson.D{{Key: "$sum", Value: "$effortSum"}}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "completionPercentage", Value: bson.D{{Key: "$multiply", Value: bson.A{
				bson.D{{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{"$totalTasks", 0}}},
					0,
					bson.D{{Key: "$divide", Value: bson.A{"$completedTasks", "$totalTasks"}}},
				}}},
				100,
			}}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
