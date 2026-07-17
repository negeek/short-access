package users

import (
	//"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/repository/v1/user"
		)

func SignUp(w http.ResponseWriter, r *http.Request){
	// read the sign up details from the request body
	newUser, err := utils.DecodeBody[user.User](r)
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}

	newUser.Id=uuid.New()
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

	// read the login details from the request body
	oldUser, err := utils.DecodeBody[user.User](r)
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}

	// validate if email and password is correct.
	_,exist:= oldUser.Authenticate()
	if exist != true{
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something Went Wrong. Check your email and password.", nil)
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



