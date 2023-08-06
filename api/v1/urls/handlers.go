package urls 

import (
	//"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/repository/v1/url"
	"github.com/negeek/short-access/repository/v1/number"
		)

type NumberStore struct {
	Number int
	Step int
	End int
}

// Depending on traffic. But i will be using up 100 numbers  before storing in DB
var numberStore=&NumberStore{0,100,100}

func Shorten( w http.ResponseWriter, r *http.Request){
	// instead of wasting number check if url exists then just give payload and also check for latest number before updating struct
	baseUrl:=os.Getenv("BASE_URL")
	url_length:=9
	body, err:= ioutil.ReadAll(r.Body)
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}

	var newUrl url.Url
	err=json.Unmarshal([]byte(body),&newUrl)
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}

	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went Wrong. Try again", nil)
		return
	}
	newUrl.UserId =userId
	_,exist:=newUrl.FindByOriginalUrl()
	if exist == true{
		utils.JsonResponse(w, true, http.StatusCreated ,"Successfully shortened url", map[string]interface{}{
			"origin":newUrl.Url,
			"slug":newUrl.ShortUrl,
			"url": baseUrl+"/"+newUrl.ShortUrl,
		})
		return
	}

	var newNum =&number.Number{}
	newNum.Step=numberStore.Step
	newNum.Id=1
	if numberStore.Number==0{
		// server is restarted or just started
		_,exist=newNum.FindById()
		if exist==false{
			err=newNum.CreateOrUpdate()
			if err!= nil{
				utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
				return
			}
			// started
			numberStore.Number=1
			numberStore.End=newNum.Number
		}else{
			// re-started
			numberStore.Number=newNum.Number+1
			err=newNum.CreateOrUpdate()
			if err!= nil{
				utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
				return
			}
			numberStore.End=newNum.Number
		}

	}else{
		if numberStore.Number >= numberStore.End{
			err=newNum.CreateOrUpdate()
			if err!= nil{
				utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
				return
			}
			numberStore.Number+=1
			numberStore.End=newNum.Number
		}else{
			numberStore.Number+=1
		}
	}

	newUrl.ShortUrl=utils.ShortAccess(numberStore.Number, url_length)
	err=newUrl.Create()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went Wrong. Try again", nil)
		return
	}
	utils.JsonResponse(w, true, http.StatusCreated ,"Successfully shortened url", map[string]interface{}{
		"origin":newUrl.Url,
		"slug":newUrl.ShortUrl,
		"url": baseUrl+"/"+newUrl.ShortUrl,
	})
	return
}

func UrlRedirect( w http.ResponseWriter, r *http.Request){
	var oldUrl =&url.Url{}
	// get the original url
	oldUrl.ShortUrl = mux.Vars(r)["slug"]
	_,exist:=oldUrl.FindByShortUrl()
	if exist != true{
		utils.JsonResponse(w, false, http.StatusBadRequest,"Something Went wrong. Make sure is Valid or Shorten again" , nil)
		return
	}
	http.Redirect(w, r, oldUrl.Url, http.StatusTemporaryRedirect)
	return
}
