# Short-Access 

Short-Access is a free and powerful URL Shortener API built with Golang.

## How to Use
#### To sign up:
`curl -X POST 'https://shrt-acc.onrender.com/api/v1/user_mgt/join/' -H 'Content-Type: application/json' -d'{"email":"patrick@gmail.com","password":"dlionking77"}'`

###### example response: 
`{"success":true,"message":"Successfully Joined","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY","email":"patrick@gmail.com"}}`



#### To get back access token:
`curl -X POST 'https://shrt-acc.onrender.com/api/v1/user_mgt/new_token/' -H 'Content-Type: application/json' -d'{"email":"patrick@gmail.com", "password":"dlionking77}'`

###### example response:
`{"success":true,"message":"Token created Successfully","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY","email":"patrick@gmail.com"}}`



#### To shorten your URL: 
`curl -X POST 'https://shrt-acc.onrender.com/api/v1/url_mgt/shorten/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY' -d '{"original_url":"https://pkg.go.dev/net/http#pkg-constants"}'`

###### example response: 
`{"success":true,"message":"Successfully shortened url","data":{"id":1,"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp","short_access":"https://shrt-acc.onrender.com/00000001u","is_custom":false,"access_count":0,"expire_at":"0001-01-01T00:00:00Z","date_created":"2023-06-02T18:19:40.150795Z","date_updated":"2023-06-02T19:15:17.617083581Z"}}`



#### To create custom URL: 
`curl -X POST 'https://shrt-acc.onrender.com/api/v1/url_mgt/custom/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY' -d '{"original_url":"https://pkg.go.dev/net/http#pkg-constants","short_url":"negeek"}'`

###### example response: 
`{"success":true,"message":"Successfully created custom url","data":{"id":2,"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp","short_access":"https://shrt-acc.onrender.com/nethttp","is_custom":true,"access_count":0,"expire_at":"0001-01-01T00:00:00Z","date_created":"2023-07-02T18:19:40.150795004Z","date_updated":"2023-07-02T18:19:40.150797413Z"}}`



#### To set URL Expiry:
###### To set the time for url to expire: From the example request below. The url with id of 2 is set to expire 40 seconds from now. Other options for "time_unit" are "y", "mo", "d", "h", "m" denoting year, month, day, hour, minute  respectively.
`curl -X POST 'https://shrt-acc.onrender.com/api/v1/url_mgt/url_expiry/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY' -d '{"time_unit":"s","time_value":40,"url_id":2}'`

###### example response
`{"success":true,"message":"Successfully set url expiry","data":{"id":2,"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp","short_access":"https://shrt-acc.onrender.com/nethttp","is_custom":true,"access_count":0,"expire_at":"2023-10-02T19:15:57.615924048Z",""date_created":"2023-07-02T18:19:40.150795004Z","date_updated":"2023-07-02T18:19:40.150797413Z"}}`



#### To get list or filter your URLs: 
`curl -X GET 'https://shrt-acc.onrender.com/api/v1/url_mgt/' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY'`

`curl -X GET 'https://shrt-acc.onrender.com/api/v1/url_mgt/?id=2&short_url=nethttp' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjQ5MTc1OGYyLWM3OGYtNDE3MC05MDI0LWEzOWU5NTIxMjM0ZCIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.FVLuSkPnIHBaS46aFplaaDBzJc4IXM9hJ7xCnL8ZZyY'`

###### example responses:
`{"success":true,"message":"","data":[{"id":1,"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp","short_access":"https://shrt-acc.onrender.com/00000001u","is_custom":false,"access_count":0,"expire_at":"0001-01-01T00:00:00Z","date_created":"2023-06-02T18:19:40.150795Z","date_updated":"2023-06-02T19:15:17.617083581Z"},{"id":2,"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp","short_access":"https://shrt-acc.onrender.com/nethttp","is_custom":true,"access_count":0,"expire_at":"0001-01-01T00:00:00Z","date_created":"2023-07-02T18:19:40.150795004Z","date_updated":"2023-07-02T18:19:40.150797413Z"}]}`

`{"success":true,"message":"","data":[{"id":2,"original_url":"https://pkg.go.dev/net/http","short_url":"nethttp","short_access":"https://shrt-acc.onrender.com/nethttp","is_custom":true,"access_count":0,"expire_at":"0001-01-01T00:00:00Z","date_created":"2023-07-02T18:19:40.150795004Z","date_updated":"2023-07-02T18:19:40.150797413Z"}]}`



#### To Update(PATCH, PUT), Delete Url:
###### for example you want to update, delete the url with id of 2
`curl -X PATCH 'https://shrt-acc.onrender.com/api/v1/url_mgt/2' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjNiYzZkODQ3LWY5M2UtNDMyYS05NjcwLWZkOTFjYzRkZmY5YSIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.vQWYoUFs0_WZUIoBdvwms8iI3u_yqpQxPfAm5bLaJBc' -d '{"access_count":3}'`

`curl -X PUT 'https://shrt-acc.onrender.com/api/v1/url_mgt/2' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjNiYzZkODQ3LWY5M2UtNDMyYS05NjcwLWZkOTFjYzRkZmY5YSIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.vQWYoUFs0_WZUIoBdvwms8iI3u_yqpQxPfAm5bLaJBc' -d '{"access_count":3}'`

`curl -X DELETE 'https://shrt-acc.onrender.com/api/v1/url_mgt/2' -H 'Content-Type: application/json' -H 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjNiYzZkODQ3LWY5M2UtNDMyYS05NjcwLWZkOTFjYzRkZmY5YSIsIkVtYWlsIjoiZGxpb25AZ21haWwuY29tIn0.vQWYoUFs0_WZUIoBdvwms8iI3u_yqpQxPfAm5bLaJBc' -d '{"access_count":3}'`



