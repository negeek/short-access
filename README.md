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

#### To create custom URL: 

`curl -X POST 'https://shrt-acc.onrender.com/api/v1/url/custom/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY' -d '{"url":"https://pkg.go.dev/net/http#pkg-constants","short_url":"negeek"}'`

###### example response: 

`{"success":true,"message":"Successfully created custom url","data":{"origin":"https://pkg.go.dev/net/http#pkg-constants","slug":"negeek","url":"https://shrt-acc.onrender.com/negeek"}}`

#### To get list of your URLs: 

`curl -X GET 'https://shrt-acc.onrender.com/api/v1/url/urls/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY'`

###### example response:

`{"success":true,"message":"","data":[{"id":1,"url":"https://pkg.go.dev/net/http#pkg-constants","short_url":"000000001","is_custom":false,"date_created":"2023-06-08T10:54:44.044265Z","date_updated":"2023-08-08T10:54:44.044266Z"},{"id":2,"url":"https://pkg.go.dev/net/http#pkg-constants","short_url":"negeek","is_custom":true,"date_created":"2023-08-08T10:59:33.817736Z","date_updated":"2023-08-08T10:59:33.817737Z"}]}`

#### To get back access token:

`curl -X POST 'https://shrt-acc.onrender.com/api/v1/user_mgt/new_token/' -H 'Content-Type: application/json' -d'{"email":"patrick@gmail.com", "password":"dlionking77}'`

###### example response:

`{"success":true,"message":"Token created Successfully","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY","email":"patrick@gmail.com"}}`



