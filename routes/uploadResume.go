package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"assn.com/jwt"
	"assn.com/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UploadResume struct {
	Collection *mongo.Collection
}

func (c *UploadResume) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	godotenv.Load(".env")

	token := r.URL.Query().Get("jwt")

	payload, err := jwt.ParseJWT(token)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	user, err := FindUserByPayload(&payload, c.Collection)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("https://api.apilayer.com/resume_parser/upload?apikey=%v", os.Getenv("APIKEY")), "application/octet-stream", r.Body)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	fmt.Println(buf.String())

	profile := models.ProfileStruct{Education: []models.EducationStruct{}, Experience: []models.ExperienceStruct{}}
	err = json.Unmarshal([]byte(buf.String()), &profile)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	user.Profile = profile

	_, err = c.Collection.ReplaceOne(context.Background(), bson.M{"email": user.Email}, user)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	data, _ := json.Marshal(profile)
	io.Writer.Write(w, data)
}

func FindUserByPayload(payload *jwt.Payload, Coll *mongo.Collection) (*models.UserStruct, error) {
	user := models.UserStruct{}
	err := Coll.FindOne(context.Background(), bson.M{
		"email": payload.Email,
	}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
