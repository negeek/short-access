package users

import (
	//"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/repository/v1/user"
		)

func SignUp(w http.ResponseWriter, r *http.Request){
	body, err:= ioutil.ReadAll(r.Body)
	// if there is error reading body set statuscode and print out the error
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return	
	}

	// create user
	var newUser user.User
	err=json.Unmarshal([]byte(body),&newUser)
	newUser.Id=uuid.New()

	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return	
	}
	
	// check if email exists. Tell user to input another email
	emailExists:=newUser.EmailExists()
	if emailExists==true{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Email already exist", nil)
		return
	}

	// Insert the new user into the database
	err=newUser.Create()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return	
	}

	//create token using the user details
	signedToken, errToken:=utils.CreateJwtToken(newUser.Id,newUser.Email)
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

	body, err:= ioutil.ReadAll(r.Body)
	// if there is error reading body set statuscode and print out the error
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return	
	}

	var oldUser user.User
	err=json.Unmarshal([]byte(body),&oldUser)

	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}

	// validate if email exist in db and get the user_id to create token.
	_,exist:= oldUser.FindByEmail()
	if exist != true{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something Went Wrong. Make sure Email used to signUp is what is provided", nil)
		return
	}

	//create token using the user details
	signedToken, errToken:=utils.CreateJwtToken(oldUser.Id,oldUser.Email)
	if errToken != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Token Error", nil)
		return	
	}

	// now return user token
	utils.JsonResponse(w, true, http.StatusCreated , "Token created Successfully", map[string]interface{}{"email":oldUser.Email,
	"access_token":signedToken})
	return	

}



