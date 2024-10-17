package jwt

import (
	"encoding/base64"
	"encoding/json"

	"assn.com/models"
)

type JWTResp struct {
	JWT string
}

type Headers struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Payload struct {
	Email    string `json:"email"`
	UserType string `json:"userType"`
}

func CreateJWT(user *models.UserStruct) string {
	payloadArr, _ := json.Marshal(Payload{Email: user.Email, UserType: user.UserType})
	headersArr, _ := json.Marshal(Headers{Alg: "Base64", Typ: "JWT"})

	payloadStr := base64.StdEncoding.EncodeToString(payloadArr)
	headersStr := base64.StdEncoding.EncodeToString(headersArr)

	return headersStr + "." + payloadStr
}

func ParseJWT(jwt string) (Payload, error) {
	var split int

	for idx, val := range jwt {
		if string(val) == "." {
			split = idx
			break
		}
	}

	decodedStr, _ := base64.StdEncoding.DecodeString(jwt[split+1:])

	data := Payload{}
	err := json.Unmarshal([]byte(decodedStr), &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
