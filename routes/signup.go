package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"assn.com/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandle struct {
	Collection *mongo.Collection
}

func (c *SignupHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := r.URL.Query()

	user, err := DecodeToUser(&params, &w)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	insertErr := insertIntoUsers(c.Collection, *user)
	if insertErr != nil {
		RaiseError(&w, insertErr)
		return
	}

	data, _ := json.Marshal(models.SuccessResp{Message: "successfully created user, login to get your JWT"})

	io.Writer.Write(w, data)

}

func DecodeToUser(params *url.Values, w *http.ResponseWriter) (*models.UserStruct, error) {
	user := models.UserStruct{}

	name := params.Get("name")
	email := params.Get("email")
	userType := params.Get("userType")
	profileHeadline := params.Get("profileHeadline")
	address := params.Get("address")
	password := params.Get("password")

	if name == "" || password == "" || address == "" || profileHeadline == "" || userType == "" || email == "" {
		err := errors.New("missing params")
		return &user, err
	}

	if userType != "admin" && userType != "applicant" {
		err := errors.New("invalid user type")
		return &user, err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 5)

	user.Address = address
	user.Email = email
	user.Name = name
	user.ProfileHeadline = profileHeadline
	user.UserType = userType
	user.Password = string(hashed)

	user.Profile.Name = name
	user.Profile.Email = email

	return &user, nil
}

func RaiseError(w *http.ResponseWriter, err error) {
	data, _ := json.Marshal(models.ErrorResp{Message: err.Error()})

	io.Writer.Write(*w, data)
}

func insertIntoUsers(coll *mongo.Collection, user models.UserStruct) error {

	userFound := models.UserStruct{}
	existsErr := coll.FindOne(context.Background(), bson.M{
		"email": user.Email,
	}).Decode(&userFound)
	if existsErr != mongo.ErrNoDocuments {
		return errors.New("user with same email exists")
	}

	_, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}
