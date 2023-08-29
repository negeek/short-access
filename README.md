# Short-Access 

Short-Access is a free and powerful URL Shortener built with Golang.

## How to Use
#### To sign up:

`curl -X POST 'https://shrt-acc.onrender.com/api/v1/user_mgt/join/' -H 'Content-Type: application/json' -d'{"email":"patrick@gmail.com","password":"dlionking77"}'`

###### example response: 

`{"success":true,"message":"Successfully Joined","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY","email":"patrick@gmail.com"}}`

#### To shorten your URL: 

`curl -X POST 'https://shrt-acc.onrender.com/api/v1/url/shorten/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY' -d '{"url":"https://pkg.go.dev/net/http#pkg-constants"}'`

###### example response: 

`{"success":true,"message":"Successfully shortened url","data":{"origin":"https://pkg.go.dev/net/http#pkg-constants","slug":"00000001u","url":"https://shrt-acc.onrender.com/00000001u"}}`


#### To get back access token:
`curl -X POST 'https://shrt-acc.onrender.com/api/v1/user_mgt/new_token/' -H 'Content-Type: application/json' -d'{"email":"patrick@gmail.com"}'`

###### example response:
`{"success":true,"message":"Token created Successfully","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY","email":"patrick@gmail.com"}}`



