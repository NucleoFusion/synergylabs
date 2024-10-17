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

type GetJobByIdHandle struct {
	Collection *mongo.Collection
}

func (c *GetJobByIdHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	job := models.JobStruct{}
	err = c.Collection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&job)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	data, _ := json.Marshal(job)

	io.Writer.Write(w, data)
}
