package utils
import (
	//"fmt"
	"errors"
	"reflect"
	"time"
)
//function to help with datecreated and dateupdated fields of db tables.
func Time(strct interface{}, new ...bool) error {
	t := reflect.TypeOf(strct)
	v := reflect.ValueOf(strct).Elem()
	// validate if strct is a pointer and struct
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("strct must be a pointer to a struct")
	}

	// validate if the datecreated and dateupdated fields are in strct and are of type time.Time
	dateCreatedField, has_created := t.Elem().FieldByName("DateCreated")
	dateUpdatedField, has_updated := t.Elem().FieldByName("DateUpdated")
	if has_created==false || has_updated==false {
		return errors.New("strct must have DateCreated and DateUpdated fields")
	}
	
	if dateCreatedField.Type.Kind() != reflect.TypeOf(time.Time{}).Kind() || dateUpdatedField.Type.Kind() != reflect.TypeOf(time.Time{}).Kind() {
		return errors.New("strct DateCreated and DateUpdated fields must be of type time.Time")
	}

	// set the time for the fields based on new arguement value
	if len(new) > 0 && new[0] {
		// Set the "DateUpdated" field to current UTC time
		v.FieldByName("DateCreated").Set(reflect.ValueOf(time.Now().UTC()))
	}
	v.FieldByName("DateUpdated").Set(reflect.ValueOf(time.Now().UTC()))
	return nil
	
}

func ExpiryDateTime(timeUnit string, timeValue int) (time.Time, error){
	/*set expiry datetime based on time unit and time value. useful for setting expiry_at for url table.*/
	start:=time.Now().UTC()
	var expire_at time.Time
	switch timeUnit {
	case "y":
		expire_at= start.AddDate(timeValue, 0, 0)
		return expire_at, nil
	case "mo":
		expire_at= start.AddDate(0, timeValue, 0)
		return expire_at, nil
	case "d":
		expire_at= start.AddDate(0, 0, timeValue)
		return expire_at, nil
	case "h":
		expire_at= start.Add(time.Duration(timeValue) * time.Hour) 
		return expire_at, nil
	case "m":
		expire_at= start.Add(time.Duration(timeValue) * time.Minute)
		return expire_at, nil
	case "s":
		expire_at= start.Add(time.Duration(timeValue) * time.Second)
		return expire_at, nil
	default:
		return time.Time{}, errors.New("Unsupported time unit. Use any of the following ['y','mo','d','h','m','s']. Which denotes year, month, day, hour, minute, second.")
	}
		
}
