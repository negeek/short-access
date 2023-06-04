package middlewares
import (
	//"fmt"
	"net/http"
	"context"
	"github.com/negeek/short-access/db"
	"strings"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/utils"
		)

var dbPool 	*pgxpool.Pool
var dbErr 	error

// i need to know if the user is authenticated and if the token provided is correct.
func AuthenticationMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // get the token 
		dbPool, dbErr = db.Connect()
		if dbErr != nil {
			utils.JsonResponse(w, false, http.StatusInternalServerError , dbErr.Error(), nil)
			return
		}
		bearerToken:= r.Header.Get("Authorization")
		if bearerToken==""{
			utils.JsonResponse(w, false, http.StatusUnauthorized , "Provide Auth Token", nil)
			return	
		}
		// check database if token actually exist and then the associated user should be passed as context
		bearerTokenArr := strings.Split(bearerToken, " ")
		bearer, token:=bearerTokenArr[0], bearerTokenArr[1]
		if bearer!="Bearer"{
			utils.JsonResponse(w, false, http.StatusUnauthorized , "Invalid Token", nil)
			return	
		}
        var userId string
		dbErr = dbPool.QueryRow(context.Background(),  "select user_id from tokens where token=$1", token).Scan(&userId)
		if dbErr != nil {
			utils.JsonResponse(w, false, http.StatusUnauthorized , dbErr.Error(), nil)	
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "user", userId)
		rWithUser := r.WithContext(ctxWithUser)
        handler.ServeHTTP(w, rWithUser)
    })
}