package databaseUtils

import (
	"TestGoLandProject/global_const"
	"TestGoLandProject/graph/model"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"
)

// region Pagination

func AddPaginationByCreatedAtToQuery(selectRequest sq.SelectBuilder, paginationSettings model.ListByTimeSortPaginationSettings) sq.SelectBuilder {
	withPagination := selectRequest.Limit(uint64(paginationSettings.Count))

	if paginationSettings.Ascending {
		withPagination = withPagination.Where(sq.Gt{"created_at": paginationSettings.LastItemCreatedAt}).OrderBy("created_at ASC")
	} else {
		withPagination = withPagination.Where(sq.Lt{"created_at": paginationSettings.LastItemCreatedAt}).OrderBy("created_at DESC")
	}

	return withPagination
}

// endregion

// region Image

func GetImageById(imageId string) (*model.Image, error) {
	// TODO Реализовать сохранение изобраения
	return &model.Image{}, nil
}

func GetImagesFromFiles(imageIds *[]string) (*[]*model.Image, error) {
	var images []*model.Image

	for _, imageId := range *imageIds {
		newImage, err := GetImageById(imageId)
		if err != nil {
			return nil, fmt.Errorf("can't get image %s: %w", imageId, err)
		}

		images = append(images, newImage)
	}

	return &images, nil
}

func DeleteImageById(imageId string) error {
	// TODO Реализовать удаление изобраения
	return nil
}

func ClearTemporaryImageTable(database sq.StatementBuilderType) error {
	imageIds, err := getEntityIdsFromDatabase(TemporaryImageFieldsSelect(database).Where(sq.Lt{"created_at": time.Time{}.Add(time.Hour * 24)}))

	if err != nil {
		return fmt.Errorf("can't get temporary imageIds: %w", err)
	}

	for _, imageId := range *imageIds {
		err = DeleteImageById(imageId)
		if err != nil {
			return fmt.Errorf("can't delete temporary image %s from files: %w", imageId, err)
		}

		_, err = database.Delete("Temporary_image").Where(sq.Eq{"image_id": imageId}).Exec()
		if err != nil {
			return fmt.Errorf("can't delete temporary image %s from database: %w", imageId, err)
		}
	}

	return nil
}

// endregion

// region Get entity

func getEntityFromDatabase(selectRequest sq.SelectBuilder, entityId string, entityFields ...any) error {
	entityRow := selectRequest.Where(sq.Eq{"id": entityId}).QueryRow()

	err := entityRow.Scan(entityFields)
	if err != nil {
		return fmt.Errorf("can't scan entity from database: %w", err)
	}

	return nil
}

func getEntityIdsWithPaginationFromDatabase(selectRequest sq.SelectBuilder, paginationSettings model.ListByTimeSortPaginationSettings) (*[]string, error) {
	return getEntityIdsFromDatabase(AddPaginationByCreatedAtToQuery(selectRequest, paginationSettings))
}

func getEntityIdsFromDatabase(selectRequest sq.SelectBuilder) (*[]string, error) {
	var entities []string

	rows, err := selectRequest.Query()
	if err != nil {
		return nil, fmt.Errorf("can't get entities: %w", err)
	}

	for rows.Next() {
		var entity string

		err = rows.Scan(&entity)
		if err != nil {
			return nil, fmt.Errorf("can't scan entity: %w", err)
		}

		entities = append(entities, entity)
	}

	return &entities, nil
}

// endregion

// region User

func GetUserFromDatabase(database sq.StatementBuilderType, userId string) (*model.User, error) {
	user := model.User{}

	err := getEntityFromDatabase(UserFieldsSelect(database), userId, []any{UserFields(&user)})
	if err != nil {
		return nil, err
	}

	return &user, nil
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

func GetUserFriendsFromDatabase(database sq.StatementBuilderType, userId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.User, error) {
	var users []*model.User

	query := UserAsFieldsSelect(database)
	query = query.From("User_friend_link f1")
	query = query.Join("User_friend_link f2 on f1.requester_id = f2.requested_id and f1.requested_id = f2.requester_id")
	query = query.Join("User u on f1.requested_id = u.id")
	query = query.Where(sq.Eq{"f1.requester_id": userId})
	query = query.Where(sq.NotEq{"f1.requested_id": userId})

	rows, err := AddPaginationByCreatedAtToQuery(query, paginationSettings).Query()
	if err != nil {
		return nil, fmt.Errorf("can't get user %s frineds from database: %w", userId, err)
	}

	for rows.Next() {
		user := model.User{}

		err = rows.Scan(UserFields(&user))
		if err != nil {
			return nil, fmt.Errorf("can't scan user %s friend from database: %w", userId, err)
		}

		users = append(users, &user)
	}

	return &users, nil
}

func GetUserFriendRequestsFromDatabase(database sq.StatementBuilderType, userId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.User, error) {
	var users []*model.User

	query := UserAsFieldsSelect(database)
	query = query.From("User_friend_link f1")
	query = query.LeftJoin("User_friend_link f2 on f1.requester_id = f2.requested_id and f1.requested_id = f2.requester_id")
	query = query.Join("User u on f1.requested_id = u.id")
	query = query.Where(sq.Eq{"f1.requester_id": userId})
	query = query.Where("f2.requester_id is null")

	rows, err := AddPaginationByCreatedAtToQuery(query, paginationSettings).Query()
	if err != nil {
		return nil, fmt.Errorf("can't get user %s friend requests from database: %w", userId, err)
	}

	for rows.Next() {
		user := model.User{}

		err = rows.Scan(UserFields(&user))
		if err != nil {
			return nil, fmt.Errorf("can't scan user %s friend request from database: %w", userId, err)
		}

		users = append(users, &user)
	}

	return &users, nil
}

func GetUserFriendInvitesFromDatabase(database sq.StatementBuilderType, userId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.User, error) {
	var users []*model.User

	query := UserAsFieldsSelect(database)
	query = query.From("User_friend_link f1")
	query = query.LeftJoin("User_friend_link f2 on f1.requester_id = f2.requested_id and f1.requested_id = f2.requester_id")
	query = query.Join("User u on f1.requester_id = u.id")
	query = query.Where(sq.Eq{"f1.requested_id": userId})
	query = query.Where("f2.requester_id is null")

	rows, err := AddPaginationByCreatedAtToQuery(query, paginationSettings).Query()
	if err != nil {
		return nil, fmt.Errorf("can't get user %s friend invites from database: %w", userId, err)
	}

	for rows.Next() {
		user := model.User{}

		err = rows.Scan(UserFields(&user))
		if err != nil {
			return nil, fmt.Errorf("can't scan user %s friend invites from database: %w", userId, err)
		}

		users = append(users, &user)
	}

	return &users, nil
}

// endregion

// region Post entities

// region Post

func GetPostFromDatabase(database sq.StatementBuilderType, postId string) (*model.Post, error) {
	post := model.Post{}
	authorId := new(string)

	err := getEntityFromDatabase(PostFieldsSelect(database), postId, []any{PostFields(&post, authorId)})
	if err != nil {
		return nil, fmt.Errorf("can't scan post %s from database: %w", postId, err)
	}

	author, err := GetUserFromDatabase(database, *authorId)
	if err != nil {
		return nil, fmt.Errorf("can't get user for post %s: %w", postId, err)
	}

	comments, err := GetPostCommentsFromDatabase(database, postId, globalConst.DefaultListPaginationSettings)
	if err != nil {
		return nil, fmt.Errorf("can't get user for post %s: %w", postId, err)
	}

	likeCount, err := GetPostLikesCount(database, postId)
	if err != nil {
		return nil, fmt.Errorf("can't get post %s likes: %w", postId, err)
	}

	post.Author = author
	post.LikeCount = *likeCount
	post.Comments = *comments

	return &post, nil
}

// endregion

// region Post comment

func GetPostCommentsFromDatabase(database sq.StatementBuilderType, postId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.PostComment, error) {
	var postComments []*model.PostComment

	rows, err := AddPaginationByCreatedAtToQuery(PostCommentFieldsSelect(database).Where(sq.Eq{"post_id": postId, "deleted_at": nil}), paginationSettings).Query()
	if err != nil {
		return nil, fmt.Errorf("can't get post comment ids: %w", err)
	}

	for rows.Next() {
		var postComment model.PostComment
		authorId := new(string)
		linkedUserId := new(string)

		err = rows.Scan(PostCommentFields(&postComment, authorId, linkedUserId))
		if err != nil {
			return nil, fmt.Errorf("can't scan post comment: %w", err)
		}

		author, err := GetUserFromDatabase(database, *authorId)
		if err != nil {
			return nil, fmt.Errorf("can't get author of post %s: %w", postComment.ID, err)
		}

		postComment.Author = author

		if linkedUserId != nil {
			linkedUser, err := GetUserFromDatabase(database, *authorId)
			if err != nil {
				return nil, fmt.Errorf("can't get linked user of post %s: %w", postComment.ID, err)
			}

			postComment.LinkedUser = linkedUser
		}

		postComments = append(postComments, &postComment)
	}

	return &postComments, nil
}

// endregion

// region Post like

func GetPostLikesCount(database sq.StatementBuilderType, postId string) (*int, error) {
	likesCountString := new(int)
	likeCountRow := database.Select("count(*)").From("Post_like").Where(sq.Eq{"post_id": postId}).QueryRow()

	err := likeCountRow.Scan(&likesCountString)
	if err != nil {
		return nil, fmt.Errorf("can't get like row count by post %s: %w", postId, err)
	}

	return likesCountString, nil
}

// endregion

// endregion

// region Game find post entities

// region Game find post

func GetGameFindPostFromDatabase(database sq.StatementBuilderType, gameFindPostId string) (*model.GameFindPost, error) {
	gameFindPost := model.GameFindPost{}
	authorId := new(string)

	err := getEntityFromDatabase(GameFindPostFieldsSelect(database), gameFindPostId, []any{GameFindPostFields(&gameFindPost, authorId)})
	if err != nil {
		return nil, fmt.Errorf("can't scan game find post %s from database: %w", gameFindPostId, err)
	}

	author, err := GetUserFromDatabase(database, *authorId)
	if err != nil {
		return nil, fmt.Errorf("can't get user for game find post %s: %w", gameFindPostId, err)
	}

	comments, err := GetGameFindPostCommentsFromDatabase(database, gameFindPostId, globalConst.DefaultListPaginationSettings)
	if err != nil {
		return nil, fmt.Errorf("can't get user for game find post %s: %w", gameFindPostId, err)
	}

	likeCount, err := GetGameFindPostLikesCount(database, gameFindPostId)
	if err != nil {
		return nil, fmt.Errorf("can't get game find post %s likes: %w", gameFindPostId, err)
	}

	kinks, err := GetGameFindPostKinks(database, gameFindPostId)
	if err != nil {
		return nil, fmt.Errorf("can't get game find post %s kinks: %w", gameFindPostId, err)
	}

	gameFindPost.Author = author
	gameFindPost.LikeCount = *likeCount
	gameFindPost.Comments = *comments
	gameFindPost.Kinks = *kinks

	return &gameFindPost, nil
}

// endregion

// region Game find post comment

func GetGameFindPostCommentsFromDatabase(database sq.StatementBuilderType, gameFindPostId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.GameFindPostComment, error) {
	var gameFindPostComments []*model.GameFindPostComment

	rows, err := AddPaginationByCreatedAtToQuery(GameFindPostCommentFieldsSelect(database).Where(sq.Eq{"game_find_post_id": gameFindPostId, "deleted_at": nil}), paginationSettings).Query()
	if err != nil {
		return nil, fmt.Errorf("can't get game find post comment ids: %w", err)
	}

	for rows.Next() {
		var gameFindPostComment model.GameFindPostComment
		authorId := new(string)
		linkedUserId := new(string)

		err = rows.Scan(GameFindPostCommentFields(&gameFindPostComment, authorId, linkedUserId))
		if err != nil {
			return nil, fmt.Errorf("can't scan game find post comment: %w", err)
		}

		author, err := GetUserFromDatabase(database, *authorId)
		if err != nil {
			return nil, fmt.Errorf("can't get author of game find post %s: %w", gameFindPostComment.ID, err)
		}

		gameFindPostComment.Author = author

		if linkedUserId != nil {
			linkedUser, err := GetUserFromDatabase(database, *authorId)
			if err != nil {
				return nil, fmt.Errorf("can't get linked user of game find post %s: %w", gameFindPostComment.ID, err)
			}

			gameFindPostComment.LinkedUser = linkedUser
		}

		gameFindPostComments = append(gameFindPostComments, &gameFindPostComment)
	}

	return &gameFindPostComments, nil
}

// endregion

// region Game find post like

func GetGameFindPostLikesCount(database sq.StatementBuilderType, postId string) (*int, error) {
	likesCountString := new(int)
	likeCountRow := database.Select("count(*)").From("Game_find_post_like").Where(sq.Eq{"game_find_post_id": postId}).QueryRow()

	err := likeCountRow.Scan(&likesCountString)
	if err != nil {
		return nil, fmt.Errorf("can't get like row count by game find post %s: %w", postId, err)
	}

	return likesCountString, nil
}

// region Game find post like

func GetGameFindPostKinks(database sq.StatementBuilderType, postId string) (*[]string, error) {
	kinkCodes, err := getEntityIdsFromDatabase(GameFindPostKinkFieldsSelect(database).Where(sq.Eq{"post_id": postId}))
	if err != nil {
		return nil, fmt.Errorf("can't get kinks for game find post %s from database: %w", postId, err)
	}

	return kinkCodes, nil
}

// endregion

// endregion

// region Chats

// region Dialog

func GetDialogFromDatabase(database sq.StatementBuilderType, dialogId string) (*model.Dialog, error) {
	dialog := model.Dialog{}

	err := getEntityFromDatabase(DialogFieldsSelect(database), dialogId, []any{DialogFields(&dialog)})
	if err != nil {
		return nil, fmt.Errorf("can't scan dialog %s from database: %w", dialogId, err)
	}

	return &dialog, nil
}

func GetDialogsFromDatabase(database sq.StatementBuilderType, ownerId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.Dialog, error) {
	var dialogs []*model.Dialog

	dialogIds, err := getEntityIdsWithPaginationFromDatabase[model.Dialog](DialogFieldsSelect(database).Where(sq.Eq{"user_id": ownerId}), paginationSettings)
	if err != nil {
		return nil, fmt.Errorf("can't scan user %s dialogIds from database: %w", ownerId, err)
	}

	for _, dialogId := range *dialogIds {
		dialog, err := GetDialogFromDatabase(database, dialogId)
		if err != nil {
			return nil, fmt.Errorf("can't get dialog %s from database: %w", dialogId, err)
		}

		dialogs = append(dialogs, dialog)
	}

	return &dialogs, nil
}

// endregion

// region Message

func GetMessageFromDatabase(database sq.StatementBuilderType, messageId string, linkLevel int) (*model.Message, error) {
	message := model.Message{}
	linkedGameFindPostId := new(string)

	err := getEntityFromDatabase(MessageFieldsSelect(database), messageId, []any{MessageFields(&message, linkedGameFindPostId)})
	if err != nil {
		return nil, fmt.Errorf("can't scan message %s from database: %w", messageId, err)
	}

	gameFindPost, err := GetGameFindPostFromDatabase(database, *linkedGameFindPostId)
	if err != nil {
		return nil, fmt.Errorf("can't scan message %s from database: %w", messageId, err)
	}

	images, err := GetMessageImagesFromDatabase(database, messageId)
	if err != nil {
		return nil, fmt.Errorf("can't get images of message %s: %w", messageId, err)
	}

	linkedMessages := &[]*model.Message{}

	if linkLevel > 0 {
		linkedMessages, err = GetLinkedMessagesFromDatabase(database, messageId, linkLevel-1)
		if err != nil {
			return nil, fmt.Errorf("can't get linked messages of message %s: %w", messageId, err)
		}
	}

	message.LinkedGameFindPost = gameFindPost
	message.Images = *images
	message.LinkedMessages = *linkedMessages

	return &message, nil
}

func GetMessageImagesFromDatabase(database sq.StatementBuilderType, messageId string) (*[]*model.Image, error) {
	imageIds, err := getEntityIdsFromDatabase[string](MessageImageFieldsSelect(database).Where(sq.Eq{"message_id": messageId}))
	if err != nil {
		return nil, fmt.Errorf("can't get images for message %s from database: %w", messageId, err)
	}

	images, err := GetImagesFromFiles(imageIds)
	if err != nil {
		return nil, fmt.Errorf("can't get images for message %s from files: %w", messageId, err)
	}

	return images, nil
}

func GetLinkedMessagesFromDatabase(database sq.StatementBuilderType, messageId string, linkLevel int) (*[]*model.Message, error) {
	var linkedMessages []*model.Message

	linkedMessageIds, err := getEntityIdsFromDatabase[string](MessageLinkedMessageFieldsSelect(database).Where(sq.Eq{"new_message_id": messageId}))
	if err != nil {
		return nil, fmt.Errorf("can't get linked message ids for message %s from database: %w", messageId, err)
	}

	for _, linkedMessageId := range *linkedMessageIds {
		message, err := GetMessageFromDatabase(database, linkedMessageId, linkLevel)
		if err != nil {
			return nil, fmt.Errorf("can't get linked message %s: %w", linkedMessageId, err)
		}

		linkedMessages = append(linkedMessages, message)
	}

	return &linkedMessages, nil
}

func GetDialogMessagesFromDatabase(database sq.StatementBuilderType, dialogId string, paginationSettings model.ListByTimeSortPaginationSettings) (*[]*model.Message, error) {
	var messages []*model.Message

	rows, err := AddPaginationByCreatedAtToQuery(MessageFieldsSelect(database).Where(sq.Eq{"dialog_id": dialogId, "deleted_at": nil}), paginationSettings).Query()
	if err != nil {
		return nil, fmt.Errorf("can't get message ids: %w", err)
	}

	for rows.Next() {
		var message model.Message
		linkedGameFindPostId := new(string)

		err = rows.Scan(MessageFields(&message, linkedGameFindPostId))
		if err != nil {
			return nil, fmt.Errorf("can't scan message: %w", err)
		}

		linkedGameFindPost, err := GetGameFindPostFromDatabase(database, *linkedGameFindPostId)
		if err != nil {
			return nil, fmt.Errorf("can't get linked game find post: %w", err)
		}

		message.LinkedGameFindPost = linkedGameFindPost

		messages = append(messages, &message)
	}

	return &messages, nil
}

// endregion

// endregion
