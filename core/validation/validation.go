package validation

import (
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/utils"
	"TestGoLandProject/graph"
	"TestGoLandProject/graph/model"
	resolverUtils "TestGoLandProject/graph/resolvers/utils"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"golang.org/x/net/context"
	"net/http"
	"net/mail"
	"strings"
)

var databaseInstance *database.Database = nil

func ImplementDirectives(root *graph.DirectiveRoot, database database.Database) error {
	databaseInstance = &database
	root.MaxLength = maxLengthDirective
	root.MinLength = minLengthDirective
	root.MaxLength = maxLengthDirective
	root.MinLength = minLengthDirective
	root.MaxValue = maxValueDirective
	root.MinValue = minValueDirective
	root.NotEmptyString = notEmptyStringDirective
	root.Authenticated = authenticatedDirective

	return nil
}

func maxLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	fieldLong := len(fieldValue)

	if fieldLong > *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field %d is characters long, expected maximum %d", fieldName, fieldLong, *value))
	}

	return next(ctx)
}

func minLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	fieldLong := len(fieldValue)

	if fieldLong < *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field %d is characters long, expected minimum %d", fieldName, fieldLong, *value))
	}

	return next(ctx)
}

func maxArrayLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].([]any)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	fieldLong := len(fieldValue)

	if fieldLong > *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field has %d items, expected maximum %d", fieldName, fieldLong, *value))
	}

	return next(ctx)
}

func minArrayLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].([]any)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	fieldLong := len(fieldValue)

	if fieldLong < *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field has %d items, expected minimum %d", fieldName, fieldLong, *value))
	}

	return next(ctx)
}

func maxValueDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(int)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	if fieldValue > *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field has value = %d, but expected maximum %d", fieldName, fieldValue, *value))
	}

	return next(ctx)
}

func minValueDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(int)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	if fieldValue < *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field has value = %d, but expected minimum %d", fieldName, fieldValue, *value))
	}

	return next(ctx)
}

func notEmptyStringDirective(ctx context.Context, obj any, next graphql.Resolver, state *bool) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	trimmedValue := strings.TrimSpace(fieldValue)

	if len(trimmedValue) == 0 && *state {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field can't contains only whitespace or cannot be empty", fieldName))
	}

	return next(ctx)
}

func authenticatedDirective(ctx context.Context, obj any, next graphql.Resolver, state *bool) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	err = resolverUtils.CheckUserAuthFromContext(ctx, *databaseInstance)
	if err != nil {
		return nil, err
	}

	return next(ctx)
}

func emailDirective(ctx context.Context, obj any, next graphql.Resolver, state *bool) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	_, err = mail.ParseAddress(fieldValue)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s field has invalid email pattern: %s", fieldName, fieldValue))
	}

	return next(ctx)
}

func catalogItemCodeDirective(ctx context.Context, obj any, next graphql.Resolver, catalog *string) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, "Can't get gin context")
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("field %s is not found in object or has other data type", fieldName))
	}

	catalogItem := model.CatalogItem{}

	catalogItemRow := databaseInstance.QueryRow("select * from Catalog_item ci join Catalog c where ci.code = $1 and c.name = $2", fieldValue, *catalog)
	err = catalogItemRow.Scan(&catalogItem.ID, &catalogItem.Code, &catalogItem.Value, &catalogItem.CatalogID, &catalogItem.IsActive)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Sprintf("%s catalog doesn't contains item with code %s", *catalog, fieldValue))
	}

	return next(ctx)
}
