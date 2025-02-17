package resolver

import (
	sq "github.com/Masterminds/squirrel"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db sq.StatementBuilderType
}
