package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileStruct struct {
	Skills     []string           `json:"skills"`
	Education  []EducationStruct  `json:"education"`
	Experience []ExperienceStruct `json:"experience"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
}

type UserStruct struct {
	Name            string        `json:"name"`
	Email           string        `json:"email"`
	Address         string        `json:"address"`
	UserType        string        `json:"userType"`
	Password        string        `json:"password"`
	ProfileHeadline string        `json:"profileHeading"`
	Profile         ProfileStruct `json:"profile"`
}

type JobStruct struct {
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	PostedOn          time.Time `json:"postedOn"`
	TotalApplications int       `json:"totalApplications"`
	CompanyName       string    `json:"companyName"`
}

type ErrorResp struct {
	Message string `json:"message"`
}

type SuccessResp struct {
	Message string `json:"message"`
}

type EducationStruct struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ExperienceStruct struct {
	Dates []string `json:"dates"`
	Name  string   `json:"name"`
	Url   string   `json:"url"`
}

type ProfileBufferStruct struct {
	Education  []string `json:"education"`
	Experience []string `json:"experience"`
}

type FindUserResp struct {
	User UserStruct
	Err  error
}

type JobStructResp struct {
	Id  primitive.ObjectID `bson:"_id" json:"id"`
	Job JobStruct          `json:"job"`
}

type UserStructResp struct {
	Id              string        `bson:"_id" json:"id"`
	Name            string        `json:"name"`
	Email           string        `json:"email"`
	Address         string        `json:"address"`
	UserType        string        `json:"userType"`
	Password        string        `json:"password"`
	ProfileHeadline string        `json:"profileHeading"`
	Profile         ProfileStruct `json:"profile"`
}

type JobArrResp struct {
	Id                string    `bson:"_id" json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	PostedOn          time.Time `json:"postedOn"`
	TotalApplications int       `json:"totalApplications"`
	CompanyName       string    `json:"companyName"`
}
