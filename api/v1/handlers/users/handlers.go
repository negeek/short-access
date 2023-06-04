package users

import (
	//"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/negeek/short-access/db"
	"github.com/negeek/short-access/utils"
		)

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		dbPool, connErr := db.Connect()
		if connErr != nil {
			utils.JsonResponse(w, false, http.StatusInternalServerError , connErr.Error(), nil)
			return	
		}
		body, err:= ioutil.ReadAll(r.Body)
		// if there is error reading body set statuscode and print out the error

		if err != nil{
			utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
			return	
		}

		// create user
		var newUser *User
		newUserId:= uuid.New()
		jsErr:=json.Unmarshal([]byte(body),&newUser)

		if jsErr != nil{
			utils.JsonResponse(w, false, http.StatusBadRequest , jsErr.Error(), nil)
			return	
		}

		// Insert the new user into the database
		_, dbErr := dbPool.Exec(context.Background(), "INSERT INTO users (id, password, email) VALUES ($1, $2, $3)",newUserId, newUser.Password, newUser.Email)
		if dbErr != nil {
			utils.JsonResponse(w, false, http.StatusBadRequest , dbErr.Error(), nil)
			return	
		}

		//create token using the user details
		var (
			key []byte
			jwtToken   *jwt.Token
			signedToken string
		)
		key = []byte(os.Getenv("AUTH_KEY"))
		jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{ 
			"id": newUserId,
			"email": newUser.Email, 
		})
		signedToken, errToken := jwtToken.SignedString(key) 
		if errToken != nil{
			utils.JsonResponse(w, false, http.StatusBadRequest , "Token Error", nil)
			return	
		}

		// store it in db
		_, dbErr2 := dbPool.Exec(context.Background(), "INSERT INTO tokens (user_id, token) VALUES ($1, $2)",newUserId,signedToken)
		if dbErr2 != nil {
			utils.JsonResponse(w, false, http.StatusBadRequest , dbErr2.Error(), nil)
			return	
		}

		// now return user token
		utils.JsonResponse(w, true, http.StatusCreated , "Successfully Joined", map[string]interface{}{"email":newUser.Email,
		"access_token":signedToken})
		return	
	}

}

