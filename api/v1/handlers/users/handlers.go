package users

import (
	//"fmt"
	"net/http"
	"io/ioutil"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/negeek/short-access/db"
	"github.com/negeek/short-access/utils"
		)

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request){
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
	
	// check if email exists. Tell user to input another email
	var emailExists bool
	emailErr:=dbPool.QueryRow(context.Background(),  "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", newUser.Email).Scan(&emailExists)
	if emailErr != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , emailErr.Error(), nil)
		return	
	}
	if emailExists==true{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Email already exist", nil)
		return
	}

	// Insert the new user into the database
	_, dbErr := dbPool.Exec(context.Background(), "INSERT INTO users (id, password, email) VALUES ($1, $2, $3)",newUserId, newUser.Password, newUser.Email)
	if dbErr != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , dbErr.Error(), nil)
		return	
	}

	//create token using the user details
	signedToken, errToken:=utils.CreateJwtToken(newUserId,newUser.Email)
	if errToken != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Token Error", nil)
		return	
	}

	// now return user token
	utils.JsonResponse(w, true, http.StatusCreated , "Successfully Joined", map[string]interface{}{"email":newUser.Email,
	"access_token":signedToken})
	return	
}

func NewToken(w http.ResponseWriter, r *http.Request){
	/*since the token doesn't expire. Get the email the user used on the api to signup and then get user_id and generate
	the token. The token will always be the same due to the claims */

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

	// get the email from the body
	type Body struct{
		Email string `json:"email"`
	}

	var emailBody *Body
	jsErr:=json.Unmarshal([]byte(body),&emailBody)

	if jsErr != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , jsErr.Error(), nil)
		return
	}

	// validate if email exist in db and get the user_id to create token.
	var userId uuid.UUID

	dbErr := dbPool.QueryRow(context.Background(), "select id from users where email=$1", emailBody.Email).Scan(&userId)
	if dbErr != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something Went Wrong. Make sure Email used to signUp is what is provided", nil)
		return
	}

	//create token using the user details
	signedToken, errToken:=utils.CreateJwtToken(userId,emailBody.Email)
	if errToken != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Token Error", nil)
		return	
	}

	// now return user token
	utils.JsonResponse(w, true, http.StatusCreated , "Token created Successfully", map[string]interface{}{"email":emailBody.Email,
	"access_token":signedToken})
	return	

}



