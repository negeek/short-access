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

type NumberStore struct {
	Start int
	Number int
	Step int
	End int
}

// Depending on traffic. But i will be using up 100 numbers  before storing in DB
var numberStore=&NumberStore{0,0,100,100}

func Shorten( w http.ResponseWriter, r *http.Request){
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
	// Handle 100 requests at once
	if numberStore.Number==0{
		// Probably server just started or it was shutdown and started again
		// get latest number
		var num int
		dbErr:=dbPool.QueryRow(context.Background(),  "select coalesce(max(number), 0) FROM numbers").Scan(&num)
		if dbErr != nil {
			utils.JsonResponse(w, false, http.StatusBadRequest, dbErr.Error(), nil)
			return
		}
		// update the struct
		numberStore.Start=num+1
		numberStore.Number=numberStore.Start
		numberStore.End=numberStore.Start+numberStore.Step
		
		// Insert the End number in numbers table in db
		_, dbErr1 := dbPool.Exec(context.Background(), "INSERT INTO numbers (number) VALUES ($1)",numberStore.End)
			if dbErr1 != nil {
				utils.JsonResponse(w, false, http.StatusBadRequest, dbErr1.Error(), nil)
				return
				
			}
			
	}else{
		// Server is still on so number still retains its original value
		// check if number has reached End. Or it has handled 100 requests
		if numberStore.Number>=numberStore.End{
			// update the struct
			numberStore.Start=numberStore.Number+1
			numberStore.Number=numberStore.Start
			numberStore.End=numberStore.Start+numberStore.Step
			_, dbErr1 := dbPool.Exec(context.Background(), "INSERT INTO numbers (number) VALUES ($1)",numberStore.End)
			if dbErr1 != nil {
				utils.JsonResponse(w, false, http.StatusBadRequest, dbErr1.Error(), nil)
				return	
			}
			
		}else{

			// number has n0t handled 100 requests yet. So just keep increasing it
			numberStore.Number+=1
		}
		
	}
	// since number is gotten. Convert to base62 to get slug and store as new url.
	userId := r.Context().Value("user")
	newShortUrl:= utils.ShortAccess(numberStore.Number, url_length)
	_, dbErr1 := dbPool.Exec(context.Background(), "INSERT INTO urls (user_id, original_url, short_url) VALUES ($1, $2, $3)",userId, url.Url, newShortUrl)
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

func UrlRedirect( w http.ResponseWriter, r *http.Request){
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
