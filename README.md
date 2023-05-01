Short-Access is a free and powerful URL Shortener built with Golang.

### How to Use
##### To sign up:

`curl -X POST 'http://localhost:8080/api/v1/user_mgt/join/' -H 'Content-Type: application/json' -d'{"email":"dlionking77@gmail.com","password":"dlionking77"}'`

###### example response: 

`{"data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBhdHJpY2tAZ21haWwuY29tIiwiaWQiOiI5MWEyMWU3NS0yMWU5LTQwYzAtOTk3MS0yNTBiN2UwMzE4NDEifQ.nhIQKPJrgsGbWQqCdSGzwrkQUgzSmeLhJ3XXgsn1xJI","email":"patrick@gmail.com"},"success":true}`

##### To shorten your URL: 

`curl -X POST 'http://localhost:8080/api/v1/url/shorten/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBhdHJpY2tAZ21haWwuY29tIiwiaWQiOiI5MWEyMWU3NS0yMWU5LTQwYzAtOTk3MS0yNTBiN2UwMzE4NDEifQ.nhIQKPJrgsGbWQqCdSGzwrkQUgzSmeLhJ3XXgsn1xJI' -d '{"url":"https://pkg.go.dev/net/http#pkg-constants"}'`

###### example response: 

`{"data":{"origin":"https://pkg.go.dev/net/http#pkg-constants","slug":"2","url":"http://localhost:8080/2"},"success":true}`
