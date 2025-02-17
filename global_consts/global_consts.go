package global_consts

import "TestGoLandProject/graph/model"

var (
	ApiUrl                        = "/api"
	QueryUrl                      = ApiUrl + "/query"
	JwtSeed                       = "my_secret_code"
	DefaultListPaginationSettings = model.ListByTimeSortPaginationSettings{
		Count:             10,
		LastItemCreatedAt: 0,
		Ascending:         true,
	}
)
