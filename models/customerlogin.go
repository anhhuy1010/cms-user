package models

import (
	"context"
	"log"
	"time"

	"github.com/anhhuy1010/cms-user/database"
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
