package users

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/negeek/short-access/db"
		)

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		dbPool, connErr := db.Connect()
		if connErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": connErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusInternalServerError )
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error: %s\n", connErr)
			return
		}
		body, err:= ioutil.ReadAll(r.Body)
		// if there is error reading body set statuscode and print out the error

		if err != nil{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": err.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Could not read body: %s\n", err)
			return
		}

		// create user
		var newUser *User
		newUserId:= uuid.New()
		jsErr:=json.Unmarshal([]byte(body),&newUser)

		if jsErr != nil{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": jsErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Error unmarshaling json: %s\n", jsErr)
			return
		}

		// Insert the new user into the database
		_, dbErr := dbPool.Exec(context.Background(), "INSERT INTO users (id, password, email) VALUES ($1, $2, $3)",newUserId, newUser.Password, newUser.Email)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error1: %s\n", dbErr)
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
			fmt.Printf("token error: %s\n", errToken)
			return
		}

		// store it in db
		_, dbErr2 := dbPool.Exec(context.Background(), "INSERT INTO tokens (user_id, token) VALUES ($1, $2)",newUserId,signedToken)
		if dbErr2 != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr2.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error2: %s\n", dbErr2)
			return
		}

		// now return user token
		response:=map[string]interface{}{
			"success": true,
			"data":map[string]string{
				"email":newUser.Email,
				"access_token":signedToken,
			},
		}
		responseJson,rjsErr:=json.Marshal(response)
		if rjsErr != nil{
			fmt.Printf("Error marshaling response json: %s\n", rjsErr)
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
		return
	}

}

func SignIn(w http.ResponseWriter, r *http.Request){
	// get the email and password for a post request
	if r.Method == "POST"{
		dbPool, connErr := db.Connect()
		if connErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": connErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusInternalServerError )
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error: %s\n", connErr)
			return
		}
		body, err:= ioutil.ReadAll(r.Body)
		// if there is error reading body set statuscode and print out the error

		if err != nil{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": err.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Could not read body: %s\n", err)
			return
		}
		var user *User
		jsErr:=json.Unmarshal([]byte(body),&user)

		if jsErr != nil{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": jsErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Error unmarshaling json: %s\n", jsErr)
			return
		}

		//validating email and password
		var userId string
		var email string

		dbErr := dbPool.QueryRow(context.Background(), "select id, email from users where email=$1 and password=$2", user.Email, user.Password ).Scan(&userId, &email)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error1: %s\n", dbErr)
			return
		}
	
		// since email and password exist. Return the access token
		var userToken string
		dbErr = dbPool.QueryRow(context.Background(),  "select token from tokens where user_id=$1", userId).Scan(&userToken)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error2: %s\n", dbErr)
			return
		}

		// now return user token
		response:=map[string]interface{}{
			"success": true,
			"data":map[string]string{
				"email":email,
				"access_token":userToken,
			},
		}
		responseJson,rjsErr:=json.Marshal(response)
		if rjsErr != nil{
			fmt.Printf("Error marshaling response json: %s\n", rjsErr)
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
		return
	}

}