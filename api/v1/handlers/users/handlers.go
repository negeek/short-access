package users

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
		)

type Token struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccessToken string  `json:"access_token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func SignUp(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool){
	// get the email and password for a post request
	if r.Method == "POST"{
		body, err:= ioutil.ReadAll(r.Body)
		// if there is error reading body set statuscode and print out the error

		if err != nil{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": err,

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Could not read body: %s\n", err)
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
				"message": jsErr,

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Error unmarshaling json: %s\n", jsErr)
		}

		// Insert the new user into the database
		userResult, dbErr := dbPool.Exec(context.Background(), "INSERT INTO users (id, email, password) VALUES ($1, $2)",newUserId, newUser.Email, newUser.Password)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr,

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error: %s\n", dbErr)
		}

		//create token using the user details
		var (
			key []byte
			jwtToken   *jwt.Token
			signedToken string
		)
		key = []byte(os.Getenv("AUTH_KEY"))
		t = jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{ 
			"id": newUserId,
			"email": newUser.Email, 
		})
		signedToken, errToken := t.SignedString(key) 
		if errToken != nil{
			fmt.Printf("token error: %s\n", errToken)
		}



		// store it in db
		userResult, dbErr := dbPool.Exec(context.Background(), "INSERT INTO tokens (user_id, token) VALUES ($1, $2)",newUserId,signedToken)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr,

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error: %s\n", dbErr)
		}

		// now return user token

		response:=map[string]interface{}{
			"success": true,
			"data":map[string]string{
				"email":newUser.Email,
				"access_token":signedToken,
			},
		}
		responseJson,rjsErr=json.Marshal(response)
		if rjsErr != nil{
			fmt.Printf("Error marshaling response json: %s\n", rjsErr)
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
	}

}