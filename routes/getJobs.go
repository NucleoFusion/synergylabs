package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"assn.com/jwt"
	"assn.com/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetJobsHandle struct {
	Collection *mongo.Collection
}

func (c *GetJobsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query().Get("jwt")

	_, err := jwt.ParseJWT(token)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	result, err := c.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		RaiseError(&w, err)
		return
	}

	jobArr, err := getJobsFromCursor(result)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	data, _ := json.Marshal(jobArr)

	io.Writer.Write(w, data)
}

func getJobsFromCursor(res *mongo.Cursor) ([]models.JobArrResp, error) {
	arr := []models.JobArrResp{}

	defer res.Close(context.Background())

	err := res.All(context.TODO(), &arr)
	if err != nil {
		return arr, err
	}

	return arr, nil
}
