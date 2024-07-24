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

type CustomersSignUp struct {
	Uuid      string    `json:"uuid,omitempty" bson:"uuid"`
	Name      string    `json:"name" bson:"name"`
	Age       int       `json:"age" bson:"age"`
	Password  string    `json:"password" bson:"password"`
	Email     string    `json:"email" bson:"email"`
	Username  string    `json:"username,omitempty" bson:"username"`
	IsActive  int       `json:"is_active" bson:"is_active"`
	IsDelete  int       `json:"is_delete" bson:"is_delete"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (u *CustomersSignUp) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("customersSignUp")
}

func (u *CustomersSignUp) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*CustomersSignUp, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var customersSignUp []*CustomersSignUp
	for cursor.Next(context.TODO()) {
		var elem CustomersSignUp
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		customersSignUp = append(customersSignUp, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return customersSignUp, nil
}

func (u *CustomersSignUp) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*CustomersSignUp, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var customers []*CustomersSignUp
	for cursor.Next(context.TODO()) {
		var elem CustomersSignUp
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println("[Decode] PopularCuisine:", err)
			log.Println("-> #", elem.Uuid)
			continue
		}

		customers = append(customers, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return customers, nil
}

func (u *CustomersSignUp) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (u *CustomersSignUp) FindOne(conditions map[string]interface{}) (*CustomersSignUp, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *CustomersSignUp) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (u *CustomersSignUp) InsertMany(Customers []interface{}) ([]interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertMany(context.TODO(), Customers)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

func (u *CustomersSignUp) Update() (int64, error) {
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

func (u *CustomersSignUp) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CustomersSignUp) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := u.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CustomersSignUp) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := u.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
