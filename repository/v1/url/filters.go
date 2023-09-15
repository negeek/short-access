package url
// planning to create a function that will dynamically filter the url table. So i can hit the urls endpoint
// and filter by is_custom or by id or by date_created and so on
// if successful it will be the foundational code for other tables

// dynamically construct the strings.
// also think of those fields that are not strings

func UrlFilter(queryParams map[string]{}interface, url *Url)([]Url, error){
	//make sure data type of params correspond to the fields datatype in url table
	//construct query
	//return query
    structType := reflect.TypeOf(url)
	var queryValues []interface{}
	// pre-construct query
	queryValues = append(queryValues, url.UserId)
	query:="SELECT id, original_url,short_url,is_custom,date_created, date_updated FROM urls WHERE user_id=$1 "
    for key, values := range queryParams {
		// complete query
		query+="and "+key+"=$" + strconv.Itoa(len(queryValues)+1)
	
		// convert params type to corresponding url table field type.
        convertedValue, err := ConvertToFieldType(values[0], structType, key)
        if err != nil {
            // Handle the error
            err=fmt.Printf("Error converting parameter to right type %s: %v\n", key, err)
			return nil, err
    	}
		queryValues = append(queryValues, convertedValue)
	}
	fmt.Println("query: ",query)
	rows, err := db.PostgreSQLDB.Query(context.Background(), query, queryValues...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var userUrls []Url
	for rows.Next() {
		var url Url
		err := rows.Scan(&url.Id, &url.Url, &url.ShortUrl, &url.IsCustom, &url.DateCreated, &url.DateUpdated)
		if err != nil {
			return nil, err
		}
		userUrls = append(userUrls, url)
	}
	return userUrls,nil
}