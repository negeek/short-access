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
		// verify the jwt token
		bearerTokenArr := strings.Split(bearerToken, " ")

		if len(bearerTokenArr) != 2 {
			utils.JsonResponse(w, false, http.StatusUnauthorized, "Invalid Authorisation Header", nil)
			return
		}
		bearer, token:=bearerTokenArr[0], bearerTokenArr[1]
		if bearer!="Bearer"{
			utils.JsonResponse(w, false, http.StatusUnauthorized , "Invalid Authorisation Header", nil)
			return	
		}
		
		claim,err:= utils.VerifyJwt(token)
		if err != nil{
			utils.JsonResponse(w, false, http.StatusUnauthorized , "Invalid Token", nil)
			return
		}

		// verify claims 
		// check db if user actually exist
		var exists bool
		dbErr = dbPool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", claim.ID).Scan(&exists)
		if dbErr != nil {
			utils.JsonResponse(w, false, http.StatusUnauthorized , dbErr.Error(), nil)	
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "user", claim.ID)
		rWithUser := r.WithContext(ctxWithUser)
        handler.ServeHTTP(w, rWithUser)
    })
}