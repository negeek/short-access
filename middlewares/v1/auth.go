package v1
import (
	//"fmt"
	"net/http"
	"context"
	"strings"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/repository/v1/user"
		)

func AuthenticationMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		var oldUser =&user.User{}
		oldUser.Email=claim.Email
		exist:= oldUser.EmailExists()
		if exist != true {
			utils.JsonResponse(w, false, http.StatusUnauthorized ,"Invalid User", nil)	
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "user", claim.ID)
		rWithUser := r.WithContext(ctxWithUser)
        handler.ServeHTTP(w, rWithUser)
    })
}