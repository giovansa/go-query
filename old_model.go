package main

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

type batchStructModel struct {
	values []structModel
	err error
}
