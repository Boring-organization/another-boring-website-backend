package databaseUtils

import (
	"TestGoLandProject/graph/model"
	sq "github.com/Masterminds/squirrel"
)

var UserAsFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("u.id", "u.nickname", "u.description", "u.sex", "u.is_admin", "u.created_at", "u.edited_at", "user.deleted_at").From("User")
}
var UserFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "nickname", "description", "sex", "is_admin", "created_at", "edited_at", "deleted_at").From("User")
}

func UserFields(user *model.User) (*string, *string, *string, *string, *bool, *int, *int, *int) {
	return &user.ID, &user.Nickname, user.Description, &user.Sex, &user.IsAdmin, &user.CreatedAt, user.EditedAt, user.DeletedAt
}

var PostFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "author_id", "created_at", "edited_at", "deleted_at").From("Post")
}

func PostFields(post *model.Post, postAuthorId *string) (*string, *string, *string, *int, *int, *int) {
	return &post.ID, &post.Text, postAuthorId, &post.CreatedAt, post.EditedAt, post.DeletedAt
}

var GameFindPostFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "author_id", "is_likable", "is_anonymously", "created_at", "edited_at", "deleted_at").From("Game_find_post")
}

func GameFindPostFields(post *model.GameFindPost, postAuthorId *string) (*string, *string, *string, *bool, *bool, *int, *int, *int) {
	return &post.ID, &post.Text, postAuthorId, &post.IsLikable, &post.IsAnonymously, &post.CreatedAt, post.EditedAt, post.DeletedAt
}

var GameFindPostKinkFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("catalog_item_code").From("Game_find_post_kink")
}

var GalleryFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "name", "author_id", "is_closed", "created_at", "edited_at", "deleted_at").From("Gallery")
}

func GalleryFields(gallery *model.Gallery, galleryAuthorId *string) (*string, *string, *string, *bool, *int, *int, *int) {
	return &gallery.ID, &gallery.Name, galleryAuthorId, &gallery.IsClosed, &gallery.CreatedAt, gallery.EditedAt, gallery.DeletedAt
}

var CatalogItemFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "code", "value", "catalog_id", "is_active").From("Catalog_item")
}

func CatalogItemFields(catalogItem *model.CatalogItem, catalogId *string) (*string, *string, *string, *string, *bool) {
	return &catalogItem.ID, &catalogItem.Code, &catalogItem.Value, catalogId, &catalogItem.IsActive
}

var AuthDataFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("device_id", "device_name", "used_at", "expired_at").From("Auth_data")
}

func AuthDataFields(authData *model.AuthData) (*string, *string, *int, *int) {
	return &authData.DeviceName, &authData.DeviceName, &authData.UsedAt, &authData.ExpiredAt
}

var DialogFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("is_notification_on").From("Dialog_user")
}

func DialogFields(dialog *model.Dialog) *bool {
	return &dialog.IsNotificationOn
}

var MessageFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "author_id", "linked_game_find_post_id", "created_at", "edited_at", "deleted_at").From("Message")
}

func MessageFields(message *model.Message, linkedGameFindPostId *string) (*string, *string, *string, *int, *int, *int) {
	return &message.ID, &message.AuthorID, linkedGameFindPostId, &message.CreatedAt, message.EditedAt, message.DeletedAt
}

var MessageImageFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("image_id").From("Message_image")
}

var MessageLinkedMessageFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("linked_message_id").From("Message_message")
}

var TemporaryImageFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("image_id").From("Temporary_image")
}

// region Comments

var PostCommentFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "author_id", "linked_user_id", "created_at", "edited_at", "deleted_at").From("Post_comment")
}

func PostCommentFields(postComment *model.PostComment, commentAuthorId *string, linkedUserId *string) (*string, *string, *string, *string, *int, *int, *int) {
	return &postComment.ID, &postComment.Text, commentAuthorId, linkedUserId, &postComment.CreatedAt, postComment.EditedAt, postComment.DeletedAt
}

var GameFindPostCommentFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "author_id", "linked_user_id", "created_at", "edited_at", "deleted_at").From("Game_find_post_comment")
}

func GameFindPostCommentFields(postComment *model.GameFindPostComment, commentAuthorId *string, linkedUserId *string) (*string, *string, *string, *string, *int, *int, *int) {
	return &postComment.ID, &postComment.Text, commentAuthorId, linkedUserId, &postComment.CreatedAt, postComment.EditedAt, postComment.DeletedAt
}

// endregion

// region Complaints

var UserComplaintFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "user_id", "author_id", "created_at", "deleted_at").From("User_complaint")
}

func UserComplaintFields(userComment *model.UserComplaint, complaintUserId *string, complaintAuthorId *string) (*string, *string, *string, *string, *int, *int) {
	return &userComment.ID, &userComment.Text, complaintUserId, complaintAuthorId, &userComment.CreatedAt, userComment.DeletedAt
}

var PostComplaintFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "post_id", "author_id", "created_at", "deleted_at").From("Post_complaint")
}

func PostComplaintFields(postComplaint *model.PostComplaint, postId *string, complaintAuthorId *string) (*string, *string, *string, *string, *int, *int) {
	return &postComplaint.ID, &postComplaint.Text, postId, complaintAuthorId, &postComplaint.CreatedAt, postComplaint.DeletedAt
}

var PostCommentComplaintFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "comment_id", "author_id", "created_at", "deleted_at").From("Post_comment_complaint")
}

func PostCommentComplaintFields(postCommentComplaint *model.PostCommentComplaint, postCommentId *string, complaintAuthorId *string) (*string, *string, *string, *string, *int, *int) {
	return &postCommentComplaint.ID, &postCommentComplaint.Text, postCommentId, complaintAuthorId, &postCommentComplaint.CreatedAt, postCommentComplaint.DeletedAt
}

var GameFindPostComplaintFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "post_id", "author_id", "created_at", "deleted_at").From("Game_find_post_complaint")
}

func GameFindPostComplaintFields(gameFindPostComplaint *model.GameFindPostComplaint, gameFindPostId *string, complaintAuthorId *string) (*string, *string, *string, *string, *int, *int) {
	return &gameFindPostComplaint.ID, &gameFindPostComplaint.Text, gameFindPostId, complaintAuthorId, &gameFindPostComplaint.CreatedAt, gameFindPostComplaint.DeletedAt
}

var GameFindPostCommentComplaintFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "comment_id", "author_id", "created_at", "deleted_at").From("Game_find_post_comment_complaint")
}

func GameFindPostCommentComplaintFields(gameFindPostCommentComplaint *model.GameFindPostCommentComplaint, gameFindPostCommentId *string, complaintAuthorId *string) (*string, *string, *string, *string, *int, *int) {
	return &gameFindPostCommentComplaint.ID, &gameFindPostCommentComplaint.Text, gameFindPostCommentId, complaintAuthorId, &gameFindPostCommentComplaint.CreatedAt, gameFindPostCommentComplaint.DeletedAt
}

var GalleryComplaintFieldsSelect = func(builderType sq.StatementBuilderType) sq.SelectBuilder {
	return sq.Select("id", "text", "gallery_id", "author_id", "created_at", "deleted_at").From("Gallery_complaint")
}

func GalleryComplaintFields(galleryComplaint *model.GalleryComplaint, galleryIg *string, complaintAuthorId *string) (*string, *string, *string, *string, *int, *int) {
	return &galleryComplaint.ID, &galleryComplaint.Text, galleryIg, complaintAuthorId, &galleryComplaint.CreatedAt, galleryComplaint.DeletedAt
}

// endregion
