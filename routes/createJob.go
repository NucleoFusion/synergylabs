package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"assn.com/jwt"
	"assn.com/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobCreateHandle struct {
	Collection *mongo.Collection
}

func (c *JobCreateHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	job, err := DecodeToJob(params)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	insertedJob, err := c.Collection.InsertOne(context.Background(), job)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	temp := insertedJob.InsertedID.(primitive.ObjectID)

	jobResp := models.JobStructResp{Id: temp, Job: *job}

	data, _ := json.Marshal(jobResp)

	io.Writer.Write(w, data)
}

func DecodeToJob(params url.Values) (*models.JobStruct, error) {
	title := params.Get("title")
	description := params.Get("description")
	totalApplications := 0
	companyName := params.Get("companyName")
	postedOn := time.Now()

	if title == "" || description == "" || companyName == "" {
		return &models.JobStruct{}, errors.New("invalid params")
	}

	job := models.JobStruct{Title: title, Description: description, TotalApplications: totalApplications, CompanyName: companyName, PostedOn: postedOn}

	return &job, nil
}
