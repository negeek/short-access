package users
import(
	"fmt"
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"github.com/negeek/short-access/repository/v1/user"
	test "github.com/negeek/short-access/utils"
)
func TestSignUp(t *testing.T){
	test.Setup()
	jsonBody := `{"email":"delete@yahoo.com","password":"dlionking77"}`
 	bodyReader := bytes.NewReader([]byte(jsonBody))
	req, err := http.NewRequest("POST","/join/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	rr := httptest.NewRecorder()
    handler := http.HandlerFunc(SignUp)
	handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
	
    }

	//clean Up
	fmt.Println("cleaning up.....")
	var newUser user.User
	err=json.Unmarshal([]byte(jsonBody),&newUser)
	if err != nil {
        t.Fatal(err)
    }
	newUser.TestDelete()
	fmt.Println("cleaning done")

}

func TestNewToken(t *testing.T){
	test.Setup()
	// first create user
	jsonBody := `{"email":"delete@yahoo.com","password":"dlionking77"}`
 	bodyReader := bytes.NewReader([]byte(jsonBody))
	req, err := http.NewRequest("POST","/join/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	rr := httptest.NewRecorder()
    handler := http.HandlerFunc(SignUp)
	handler.ServeHTTP(rr, req)
    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
    }

	// now get new token
 	bodyReader = bytes.NewReader([]byte(jsonBody))
	req, err = http.NewRequest("POST","/new_token/",bodyReader)
	if err != nil {
        t.Fatal(err)
    }
	rr = httptest.NewRecorder()
    handler = http.HandlerFunc(NewToken)
	handler.ServeHTTP(rr, req)
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
		newUser.TestDelete()
		fmt.Println("cleaning done")
	}()

}