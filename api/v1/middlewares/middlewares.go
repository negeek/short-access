package middlewares
import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"context"
	"github.com/negeek/short-access/db"
	"strings"
	"github.com/jackc/pgx/v4/pgxpool"
		)

var dbPool 	*pgxpool.Pool
var dbErr 	error

// i need to know if the user is authenticated and if the token provided is correct.
func AuthenticationMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bypassUrls := []string{"/api/v1/user_mgt/join/", "/"}
        
        // Check if the requested URL matches any of the bypass URLs
        for _, url := range bypassUrls {
            if strings.HasPrefix(r.URL.Path, url) {
                // The requested URL matches a bypass URL, so skip the authentication middleware
                handler.ServeHTTP(w, r)
                return
            }
        }

        // get the token 
		dbPool, dbErr = db.Connect()
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr.Error(),

			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusInternalServerError )
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("db connection error: %s\n", dbErr)
			return
		}

		bearerToken:= r.Header.Get("Authorization")
		if bearerToken==""{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Provide auth token")
			return	
		}
		// check database if token actually exist and then the associated user should be passed as context
		bearerTokenArr := strings.Split(bearerToken, " ")
		bearer, token:=bearerTokenArr[0], bearerTokenArr[1]
		if bearer!="Bearer"{
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
			}
			responseJson,_:=json.Marshal(response)
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, fmt.Sprintf("%s\n",responseJson))
			fmt.Printf("Invalid token")
			return	
		}

        var userId string
		dbErr = dbPool.QueryRow(context.Background(),  "select user_id from tokens where token=$1", token).Scan(&userId)
		if dbErr != nil {
			response:=map[string]interface{}{
				"success": false,
				"data":map[string]string{
				},
				"message": dbErr.Error(),

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