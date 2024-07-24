package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/cms-user/database"
	"github.com/anhhuy1010/cms-user/helpers/util"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"

	"github.com/anhhuy1010/cms-user/constant"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomersLogin struct {
	CustomerUuid string    `json:"customer_uuid" bson:"customer_uuid"`
	Uuid         string    `json:"uuid" bson:"uuid"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	Token        string    `json:"token"`
}

func (u *CustomersLogin) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("customersLogin")
}

func (u *CustomersLogin) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*CustomersLogin, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var customersLogin []*CustomersLogin
	for cursor.Next(context.TODO()) {
		var elem CustomersLogin
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		customersLogin = append(customersLogin, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return customersLogin, nil
}

func (u *CustomersLogin) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*CustomersLogin, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var customersLogin []*CustomersLogin
	for cursor.Next(context.TODO()) {
		var elem CustomersLogin
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println("[Decode] PopularCuisine:", err)
			log.Println("-> #", elem.Uuid)
			continue
		}

		customersLogin = append(customersLogin, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return customersLogin, nil
}

func (u *CustomersLogin) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (u *CustomersLogin) FindOne(conditions map[string]interface{}) (*CustomersLogin, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *CustomersLogin) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (u *CustomersLogin) InsertMany(Customers []interface{}) ([]interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertMany(context.TODO(), Customers)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

func (u *CustomersLogin) Update() (int64, error) {
	coll := u.Model()

	condition := make(map[string]interface{})
	condition["uuid"] = u.Uuid

	u.UpdatedAt = util.GetNowUTC()
	updateStr := make(map[string]interface{})
	updateStr["$set"] = u

	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CustomersLogin) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CustomersLogin) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := u.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CustomersLogin) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := u.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
