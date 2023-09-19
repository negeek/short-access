package urls
import(
	"fmt"
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"github.com/negeek/short-access/repository/v1/user"
	"github.com/negeek/short-access/repository/v1/url"
	"github.com/negeek/short-access/api/v1/users"
	test "github.com/negeek/short-access/utils"
	//"time"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
)

func TestShorten(t *testing.T){
	test.Setup()
	//create user first
	jsonBody := `{"email":"delet@yahoo.com","password":"dlionking77"}`
 	bodyReader := bytes.NewReader([]byte(jsonBody))
	req, err := http.NewRequest("POST","/join/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	rr := httptest.NewRecorder()
    handler := http.HandlerFunc(users.SignUp)
	handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
		
    }

	var resp=test.Response{}
	err=json.Unmarshal([]byte(rr.Body.String()),&resp)
	if err != nil {
        t.Fatal(err)
    }
	// resp.Data is of type interface{} and contains the data
	data, ok := resp.Data.(interface{})
	if !ok {
		fmt.Println("Data is not of type interface{}")
		
	}

	token, ok := data["access_token"].(string)
	if !ok {
		fmt.Println("token not found or not of type string.")
		
	}

	//shorten
	jsonBody2 := `{"original_url":"https://www.test.com/test"}`
 	bodyReader = bytes.NewReader([]byte(jsonBody2))
	req, err = http.NewRequest("POST","/shorten/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	req.Header.Set("Authorization", "Bearer "+token)

	rr = httptest.NewRecorder()
    handler2 := v1middlewares.AuthenticationMiddleware(http.HandlerFunc(Shorten))
	handler2.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)		
    }

	//clean Up
	defer func(){
		fmt.Println("cleaning up.....")
		var newUser user.User
		err=json.Unmarshal([]byte(jsonBody),&newUser)
		if err != nil {
			t.Fatal(err)
		}
		newUser.Delete()

		var newUrl url.Url
		err=json.Unmarshal([]byte(jsonBody2),&newUrl)
		if err != nil {
			t.Fatal(err)
		}
		newUrl.Delete()
		fmt.Println("cleaning done")
	}()
}

func TestCustomUrl(t *testing.T){
	test.Setup()
	//create user first
	jsonBody := `{"email":"delet@yahoo.com","password":"dlionking77"}`
 	bodyReader := bytes.NewReader([]byte(jsonBody))
	req, err := http.NewRequest("POST","/join/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	rr := httptest.NewRecorder()
    handler := http.HandlerFunc(users.SignUp)
	handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
		
    }

	var resp=test.Response{}
	err=json.Unmarshal([]byte(rr.Body.String()),&resp)
	if err != nil {
        t.Fatal(err)
    }
	// resp.Data is of type interface{} and contains the data
	data, ok := resp.Data.(interface{})
	if !ok {
		fmt.Println("Data is not of type interface{}")
		
	}

	token, ok := data["access_token"].(string)
	if !ok {
		fmt.Println("token not found or not of type string.")
		
	}

	//shorten
	jsonBody2 := `{"original_url":"https://www.test.com/test","short_url":"test"}`
 	bodyReader = bytes.NewReader([]byte(jsonBody2))
	req, err = http.NewRequest("POST","/custom/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	req.Header.Set("Authorization", "Bearer "+token)

	rr = httptest.NewRecorder()
    handler2 := v1middlewares.AuthenticationMiddleware(http.HandlerFunc(CustomUrl))
	handler2.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)		
    }

	//clean Up
	defer func(){
		fmt.Println("cleaning up.....")
		var newUser user.User
		err=json.Unmarshal([]byte(jsonBody),&newUser)
		if err != nil {
			t.Fatal(err)
		}
		newUser.Delete()

		var newUrl url.Url
		err=json.Unmarshal([]byte(jsonBody2),&newUrl)
		if err != nil {
			t.Fatal(err)
		}
		newUrl.Delete()
		fmt.Println("cleaning done")
	}()
}