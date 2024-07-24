package models

import (
	"context"
	"time"

	"github.com/anhhuy1010/cms-user/database"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"

	"github.com/anhhuy1010/cms-user/constant"
)

type AdminSignUp struct {
	Uuid      string    `json:"uuid" bson:"uuid"`
	Password  string    `json:"password" bson:"password"`
	Username  string    `json:"username" bson:"username"`
	IsActive  int       `json:"is_active" bson:"is_active"`
	IsDelete  int       `json:"is_delete" bson:"is_delete"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (u *AdminSignUp) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("adminSignUp")
}

func (u *AdminSignUp) FindOne(conditions map[string]interface{}) (*AdminSignUp, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *AdminSignUp) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}
