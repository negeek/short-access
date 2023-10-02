package url

import(
	//"fmt"
	"context"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/db"
	"github.com/jackc/pgx/v4"
	"time"
)

func (u *Url) Create() error {
	query,queryValues,err:=utils.CRUDQueryBuild(u,u.TableName(),"create")
	if err != nil {
		return err
	}
	query+="RETURNING id"
	err2 := db.PostgreSQLDB.QueryRow(context.Background(), query, queryValues...).Scan(&u.Id)
	if err2 != nil {
		return err2
	}
	return nil
}

func (u *Url) Update() error{
	// dyanmically update url table
	query,queryValues,err:=utils.CRUDQueryBuild(u,u.TableName(),"update")
	if err != nil {
		return err
	}
	_,err2 := db.PostgreSQLDB.Exec(context.Background(), query, queryValues...)
	if err2 != nil {
		return err2
	}
	return nil
}

func (u *Url) Delete() error {
	query,queryValues,err:=utils.CRUDQueryBuild(u,u.TableName(),"delete")
	if err != nil {
		return err
	}
	_,err2 := db.PostgreSQLDB.Exec(context.Background(), query, queryValues...)
	if err2 != nil {
		return err2
	}
	return nil
	
}

func(u *Url) FindById()(error,bool){
	query,queryValues,err:=utils.CRUDQueryBuild(u,u.TableName(),"retrieve")
	if err != nil {
		return err, false
	}
	err2:=db.PostgreSQLDB.QueryRow(context.Background(), query,queryValues...).Scan(&u.Id, &u.OriginalUrl, &u.ShortUrl,&u.ShortAccess, &u.IsCustom, &u.AccessCount, &u.ExpireAt, &u.DateCreated, &u.DateUpdated)
	if err2 != nil {
		if err2 == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *Url) FindByOriginalUrl()(error,bool){
	query:="SELECT id,original_url,short_url,short_access,is_custom,access_count,expire_at,date_created,date_updated FROM urls WHERE original_url=$1 and user_id=$2"
	err:=db.PostgreSQLDB.QueryRow(context.Background(), query, u.OriginalUrl, u.UserId).Scan(&u.Id, &u.OriginalUrl, &u.ShortUrl,&u.ShortAccess, &u.IsCustom, &u.AccessCount, &u.ExpireAt, &u.DateCreated, &u.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *Url) FindByShortUrl()(error,bool){
	query:="SELECT id,original_url,short_url,short_access,is_custom,access_count,expire_at,date_created,date_updated FROM urls WHERE short_url=$1"
	err:=db.PostgreSQLDB.QueryRow(context.Background(), query, u.ShortUrl).Scan(&u.Id, &u.OriginalUrl, &u.ShortUrl,&u.ShortAccess, &u.IsCustom, &u.AccessCount, &u.ExpireAt, &u.DateCreated, &u.DateUpdated)
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
		err := rows.Scan(&url.Id, &url.OriginalUrl, &url.ShortUrl,&url.ShortAccess, &url.IsCustom, &url.AccessCount, &url.ExpireAt, &url.DateCreated, &url.DateUpdated)
		if err != nil {
			return nil, err
		}
		userUrls = append(userUrls, url)
	}
	return userUrls,nil
}


func (u *Url) TestDelete() error {
	// for test purpose only
	if u.ShortUrl!=""{
		query:="DELETE FROM urls WHERE short_url=$1"
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, u.ShortUrl)
		if err != nil {
			return err
		}
		
	}else{
		query:="DELETE FROM urls WHERE original_url=$1"
		_, err := db.PostgreSQLDB.Exec(context.Background(), query, u.OriginalUrl)
		if err != nil {
			return err
		}
	}
	return nil	
}

func (u *Url) Expired() bool{
	if u.ExpireAt.IsZero(){
		return false
	}
	return u.ExpireAt.Before(time.Now().UTC())
}