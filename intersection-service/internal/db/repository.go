package db

import (
	"context"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoIntersectionRepo struct {
	collection *mongo.Collection
}

func NewMongoIntersectionRepo(collection *mongo.Collection) IntersectionRepository {
	return &MongoIntersectionRepo{collection: collection}
}

func (r *MongoIntersectionRepo) CreateIntersection(
	ctx context.Context,
	intersection *model.Intersection,
) (*model.Intersection, error) {
	_, err := r.collection.InsertOne(ctx, intersection)
	if err != nil {
		return nil, err
	}

	return intersection, nil
}

func (r *MongoIntersectionRepo) GetIntersectionByID(
	ctx context.Context,
	id string,
) (*model.Intersection, error) {
	var intersection model.Intersection

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

func (r *MongoIntersectionRepo) GetAllIntersections(
	ctx context.Context,
) ([]*model.Intersection, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var intersections []*model.Intersection
	for cursor.Next(ctx) {
		var intersection model.Intersection
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

func (r *MongoIntersectionRepo) UpdateIntersection(
	ctx context.Context,
	id string,
	name string,
	details model.IntersectionDetails,
) (*model.Intersection, error) {
	// TODO: Implement UpdateIntersection
	return nil, nil
}

func (r *MongoIntersectionRepo) DeleteIntersection(ctx context.Context, id string) error {
	// TODO: Implement DeleteIntersection
	return nil
}

func (r *MongoIntersectionRepo) PutOptimisation(
	ctx context.Context,
	id string,
	params model.OptimisationParameters,
) (bool, error) {
	// TODO: Implement PutOptimisation
	return true, nil
}
