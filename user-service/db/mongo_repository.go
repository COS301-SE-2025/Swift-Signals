package db

import (
	"context"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) UserRepository {
	return &MongoUserRepository{collection: collection}
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Note this was not properly implemented it is the implementation of GetUserByID
func (r *MongoUserRepository) FindByEmail(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// func (r *MongoUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var user models.User
// 	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }
// func (r *MongoUserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
// 	objID, err := primitive.ObjectIDFromHex(user.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = r.collection.UpdateByID(ctx, objID, bson.M{"$set": user})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }
// func (r *MongoUserRepository) DeleteUser(ctx context.Context, id string) error {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
// 	return err
// }
// func (r *MongoUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
// 	cursor, err := r.collection.Find(ctx, bson.M{}, &mongo.Options{
// 		Limit: int64(limit),
// 		Skip:  int64(offset),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)
//
// 	var users []*models.User
// 	for cursor.Next(ctx) {
// 		var user models.User
// 		if err := cursor.Decode(&user); err != nil {
// 			return nil, err
// 		}
// 		users = append(users, &user)
// 	}
// 	return users, nil
// }
