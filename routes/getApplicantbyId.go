package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"assn.com/jwt"
	"assn.com/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetApplicantById struct {
	Collection *mongo.Collection
}

func (c *GetApplicantById) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := r.URL.Query()

	payload, err := jwt.ParseJWT(params.Get("jwt"))
	if err != nil {
		RaiseError(&w, err)
		return
	}

	if payload.UserType != "admin" {
		RaiseError(&w, errors.New("user does not have authentication"))
		return
	}

	id := r.PathValue("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if id == "" || err != nil {
		RaiseError(&w, errors.New("id not provided or invalid id"))
		return
	}

	user := models.UserStruct{}
	err = c.Collection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	data, _ := json.Marshal(user)

	io.Writer.Write(w, data)
}
