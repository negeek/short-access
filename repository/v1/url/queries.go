package url

import(
	//"fmt"
	"context"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/db"
	"github.com/jackc/pgx/v4"
)

func (u *Url) Create() error {
	utils.Time(u,true)
	query:="INSERT INTO urls (user_id, original_url, short_url, is_custom, access_count, date_created, date_updated) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_,err := db.PostgreSQLDB.Exec(context.Background(), query, u.UserId, u.OriginalUrl, u.ShortUrl, u.IsCustom, u.AccessCount, u.DateCreated, u.DateUpdated)
	if err != nil {
		return err
	}
	return nil
}

func (u *Url) Update() error{
	// dyanmically update url table
}

func (u *Url) UpdateAccessCount() error {
	utils.Time(u,false)
	query:="UPDATE urls SET access_count = $1 WHERE user_id=$2 and short_url=$3"
	_,err := db.PostgreSQLDB.Exec(context.Background(), query, u.AccessCount, u.UserId, u.ShortUrl)
	if err != nil {
		return err
	}
	return nil


}

func (u *Url) FindByOriginalUrl()(error,bool){
	query:="SELECT short_url FROM urls WHERE original_url=$1 and user_id=$2"
	err:=db.PostgreSQLDB.QueryRow(context.Background(), query, u.OriginalUrl, u.UserId).Scan(&u.ShortUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *Url) FindByShortUrl()(error,bool){
	query:="SELECT user_id, original_url, access_count FROM urls WHERE short_url=$1"
	err:=db.PostgreSQLDB.QueryRow(context.Background(), query, u.ShortUrl).Scan(&u.UserId, &u.OriginalUrl, &u.AccessCount)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *Url) UserUrls(query string, queryValues []interface{})([]Url,error){
	rows,err:=db.PostgreSQLDB.Query(context.Background(), query, queryValues...)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	var userUrls []Url
	for rows.Next() {
		var url Url
		err := rows.Scan(&url.Id, &url.OriginalUrl, &url.ShortUrl, &url.IsCustom, &url.AccessCount, &url.DateCreated, &url.DateUpdated)
		if err != nil {
			return nil, err
		}
		userUrls = append(userUrls, url)
	}
	return userUrls,nil
}

func (u *Url) Delete() error {
	if u.ShortUrl!=""{
		query:="DELETE FROM urls WHERE short_url=$1"
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, u.ShortUrl)
		if err != nil {
			return err
		}
		
	}else{
		// this will delete every instance of the url which will affect other users, so this is for test only
		query:="DELETE FROM urls WHERE original_url=$1"
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, u.OriginalUrl)
		if err != nil {
			return err
		}
	}
	return nil
	
}
