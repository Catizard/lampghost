package data

import "github.com/guregu/null/v5"

type Filter interface {
	GenerateWhereClause() string
}

type nullFilter struct {
}

func (f nullFilter) GenerateWhereClause() string {
	return ""
}

func newNullFilter() Filter {
	return nullFilter{}
}

var NullFilter null.Value[Filter] = null.NewValue(newNullFilter(), false)

type SimpleFilter struct {
	WhereClause string
}

func (f SimpleFilter) GenerateWhereClause() string {
	return f.WhereClause
}