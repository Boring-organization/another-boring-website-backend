package globalConst

import (
	"TestGoLandProject/graph/model"
	"time"
)

var (
	ApiUrl                        = "/api"
	QueryUrl                      = ApiUrl + "/query"
	JwtSeed                       = "my_secret_code"
	TokenLiveTime                 = time.Hour * 24 * 30
	DefaultListPaginationSettings = model.ListByTimeSortPaginationSettings{
		Count:             10,
		LastItemCreatedAt: 0,
		Ascending:         true,
	}
)
