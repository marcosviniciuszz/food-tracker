package repositories

import (
	"context"
	"fmt"
	"food-tracker/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() (*OrderRepository, error) {
	collection := database.GetCollection("food-tracker", "orders")
	if collection == nil {
		return nil, fmt.Errorf("failed to get MongoDB collection")
	}
	return &OrderRepository{collection: collection}, nil
}

func (r *OrderRepository) Insert(ctx context.Context, data bson.M) error {
	// Always insert as "peding"
	data["status"] = "pending"

	_, err := r.collection.InsertOne(ctx, data)
	return err
}

func (repo *OrderRepository) GetOrders(ctx context.Context) ([]bson.M, error) {
	var orders []bson.M

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (repo *OrderRepository) ConfirmOrder(ctx context.Context, id string) error {
	// Convert to ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"status": "confirmed"}}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", id)
	}

	return nil
}

func (repo *OrderRepository) StartPreparation(ctx context.Context, id string) error {
	// Convert to ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"status": "preparation"}}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", id)
	}

	return nil
}

func (repo *OrderRepository) ReadyToPickup(ctx context.Context, id string) error {
	// Convert to ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"status": "ready"}}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", id)
	}

	return nil
}
