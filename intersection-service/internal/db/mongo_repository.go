package db

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoIntersectionRepository struct {
	collection *mongo.Collection
}

func NewMongoIntersectionRepository(collection *mongo.Collection) IntersectionRepository {
	return &MongoIntersectionRepository{collection: collection}
}

func (r *MongoIntersectionRepository) CreateIntersection(ctx context.Context, intersection *model.IntersectionResponse) (*model.IntersectionResponse, error) {
	_, err := r.collection.InsertOne(ctx, intersection)
	if err != nil {
		return nil, err
	}

	return intersection, nil
}

func (r *MongoIntersectionRepository) GetIntersectionByID(ctx context.Context, id string) (*model.IntersectionResponse, error) {
	var intersection model.IntersectionResponse

	filter := bson.M{"id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&intersection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &intersection, nil
}

func (r *MongoIntersectionRepository) GetAllIntersections(ctx context.Context) ([]*model.IntersectionResponse, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var intersections []*model.IntersectionResponse
	for cursor.Next(ctx) {
		var intersection model.IntersectionResponse
		if err := cursor.Decode(&intersection); err != nil {
			return nil, err
		}
		intersections = append(intersections, &intersection)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return intersections, nil
}
