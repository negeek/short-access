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
	query:="INSERT INTO urls (user_id, original_url, short_url, is_custom, date_created, date_updated) VALUES ($1, $2, $3, $4, $5, $6)"
	_,err := db.PostgreSQLDB.Exec(context.Background(), query, u.UserId, u.Url, u.ShortUrl, u.IsCustom, u.DateCreated, u.DateUpdated)
	if err != nil {
		return err
	}
	return nil
}

func (u *Url) FindByOriginalUrl()(error,bool){
	query:="SELECT short_url FROM urls WHERE original_url=$1 and user_id=$2"
	err:=db.PostgreSQLDB.QueryRow(context.Background(), query, u.Url, u.UserId).Scan(&u.ShortUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *Url) FindByShortUrl()(error,bool){
	query:="SELECT user_id, original_url FROM urls WHERE short_url=$1"
	err:=db.PostgreSQLDB.QueryRow(context.Background(), query, u.ShortUrl).Scan(&u.UserId, &u.Url)
	if err != nil {
		if err == pgx.ErrNoRows{
			return nil, false
		}
		return err,false
	}
	return nil, true
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
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, u.Url)
		if err != nil {
			return err
		}
	}
	return nil
	
}