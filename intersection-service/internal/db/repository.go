package db

import (
	"context"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	logger := util.LoggerFromContext(ctx)

	logger.Debug("inserting intersection")

	_, err := r.collection.InsertOne(ctx, intersection)
	if err != nil {
		return nil, errs.NewDatabaseError(
			"failed to insert intersection into collection",
			err,
			map[string]any{"intersection ID:": intersection.ID},
		)
	}

	return intersection, nil
}

func (r *MongoIntersectionRepo) GetIntersectionByID(
	ctx context.Context,
	id string,
) (*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)

	logger.Debug("finding intersection by ID")

	var intersection model.Intersection

	filter := bson.M{"id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&intersection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError(
				"intersection ID not found in collection",
				map[string]any{"intersection ID": intersection.ID},
			)
		}
		return nil, errs.NewDatabaseError(
			"failed to find intersection",
			err,
			map[string]any{"intersection ID": intersection.ID},
		)
	}

	return &intersection, nil
}

func (r *MongoIntersectionRepo) GetAllIntersections(
	ctx context.Context, limit, offset int, filter string,
) ([]*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Debug("fetching all intersections")
	// TODO: Implement advanced filtering

	var query bson.M
	if filter != "" {
		ids := strings.Split(strings.TrimSpace(filter), ",")
		var cleanIDs []string
		for _, id := range ids {
			if trimmedID := strings.TrimSpace(id); trimmedID != "" {
				cleanIDs = append(cleanIDs, trimmedID)
			}
		}

		if len(cleanIDs) > 0 {
			query = bson.M{"id": bson.M{"$in": cleanIDs}}
		} else {
			query = bson.M{}
		}
	} else {
		query = bson.M{}
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, errs.NewDatabaseError(
			"failed to find intersections",
			err,
			map[string]any{"limit": limit, "offset": offset, "filter": filter},
		)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logger.Warn("failed to close cursor:", "error", err)
		}
	}()

	var intersections []*model.Intersection
	if err = cursor.All(ctx, &intersections); err != nil {
		return nil, errs.NewDatabaseError(
			"failed to decode intersections",
			err,
			map[string]any{"limit": limit, "offset": offset, "filter": filter},
		)
	}

	return intersections, nil
}

func (r *MongoIntersectionRepo) UpdateIntersection(
	ctx context.Context,
	id string,
	name string,
	details model.IntersectionDetails,
) (*model.Intersection, error) {
	logger := util.LoggerFromContext(ctx)
	logger.Debug("updating intersection")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"name":    name,
			"details": details,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedIntersection model.Intersection

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedIntersection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.NewNotFoundError(
				"intersection ID not found for update",
				map[string]any{"intersection ID": id},
			)
		}
		return nil, errs.NewDatabaseError(
			"failed to update intersection",
			err,
			map[string]any{"intersection ID": id},
		)
	}

	return &updatedIntersection, nil
}

func (r *MongoIntersectionRepo) DeleteIntersection(ctx context.Context, id string) error {
	logger := util.LoggerFromContext(ctx)
	logger.Debug("deleting intersection")

	filter := bson.M{"id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errs.NewDatabaseError(
			"failed to delete intersection",
			err,
			map[string]any{"intersection ID": id},
		)
	}

	if result.DeletedCount == 0 {
		return errs.NewNotFoundError(
			"intersection ID not found for deletion",
			map[string]any{"intersection ID": id},
		)
	}

	return nil
}

func (r *MongoIntersectionRepo) UpdateCurrentParams(
	ctx context.Context,
	id string,
	params model.OptimisationParameters,
) error {
	logger := util.LoggerFromContext(ctx)
	logger.Debug("updating current parameters")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"current_parameters": params,
			"last_run_at":        time.Now(),
			"status":             model.Optimised,
		},
		"$inc": bson.M{
			"run_count": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.NewDatabaseError(
			"failed to update intersection current parameters",
			err,
			map[string]any{"intersection ID": id},
		)
	}

	if result.MatchedCount == 0 {
		return errs.NewNotFoundError(
			"intersection ID not found for optimisation update",
			map[string]any{"intersection ID": id},
		)
	}

	return nil
}

func (r *MongoIntersectionRepo) UpdateBestParams(
	ctx context.Context,
	id string,
	params model.OptimisationParameters,
) error {
	logger := util.LoggerFromContext(ctx)
	logger.Debug("updating best parameters")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"bestparameters": params,
			"lastrunat":      time.Now(),
			"status":         model.Optimised,
		},
		"$inc": bson.M{
			"runcount": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.NewDatabaseError(
			"failed to update intersection best parameters",
			err,
			map[string]any{"intersection ID": id},
		)
	}

	if result.MatchedCount == 0 {
		return errs.NewNotFoundError(
			"intersection ID not found for optimisation update",
			map[string]any{"intersection ID": id},
		)
	}

	return nil
}
