// package urls
// import(
// 	"fmt"
// 	"testing"
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"encoding/json"
// 	"github.com/negeek/short-access/repository/v1/user"
// 	"github.com/negeek/short-access/api/v1/users"
// 	test "github.com/negeek/short-access/utils"
// 	//"time"
// )

// func TestShorten(t *testing.T){
// 	test.Setup()
// 	//create user first
// 	jsonBody := `{"email":"delet@yahoo.com","password":"dlionking77"}`
//  	bodyReader := bytes.NewReader([]byte(jsonBody))
// 	req, err := http.NewRequest("POST","/join/",bodyReader)
// 	if err != nil {
//         t.Fatal(err)
//     }
// 	rr := httptest.NewRecorder()
//     handler := http.HandlerFunc(users.SignUp)
// 	handler.ServeHTTP(rr, req)
//     // Check the status code is what we expect.
//     if status := rr.Code; status != http.StatusCreated {
//         t.Errorf("handler returned wrong status code: got %v want %v",
//             status, http.StatusCreated)
		
//     }

	
// 	var resp=test.Response{}
// 	err=json.Unmarshal([]byte(rr.Body.String()),&resp)
// 	if err != nil {
//         t.Fatal(err)
//     }
// 	s:=resp.Data
// 	fmt.Println(s["email"])

// 	// shorten
// 	// jsonBody = `{"url":"https://www.digitalocean.com/community/tutorials/how-to-install-postgresql-on-ubuntu-20-04-quickstart"}`
//  	// bodyReader = bytes.Nfmt.Println(rr.Body)NewReader([]byte(jsonBody))
// 	// req, err = http.NewRequest("POST","/shorten/",bodyReader)
// 	// if err != nil {
//     //     t.Fatal(err)
//     // }
// 	// req.Header.Set("Authorisation", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjdhZDk4MTU3LTAxZTUtNDgyMS1iZDU0LWRjODAwZDc0MmZiZCIsIkVtYWlsIjoiZGVsZXRlQHlhaG9vLmNvbSJ9.9YHMfbU93JyFPqaqJ-fnRN4Ftm9Dqqqyh81YiJK2b58")
// 	// client := http.Client{
// 	// 	Timeout: 30 * time.Second,
// 	// }
   
// 	//  res, err := client.Do(req)
// 	//  fmt.Println("err",err)
// 	//  fmt.Println(res)
   
// 	// rr = httptest.NewRecorder()
//     // handler = http.HandlerFunc(Shorten)
// 	// handler.ServeHTTP(rr, req)
// 	// fmt.Println(rr.Body.String())
//     // // Check the status code is what we expect.
//     // if status := rr.Code; status != http.StatusCreated {
//     //     t.Errorf("handler returned wrong status code: got %v want %v",
//     //         status, http.StatusCreated)
// 	// 	return
		
//     // }

// 	//clean Up
// 	fmt.Println("cleaning up.....")
// 	var newUser user.User
// 	err=json.Unmarshal([]byte(jsonBody),&newUser)
// 	if err != nil {
//         t.Fatal(err)
//     }
// 	newUser.Delete()
// 	fmt.Println("cleaning done")



	

// }