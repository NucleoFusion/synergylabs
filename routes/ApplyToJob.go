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

type ApplyToJob struct {
	Collection *mongo.Collection
}

func (c *ApplyToJob) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := r.URL.Query()

	_, err := jwt.ParseJWT(params.Get("jwt"))
	if err != nil {
		RaiseError(&w, err)
		return
	}

	id := params.Get("job_id")
	if id == "" {
		RaiseError(&w, errors.New("invalid params"))
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		RaiseError(&w, errors.New("id not provided or invalid id"))
		return
	}

	job := models.JobArrResp{}
	err = c.Collection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&job)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	job.TotalApplications += 1

	_, err = c.Collection.UpdateOne(context.Background(), bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"totalapplications", job.TotalApplications}}}})
	if err != nil {
		RaiseError(&w, err)
		return
	}

	data, _ := json.Marshal(job)

	io.Writer.Write(w, data)
}
