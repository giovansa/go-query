package main

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

//Interface to provide functions to generate query
type providerQuery interface {
	ViewAll(table string)(query string , err error)
	Insert(table string)(query string, values []interface{}, err error)
	Delete(table string)(query string, err error)
	Update(table string)(query string, values []interface{}, err error)
	Where(operatorCondition, operationBetweenCondition string)(query string, values []interface{}, err error)
}

/*
	Defining the body function
*/
//Function to generating query for query SELECT *
func (s structModel)ViewAll(table string)(query string, err error){
	if s.err != nil{
		return "", s.err
	}
	var arrQuery []string
	query = "SELECT"
	for i, _ := range s.value{
		arrQuery = append(arrQuery, " " +s.key[i])
	}
	query = query + strings.Join(arrQuery, ",") + " FROM " + table
	return query,nil
}

//Function to generating query for query INSERT
func (s structModel)Insert(table string)(query string, values []interface{}, err error){
	if s.err != nil{
		return "", nil, s.err
	}
	var arrQuery, valArr []string
	query = "INSERT INTO " + table +"("
	queryForValues := " VALUES("
	listValue := make([]interface{}, 0)

	for i, _ := range s.value{
		arrQuery = append(arrQuery, " " + s.key[i])
		valArr = append(valArr, " $" + strconv.Itoa(i+1))
		listValue = append(listValue, s.value[i])
	}
	query = query + strings.Join(arrQuery, ",") + ")"
	queryForValues = queryForValues + strings.Join(valArr, ",") + ")"
	return query + queryForValues, listValue, nil
}

//Function to generating query for query DELETE
func (s structModel)Delete(table string)(query string, err error){
	if s.err != nil{
		return "", s.err
	}

	query = "DELETE FROM " + table
	return query, nil
}

//Function to generating query for query UPDATE
func (s structModel)Update(table string)(query string, values[]interface{}, err error){
	if s.err != nil{
		return "", nil, s.err
	}
	var arrQuery []string
	query = "UPDATE " + table + " SET"
	listValues := make([]interface{}, 0)
	for i, _ := range s.value{
		arrQuery = append(arrQuery, " " + s.key[i] + "= $" + strconv.Itoa(i+1))
		listValues = append(listValues, s.value[i])
	}
	query = query + strings.Join(arrQuery, ",")
	return query, listValues, nil
}

//Function to generating WHERE condition to the query
func (s structModel)Where(operatorCondition, operationBetweenCondition string)(query string, values[]interface{}, err error){
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



type batchQuery interface {
	InsertQuery(table string)(query string, err error)
	ValueBatch()(query string, values []interface{}, err error)
}

func (batchStructModel)InsertQuery(table string)(query string, err error){

	return "", nil
}
func (batchStructModel)ValueBatch()(query string, values []interface{}, err error){

	return "", nil, nil
}


func ValueConversion(model interface{}) batchQuery{
	if reflect.TypeOf(model).Kind() == reflect.Slice{
		/*
			Check if inside the slice is a struct
		*/
		value := reflect.TypeOf(model).Elem()
		if reflect.TypeOf(value.Field(0)).Kind() != reflect.Struct{
			return batchQuery(batchStructModel{err:errors.New("parameter must be a struct")})
		}
	}
	convertedModel := batchStructModel{}
	result := batchQuery(convertedModel)
	return result
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
		/*
			Skipping for nested struct
		*/
		if typField.Type.Kind() == reflect.Struct{
			continue
		}
		keyValue, ok := typField.Tag.Lookup("db")
		if !ok {
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

		dateVal, ok := typField.Tag.Lookup("date")
		if ok && dateVal=="now"{
			keys = append(keys, keyValue)
			vals = append(vals, "now()")
			continue
		}
		keys = append(keys, keyValue)
		vals = append(vals, valueField.Interface().(interface{}))
	}
	convertedModel := structModel{key: keys,value: vals, err: nil}
	result := providerQuery(convertedModel)
	return result
}
