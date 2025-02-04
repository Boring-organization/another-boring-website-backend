package resolverUtils

import (
	"TestGoLandProject/core/auth"
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/utils"
	"golang.org/x/net/context"
	"net/http"
)

func DeleteAuthToken(token string, database database.Database) {
	var p []byte
	scanRow := database.QueryRow("select * from Auth_data where token = $1", token)
	err := scanRow.Scan(&p, &p, &p, &p, &p)

	if err != nil {
		return
	}

	database.ExecuteOperation("delete from Auth_data where token = $1", token)
}

func GetUserIdFromContext(ctx context.Context, database database.Database) (*string, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	token := auth.GetTokenFromGinContext(ginContext)
	if token == nil {
		return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, "No auth token")
	}

	var p []byte
	scanRow := database.QueryRow("select * from Auth_data where token = $1", token)
	err = scanRow.Scan(&p, &p, &p, &p, &p)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, "Invalid auth token")
	}

	err = auth.CheckTokenValidity(*token)
	if err != nil {
		DeleteAuthToken(*token, database)
		return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, "Invalid auth token")
	}

	userId, err := auth.GetUserIdByToken(*token)
	if err != nil {
		DeleteAuthToken(*token, database)
		return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, "Can't get user id from token")
	}

	return userId, nil
}

func CheckUserAuthFromContext(ctx context.Context, database database.Database) error {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	token := auth.GetTokenFromGinContext(ginContext)
	if token == nil {
		return utils.ResponseError(ginContext, http.StatusUnauthorized, "No auth token")
	}

	var p []byte
	scanRow := database.QueryRow("select * from Auth_data where token = $1", token)
	err = scanRow.Scan(&p, &p, &p, &p, &p)
	if err != nil {
		return utils.ResponseError(ginContext, http.StatusUnauthorized, "Invalid auth token")
	}

	err = auth.CheckTokenValidity(*token)
	if err != nil {
		DeleteAuthToken(*token, database)
		return utils.ResponseError(ginContext, http.StatusUnauthorized, "Invalid auth token")
	}

	return nil
}
