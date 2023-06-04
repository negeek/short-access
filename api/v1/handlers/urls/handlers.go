package urls 

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"context"
	"encoding/json"
	"github.com/negeek/short-access/db"
	"github.com/gorilla/mux"
	"github.com/negeek/short-access/utils"
		)


func base10To62(quotient int)string{
	numMap:=map[int]string{
		0:"0",1:"1",2:"2",3:"3",4:"4",5:"5",6:"6",7:"7",8:"8",9:"9",
		10:"A",11:"B",12:"C",13:"D",14:"E",15:"F",16:"G",17:"H",18:"I",
		19:"J",20:"K",21:"L",22:"M",23:"N",24:"O",25:"P",26:"Q",27:"R",
		28:"S",29:"T",30:"U",31:"V",32:"W",33:"X",34:"Y",35:"Z", 36:"a",
		37:"b",38:"c",39:"d",40:"e",41:"f",42:"g",43:"h",44:"i",45:"j",
		46:"k",47:"l",48:"m",49:"n",50:"o",51:"p",52:"q",53:"r",54:"s",
		55:"t",56:"u",57:"v",58:"w",59:"x",60:"y",61:"z",
	}

	resStr:=""
	var rem int

	// perform conversion and add to resStr
	for{
		quotient,rem= quotient/62, quotient%62
		resStr+=numMap[rem]
		if quotient<1{
			break
		}
	}

	// reverse the resStr,that is the correct result
	resArr := []byte(resStr)
    for i, j := 0, len(resArr)-1; i < j; i, j = i+1, j-1 {
        resArr[i], resArr[j] = resArr[j], resArr[i]
    }
	resStr=string(resArr)
    return resStr
}


func Shorten( w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		baseUrl:=os.Getenv("BASE_URL")
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
		//check if the url exists
		// if not get latest id and the convert it to base62 and store the new url to db
		userId := r.Context().Value("user")
		var urlId int
		var shortUrl string
		dbErr:= dbPool.QueryRow(context.Background(),  "select id, short_url from urls where original_url=$1 and user_id=$2", url.Url, userId).Scan(&urlId, &shortUrl)
		if dbErr.Error()!="no rows in result set" {
			utils.JsonResponse(w, false, http.StatusBadRequest , dbErr.Error(), nil)
			return
		}
		if dbErr.Error()=="no rows in result set" {
			// get latest id in db
			var lastId int
			dbErr=dbPool.QueryRow(context.Background(),  "select last_value from urls_id_seq").Scan(&lastId)
			nextId:=lastId+1
			newShortUrl:=base10To62(nextId)

			// Insert the new url into the database
			_, dbErr1 := dbPool.Exec(context.Background(), "INSERT INTO urls (id, user_id, original_url, short_url) VALUES ($1, $2, $3, $4)",nextId, userId, url.Url, newShortUrl)
			if dbErr1 != nil {
				utils.JsonResponse(w, false, http.StatusBadRequest, dbErr1.Error(), nil)
				return
				
			}
			
			utils.JsonResponse(w, true, http.StatusCreated , dbErr.Error(), map[string]interface{}{
				"origin":url.Url,
				"slug":newShortUrl,
				"url": baseUrl+"/"+newShortUrl,
			})
			return
			}

		utils.JsonResponse(w, true, http.StatusOK , dbErr.Error(), map[string]interface{}{
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