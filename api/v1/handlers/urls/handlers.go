package urls 

import (
	//"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"context"
	"encoding/json"
	"github.com/negeek/short-access/db"
	"github.com/gorilla/mux"
	"github.com/negeek/short-access/utils"
		)

func Shorten( w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		baseUrl:=os.Getenv("BASE_URL")
		url_length:=9
		dbPool, connErr := db.Connect()
		if connErr != nil {
			utils.JsonResponse(w, false, http.StatusInternalServerError , connErr.Error(), nil)
			return
		}

		body, err:= ioutil.ReadAll(r.Body)
		if err != nil{
			utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
			return
		}

		type UrlBody struct{
			Url string `json:"url"`
		}

		var url *UrlBody
		jsErr:=json.Unmarshal([]byte(body),&url)

		if jsErr != nil{
			utils.JsonResponse(w, false, http.StatusBadRequest , jsErr.Error(), nil)
			return
		}

		// get the user_id from context.
		// check if the url exists
		// if not get latest id and the convert it to base62 and store the new url to db
		userId := r.Context().Value("user")
		var urlId int
		var shortUrl string
		dbErr:= dbPool.QueryRow(context.Background(), "select id, short_url from urls where original_url=$1 and user_id=$2", url.Url, userId).Scan(&urlId, &shortUrl)
		if dbErr!=nil {
			if dbErr.Error()=="no rows in result set" {
				// get latest id in db
				var lastId int
				dbErr=dbPool.QueryRow(context.Background(),  "select max(id) from urls").Scan(&lastId)
				nextId:=lastId+1
				newShortUrl:= utils.ShortAccess(nextId, url_length)
	
				// Insert the new url into the database
				_, dbErr1 := dbPool.Exec(context.Background(), "INSERT INTO urls (id, user_id, original_url, short_url) VALUES ($1, $2, $3, $4)",nextId, userId, url.Url, newShortUrl)
				if dbErr1 != nil {
					utils.JsonResponse(w, false, http.StatusBadRequest, dbErr1.Error(), nil)
					return
					
				}
				
				utils.JsonResponse(w, true, http.StatusCreated ,"Successfully shortened url", map[string]interface{}{
					"origin":url.Url,
					"slug":newShortUrl,
					"url": baseUrl+"/"+newShortUrl,
				})
				return
			}
			utils.JsonResponse(w, false, http.StatusBadRequest , dbErr.Error(), nil)
			return
		}
		utils.JsonResponse(w, true, http.StatusOK ,"success", map[string]interface{}{
			"origin":url.Url,
			"slug":shortUrl,
			"url": baseUrl+"/"+shortUrl,
		})
		return
	}
}

func UrlRedirect( w http.ResponseWriter, r *http.Request){
	if r.Method=="GET"{
		dbPool, connErr := db.Connect()
		if connErr != nil {
			utils.JsonResponse(w, false, http.StatusInternalServerError , connErr.Error(), nil)
			return
		}

		// get the original url
		shortUrl := mux.Vars(r)["slug"]
		var originalUrl string
		dbErr:= dbPool.QueryRow(context.Background(),  "select original_url from urls where short_url=$1", shortUrl).Scan(&originalUrl)
		if dbErr !=nil{
			utils.JsonResponse(w, false, http.StatusBadRequest , dbErr.Error(), nil)
			return
		}
		http.Redirect(w, r, originalUrl, http.StatusTemporaryRedirect)
		return
	}
}