package main

import (
	"reflect"
	"strings"
)

type Query interface {
	GetQuery(data interface{})(string, error)
	GetData(data interface{})(interface{}, error)
}

func SetParam(queryType, whereClause string)QueryOptions{
	return QueryOptions{
			QueryType:queryType,
			WhereClause:whereClause}
}

func NewQuery(options QueryOptions)Query{
	return &queryTemplate{
		options:options,
	}
}

func (q *queryTemplate)GetQuery(data interface{})(output string, err error){
	switch strings.ToLower(q.options.QueryType) {
	case MODE_SELECT:
	case MODE_INSERT:
	case MODE_UPDATE:
	case MODE_DELETE:
	}
	return q.options.QueryType, nil
}

func (q *queryTemplate)GetData(data interface{})(output interface{}, err error){
	return q.options.QueryType, nil
}

func constructBaseModel(data interface{}) baseModel{

	var keys []string
	var vals []interface{}

	typeReflect := reflect.TypeOf(data)
	valReflect := reflect.ValueOf(data)
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
	return baseModel{key: keys,value: vals}
}