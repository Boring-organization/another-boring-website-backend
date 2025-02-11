package resolverUtils

import (
	"TestGoLandProject/core/auth"
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/utils"
	"TestGoLandProject/graph/model"
	"database/sql"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
)

func DeleteAuthToken(token string, database database.Database) error {
	var p []byte
	scanRow := database.QueryRow("select * from Auth_data where token = $1", token)
	err := scanRow.Scan(&p, &p, &p, &p, &p)

	if err != nil {
		return err
	}

	_, err = database.ExecuteOperation("delete from Auth_data where token = $1", token)
	if err != nil {
		return err
	}

	return nil
}

func CheckUserAuthFromContext(ctx context.Context, database database.Database) error {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	token := auth.GetTokenFromGinContext(ginContext)
	if token == nil {
		return utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("no auth token: %w", err))
	}

	var p []byte
	scanRow := database.QueryRow("select * from Auth_data where token = $1", token)
	err = scanRow.Scan(&p, &p, &p, &p, &p)
	if err != nil {
		return utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("invalid auth token: %w", err))
	}

	err = auth.CheckTokenValidity(*token)
	if err != nil {
		err = DeleteAuthToken(*token, database)

		if err != nil {
			return utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("invalid auth token: %w", err))
		}

		return utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("invalid auth token: %w", err))
	}

	return nil
}

// region User
func GetUserFromDatabase(database database.Database, userPointer *model.User, userId string) error {
	userRow := database.QueryRow("select * from User where id = $1", userId)

	err := userRow.Scan(&userPointer.ID, &userPointer.Nickname, &userPointer.Login, &userPointer.Password, &userPointer.CreatedAt, &userPointer.EditedAt, &userPointer.DeletedAt, &userPointer.IsAdmin, &userPointer.LastActionAt)
	if err != nil {
		return err
	}

	userPointer.Password = nil

	return nil
}

func GetUsersFromDatabase(users *[]*model.User, rows *sql.Rows) error {
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Nickname, &user.Login, &user.Password, &user.CreatedAt, &user.EditedAt, &user.DeletedAt, &user.IsAdmin, &user.LastActionAt)

		if err != nil {
			return err
		}

		user.Password = nil

		*users = append(*users, &user)
	}

	return nil
}

func GetUserIdFromContext(ctx context.Context, database database.Database) (*string, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	token := auth.GetTokenFromGinContext(ginContext)
	if token == nil {
		return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("no auth token: %w", err))
	}

	userId, err := auth.GetUserIdByToken(*token)
	if err != nil {
		err = DeleteAuthToken(*token, database)

		if err != nil {
			return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("invalid auth token, tried to delete auth token, can't delete: %w", err))
		}

		return nil, utils.ResponseError(ginContext, http.StatusUnauthorized, fmt.Errorf("can't get user id from token: %w", err))
	}

	return userId, nil
}

func GetUserFriends(ctx context.Context, userId string, database database.Database) ([]*model.User, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	rows, err := database.Query("select u.*  from User_friend_link f1 join User_friend_link f2 on f1.requester_id = f2.requested_id and f1.requested_id = f2.requester_id join User u on f1.requested_id = u.id where f1.requester_id = $1 and f1.requested_id != $1", userId)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get user %s friends: %w", userId, err))
	}

	friends := make([]*model.User, 0)

	err = GetUsersFromDatabase(&friends, rows)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get user %s friends: %w", userId, err))
	}

	return friends, nil
}

func GetUserFriendRequests(ctx context.Context, userId string, database database.Database) ([]*model.User, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	rows, err := database.Query("select u.* from User_friend_link f1 left join User_friend_link f2 on f1.requester_id = f2.requested_id and f1.requested_id = f2.requester_id join User u on f1.requested_id = u.id where f1.requester_id = $1 and f2.requester_id is null", userId)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get user %s friends: %w", userId, err))
	}

	friends := make([]*model.User, 0)

	err = GetUsersFromDatabase(&friends, rows)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get user %s friends: %w", userId, err))
	}

	return friends, nil
}

func GetUserFriendInvites(ctx context.Context, userId string, database database.Database) ([]*model.User, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	rows, err := database.Query("select u.* from User_friend_link f1 left join User_friend_link f2 on f1.requester_id = f2.requested_id and f1.requested_id = f2.requester_id join User u on f1.requester_id = u.id where f1.requested_id = $1 and f2.requester_id is null", userId)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("Can't get user %s friends: %w", userId, err))
	}

	friends := make([]*model.User, 0)

	err = GetUsersFromDatabase(&friends, rows)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("Can't get user %s friends: %w", userId, err))
	}

	return friends, nil
}

//endregion
