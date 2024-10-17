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
	"go.mongodb.org/mongo-driver/mongo"
)

type GetApplicants struct {
	Collection *mongo.Collection
}

func (c *GetApplicants) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	result, err := c.Collection.Find(context.Background(), bson.M{"usertype": "applicant"})
	if err != nil {
		RaiseError(&w, err)
		return
	}

	applicantArr, err := getUsersFromCursor(result)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	data, _ := json.Marshal(applicantArr)

	io.Writer.Write(w, data)
}

func getUsersFromCursor(res *mongo.Cursor) ([]models.UserStructResp, error) {
	arr := []models.UserStructResp{}

	defer res.Close(context.Background())

	for res.Next(context.Background()) {
		userResp := models.UserStructResp{}

		err := res.Decode(&userResp)
		if err != nil {
			return arr, err
		}

		// id := res.ID()

		// fmt.Println(id)

		arr = append(arr, userResp)
	}

	return arr, nil
}
