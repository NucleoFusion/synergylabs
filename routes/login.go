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
	"golang.org/x/crypto/bcrypt"
)

type LoginHandle struct {
	Collection *mongo.Collection
}

func (c *LoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	password := r.URL.Query().Get("password")
	email := r.URL.Query().Get("email")
	if email == "" || password == "" {
		RaiseError(&w, errors.New("invalid params"))
		return
	}

	user, err := matchPass(c.Collection, password, email)
	if err != nil {
		RaiseError(&w, err)
		return
	}

	JWT := jwt.CreateJWT(user)

	data, _ := json.Marshal(jwt.JWTResp{JWT: JWT})

	io.Writer.Write(w, data)
}

func matchPass(coll *mongo.Collection, pass string, email string) (*models.UserStruct, error) {
	userFound := models.UserStruct{}
	err := coll.FindOne(context.Background(), bson.M{
		"email": email,
	}).Decode(&userFound)
	if err == mongo.ErrNoDocuments {
		return &userFound, errors.New("no user with same email exists")
	} else if err != nil {
		return &userFound, err
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(pass))
	if passErr != nil {
		return &userFound, errors.New("password not matched")
	}

	return &userFound, nil
}
