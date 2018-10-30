package query

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

//Interface to provide functions to generate query
type providerQuery interface {
	viewAll(table string)(query string , err error)
	insert(table string)(query string, values []interface{}, err error)
	delete(table string)(query string, err error)
	update(table string)(query string, values []interface{}, err error)
	where(operatorCondition, operationBetweenCondition string)(query string, values []interface{}, err error)
}

/*
	Defining the body function
*/
//Function to generating query for query SELECT *
func (s structModel)viewAll(table string)(query string, err error){
	if s.err != nil{
		return "", s.err
	}

	query = "SELECT"
	for i, _ := range s.value{
		keyValue := s.key[i]
		query += " " +keyValue + ","
	}
	query = query[0:(len(query)-1)]
	query += " FROM " + table
	return query,nil
}

//Function to generating query for query INSERT
func (s structModel)insert(table string)(query string, values []interface{}, err error){
	if s.err != nil{
		return "", nil, s.err
	}

	query = "INSERT INTO " + table +"("
	queryForValues := " VALUES("
	listValue := make([]interface{}, 0)

	for i, _ := range s.value{
		query += " " + s.key[i] + ","
		queryForValues += " $" + strconv.Itoa(i+1) +","
		listValue = append(listValue, s.value[i])
	}
	query = query[0:len(query)-1]
	query += ")"
	queryForValues = queryForValues[0:len(queryForValues)-1]
	queryForValues += ")"

	return query + queryForValues, listValue, nil
}

//Function to generating query for query DELETE
func (s structModel)delete(table string)(query string, err error){
	if s.err != nil{
		return "", s.err
	}

	query = "DELETE FROM " + table
	return query, nil
}

//Function to generating query for query UPDATE
func (s structModel)update(table string)(query string, values[]interface{}, err error){
	if s.err != nil{
		return "", nil, s.err
	}

	query = "UPDATE " + table + " SET"
	listValues := make([]interface{}, 0)
	for i, _ := range s.value{
		query += " " + s.key[i] + "= $" + strconv.Itoa(i+1) + ","
		listValues = append(listValues, s.value[i])
	}
	query = query[0:len(query)-1]
	return query, listValues, nil
}

//Function to generating WHERE condition to the query
func (s structModel)where(operatorCondition, operationBetweenCondition string)(query string, values[]interface{}, err error){
	if s.err != nil{
		return "", nil, s.err
	}
	query = " WHERE"
	listValues := make([]interface{}, 0)
	for i, _ := range s.value{
		query += " " + s.key[i] + " " + operatorCondition + " " + "$" + strconv.Itoa(i+1) + " " + operationBetweenCondition
		listValues = append(listValues, s.value[i])
	}
	query = query[0:(len(query)-len(operationBetweenCondition))]
	return  query, listValues, nil
}

/*
	Model that is used as key and value to create the query
	key : database column name
	value : expected value
*/
type structModel struct {
	key   []string
	value []interface{}
	err error
}

func Conversion(model interface{}) providerQuery {
	/*
		Returning error as validation to force function only accepting struct
		struct assume as reflect.Struct
	*/
	if reflect.TypeOf(model).Kind() != reflect.Struct{
		return providerQuery(structModel{err:errors.New("parameter must be a struct")})
	}

	var keys []string
	var vals []interface{}

	typeReflect := reflect.TypeOf(model)
	valReflect := reflect.ValueOf(model)
	/*
		Loop through the model to convert it to other model ('structModel')
		to be treated as key and value
	*/
	for i := 0; i < typeReflect.NumField(); i++ {
		typField := typeReflect.Field(i)
		valueField := valReflect.Field(i)
		/*
			Skip iteration if data is empty
			now empty is considered as empty string ("") or 0 if the data type is integer
		*/
		if valueField.String() == "" {
			continue
		}
		if typField.Type.String() == "int" || typField.Type.String() == "int64"{
			if valueField.Int() == 0{continue}
		}
		keyValue, ok := typField.Tag.Lookup("db")
		if !ok {
			/*
				If tag not found
				the converion will search for default tag
				the default tag wll be used, to manipulate string in the struct's attribute
				with ToUpper or ToLower
				thus the Tag should be "lower" or "upper"
				other than that struct's attribute name will be used to indentify database column
			*/
			tagDefault, ok := typField.Tag.Lookup("default")
			if !ok{
				keyValue = typField.Name
			}
			switch tagDefault {
			case "lower":
				keyValue = strings.ToLower(typField.Name)
			case "upper":
				keyValue = strings.ToUpper(typField.Name)
			default:
				keyValue = typField.Name
			}
		}

		/*
			This part is used to check the data type to be converted correctly to be used in query
			--------------------------------------------------------------------------------------
			For temporary the assumption is int, int64 and bool to be treated differently
		*/
		if valueField.Type().String() == "int" || valueField.Type().String() == "int64" {
			newInt := strconv.Itoa(int(valueField.Int()))
			keys = append(keys, keyValue)
			vals = append(vals, newInt)
		} else if valueField.Type().String() == "bool" {
			keys = append(keys, keyValue)
			vals = append(vals, strconv.FormatBool(valueField.Bool()))
		} else {
			keys = append(keys, keyValue)
			vals = append(vals, valueField.String())
		}
	}
	convertedModel := structModel{key: keys,value: vals, err: nil}
	result := providerQuery(convertedModel)
	return result
}
