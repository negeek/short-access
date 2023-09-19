package utils
import (
    "fmt"
    "reflect"
    "strconv"
	"github.com/google/uuid"
)
func ShortAccess(quotient int, resultLength int)string{
	numMap:=map[int]string{
		0:"0",1:"1",2:"2",3:"3",4:"4",5:"5",6:"6",7:"7",8:"8",9:"9",
		10:"A",11:"B",12:"C",13:"D",14:"E",15:"F",16:"G",17:"H",18:"I",
		19:"J",20:"K",21:"L",22:"M",23:"N",24:"O",25:"P",26:"Q",27:"R",
		28:"S",29:"T",30:"U",31:"V",32:"W",33:"X",34:"Y",35:"Z", 36:"a",
		37:"b",38:"c",39:"d",40:"e",41:"f",42:"g",43:"h",44:"i",45:"j",
		46:"k",47:"l",48:"m",49:"n",50:"o",51:"p",52:"q",53:"r",54:"s",
		55:"t",56:"u",57:"v",58:"w",59:"x",60:"y",61:"z",
	}

	resStr:=""
	var rem int

	// perform conversion and add to resStr
	for{
		quotient,rem= quotient/62, quotient%62
		resStr+=numMap[rem]
		if quotient<1{
			break
		}
	}

	// reverse the resStr,that is the correct result
	resArr := []byte(resStr)
    for i, j := 0, len(resArr)-1; i < j; i, j = i+1, j-1 {
        resArr[i], resArr[j] = resArr[j], resArr[i]
    }
	resStr=string(resArr)

	//pad the string if length is less than resultLength
	if len(resStr)< resultLength{
		num_zeros:=resultLength-len(resStr)
		zeros:=""
		for i:=0; i<num_zeros;i++{
			zeros+=string("0")
		}
		res:= zeros+resStr
		return res
	}
    return resStr
}

func ConvertToFieldType(value string, structType reflect.Type, key string) (interface{}, error) {
    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        jsonTag := field.Tag.Get("json")
		fmt.Println("tags: ",jsonTag)
        
        // Check if the JSON tag matches the key from the query parameters
        if jsonTag == key {
            fieldType := field.Type
			fmt.Println("fieldtype: ",fieldType)

            // Check the field type and perform the appropriate conversion
            switch fieldType.Kind() {
            case reflect.Int:
                intValue, err := strconv.Atoi(value)
                if err != nil {
                    return nil, err
                }
                return intValue, nil

            case reflect.Bool:
                boolValue, err := strconv.ParseBool(value)
                if err != nil {
                    return nil, err
                }
                return boolValue, nil

			case reflect.String:
				if fieldType == reflect.TypeOf(uuid.UUID{}) {
					uuidValue, err := uuid.Parse(value)
					if err != nil {
						return nil, err
					}
					return uuidValue, nil
				}
				return value, nil
		

            //I should Add more cases for other data types as needed

            default:
                // If the field type is not supported, return an error
                return nil, fmt.Errorf("Unsupported field type: %v", fieldType)
            }
        }
    }

    // If no matching field 
    return nil, fmt.Errorf("Field with JSON tag %s not found", key)
}