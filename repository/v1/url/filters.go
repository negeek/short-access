package url
import (
    //"fmt"
	"context"
    "reflect"
    "strconv"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/db"
)

func UrlFilter(queryParams map[string][]string, url Url)([]Url, error){
	// dynamically filter url table
    structType := reflect.TypeOf(url)
	var queryValues []interface{}
	// pre-construct query
	queryValues = append(queryValues, url.UserId)
	query:="SELECT id, original_url,short_url,is_custom,date_created,date_updated FROM urls WHERE user_id=$1"
    for key, values := range queryParams {
		// complete query
		query+=" and "+key+"=$" + strconv.Itoa(len(queryValues)+1)
	
		// convert params type to corresponding url table field type.
        convertedValue, err := utils.ConvertToFieldType(values[0], structType, key)
        if err != nil {
			return nil, err
    	}
		queryValues = append(queryValues, convertedValue)
	}
	rows, err := db.PostgreSQLDB.Query(context.Background(), query, queryValues...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var userUrls []Url
	for rows.Next() {
		var url Url
		err := rows.Scan(&url.Id, &url.OriginalUrl, &url.ShortUrl, &url.IsCustom, &url.DateCreated, &url.DateUpdated)
		if err != nil {
			return nil, err
		}
		userUrls = append(userUrls, url)
	}
	return userUrls,nil
}