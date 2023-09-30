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
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went wrong. Try again", nil)
		return
	}
	newUrl.UserId =userId
	_,exist:=newUrl.FindByOriginalUrl()
	if exist == true{
		utils.JsonResponse(w, true, http.StatusCreated ,"Successfully shortened url", map[string]interface{}{
			"original_url":newUrl.OriginalUrl,
			"short_url": baseUrl+"/"+newUrl.ShortUrl,
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
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went wrong. Try again", nil)
		return
	}
	utils.JsonResponse(w, true, http.StatusCreated ,"Successfully shortened url", map[string]interface{}{
		"original_url":newUrl.OriginalUrl,
		"short_url": baseUrl+"/"+newUrl.ShortUrl,
	})
	return
}

func UrlExpiry(w http.ResponseWriter, r *http.Request){
	// set expiry datetime of a url that has been shortened.
	body, err:= ioutil.ReadAll(r.Body)
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}
	var exp_dtl DateTimeExpiryDetail
	err=json.Unmarshal([]byte(body),&exp_dtl)
	if err != nil{
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went Wrong. Try again", nil)
		return
	}
	// ensure udate function works and findbyshorturl sets all url struct fields
	expire_at:=utils.ExpiryDateTime(exp_dtl.TimeUnit, exp_dtl.TimeValue)
	var newUrl url.Url
	newUrl.UserId =userId
	newUrl.ShortUrl=exp_dtl.ShortUrl
	_,exist:=newUrl.FindByShortUrl()
	if exist == false{
		utils.JsonResponse(w, true, http.StatusBadRequest ,"Url provided does not exist", nil)
		return
	}
	newUrl.Update()
	// end
	utils.JsonResponse(w, true, http.StatusOK ,"Successfully set url expiry", nil)
	return
}

func CustomUrl(w http.ResponseWriter, r *http.Request){
	baseUrl:=os.Getenv("BASE_URL")
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
	// check if short url exists before
	_,exist:=newUrl.FindByShortUrl()
	if exist == true{
		utils.JsonResponse(w, true, http.StatusBadRequest ,"Url provided exists already", nil)
		return
	}
	// now store new url
	newUrl.IsCustom=true
	err=newUrl.Create()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went wrong. Try again", nil)
		return
	}
	utils.JsonResponse(w, true, http.StatusCreated ,"Successfully created custom url", map[string]interface{}{
		"original_url":newUrl.OriginalUrl,
		"short_url": baseUrl+"/"+newUrl.ShortUrl,
	})
	return
}

func UrlFilter(w http.ResponseWriter, r *http.Request){
	userId, ok := r.Context().Value("user").(uuid.UUID)
	if !ok {
		utils.JsonResponse(w, false, http.StatusBadRequest , "Something went wrong. Try again", nil)
		return
	}
	var urll url.Url
	var result []url.Url
	urll.UserId =userId
	queryParams:=r.URL.Query()
	query,queryValues,err:=utils.Filter(queryParams,urll,"urls")
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}
	result,err=urll.UserUrls(query,queryValues)
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest , err.Error(), nil)
		return
	}
	utils.JsonResponse(w, true, http.StatusOK ,"", result)
	return
}

func UrlRedirect( w http.ResponseWriter, r *http.Request){
	var oldUrl =&url.Url{}
	oldUrl.ShortUrl = mux.Vars(r)["slug"]
	_,exist:=oldUrl.FindByShortUrl()
	if exist != true{
		utils.JsonResponse(w, false, http.StatusBadRequest,"Something went wrong. Make sure url is valid." , nil)
		return
	}
	oldUrl.AccessCount+=1
	err:=oldUrl.UpdateAccessCount()
	if err != nil {
		utils.JsonResponse(w, false, http.StatusBadRequest,"Something went wrong. Try again." , nil)
		return
	}
	http.Redirect(w, r, oldUrl.OriginalUrl, http.StatusTemporaryRedirect)
	return
}
