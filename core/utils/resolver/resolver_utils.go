package resolverUtils

import (
	"TestGoLandProject/core/auth"
	commonUtils "TestGoLandProject/core/utils/common"
	"TestGoLandProject/core/utils/database_utils"
	"TestGoLandProject/global_const"
	"TestGoLandProject/graph/model"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// endregion

func DeleteAuthToken(database sq.StatementBuilderType, token string) error {
	_, err := database.Delete("Auth_data").Where(sq.Eq{"token": token}).Exec()
	if err != nil {
		return fmt.Errorf("can't delete token %s frtom database: %w", token, err)
	}

	return nil
}

func CheckUserAuthFromContext(ginContext *gin.Context, database sq.StatementBuilderType) error {
	token, err := auth.GetTokenFromGinContext(ginContext)
	if err != nil {
		return fmt.Errorf("no auth token: %w", err)
	}

	scanRow := database.Select("token").Where(sq.Eq{"token": token}).QueryRow()
	err = scanRow.Scan(new(string))
	if err != nil {
		return fmt.Errorf("invalid auth token: %w", err)
	}

	err = auth.CheckTokenValidity(*token)
	if err != nil {
		err = DeleteAuthToken(database, *token)

		if err != nil {
			return fmt.Errorf("can't delete invalide auth token: %w", err)
		}

		return fmt.Errorf("invalid auth token: %w", err)
	}

	return nil
}

// region User

func GetUserIdFromContext(ginContext *gin.Context, database sq.StatementBuilderType) (*string, error) {
	token, err := auth.GetTokenFromGinContext(ginContext)
	if err != nil {
		return nil, fmt.Errorf("no auth token: %w", err)
	}

	userId, err := auth.GetUserIdByToken(*token)
	if err != nil {
		err = DeleteAuthToken(database, *token)

		if err != nil {
			return nil, fmt.Errorf("can't delete invalide auth token: %w", err)
		}

		return nil, fmt.Errorf("invalid auth token: %w", err)
	}

	return userId, nil
}

func GetUserFriends(database sq.StatementBuilderType, userId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.User, error) {
	friends, err := databaseUtils.GetUserFriendsFromDatabase(database, userId, paginationSettings)
	if err != nil {
		return nil, fmt.Errorf("can't get user %s friends: %w", userId, err)
	}

	return friends, nil
}

func GetUserFriendRequests(database sq.StatementBuilderType, userId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.User, error) {
	friendRequests, err := databaseUtils.GetUserFriendRequestsFromDatabase(database, userId, paginationSettings)
	if err != nil {
		return nil, fmt.Errorf("can't get user %s friend requests: %w", userId, err)
	}

	return friendRequests, nil
}

func GetUserFriendInvites(database sq.StatementBuilderType, userId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.User, error) {
	friendInvites, err := databaseUtils.GetUserFriendInvitesFromDatabase(database, userId, paginationSettings)
	if err != nil {
		return nil, fmt.Errorf("can't get user %s friend invites: %w", userId, err)
	}

	return friendInvites, nil
}

func UpdateUserLastActionTime(database sq.StatementBuilderType, userId string) error {
	_, err := database.Update("User").Where(sq.Eq{"id": userId}).Set("last_action_at", time.Now().UnixMilli()).Exec()
	if err != nil {
		return fmt.Errorf("can't update user last action time in database: %w", err)
	}

	return nil
}

func RefreshUserToken(database sq.StatementBuilderType, userId string) (*string, error) {
	jwtToken, err := auth.GenerateJwt(userId, globalConst.TokenLiveTime)
	if err != nil {
		return nil, fmt.Errorf("can't generate token: %w", err)
	}

	_, err = database.Update("Auth_data").Where(sq.Eq{"user_id": userId}).Set("token", jwtToken).Set("expired_at", time.Now().Add(globalConst.TokenLiveTime).UnixMilli()).Exec()
	if err != nil {
		return nil, fmt.Errorf("can't update user auth_data in database: %w", err)
	}

	return &jwtToken, nil
}

// endregion
