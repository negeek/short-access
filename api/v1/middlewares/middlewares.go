package middlewares
import (
	"fmt"
	"net/http"
	"io"
	"context"
	"github.com/negeek/short-access/db"
	"strings"
		)

dbPool, dbErr := db.Connect()
// i need to know if the user is authenticated and if the token provided is correct.

func AuthenticationMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // get the token 
		bearerToken:= r.Header.Get("Authorization")
		if bearerToken==""{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.http.StatusUnauthorized)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Could not read body: %s\n", err)
			return	
		}
		// check database if token actually exist and then the associated user should be passed as context
		bearer, token := strings.split(bearerToken, " ")
		if bearer!="Bearer"{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.http.StatusUnauthorized)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Could not read body: %s\n", err)
			return	
		}

        var userId string
		dbErr = dbPool.QueryRow(context.Background(),  "select user_id from tokens where token=$1", token).Scan(&userId)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db error: %s\n", dbErr)
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "user", userId)
		rWithUser := r.WithContext(ctxWithUser)
        handler.ServeHTTP(w, rWithUser)
    })
}