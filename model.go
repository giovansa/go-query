package main

const(
	MODE_SELECT = "select"
	MODE_INSERT = "insert"
	MODE_UPDATE = "update"
	MODE_DELETE = "delete"
)

type(
	QueryOptions struct {
		QueryType string
		WhereClause string
	}

	queryTemplate struct {
		options QueryOptions
	}

	baseModel struct {
		key []string
		value []interface{}
	}
)