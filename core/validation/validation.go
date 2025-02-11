package validation

import (
	"TestGoLandProject/core/database"
	"TestGoLandProject/core/utils"
	graph "TestGoLandProject/graph/generated"
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
	root.MaxArrayLength = maxArrayLengthDirective
	root.MinArrayLength = minArrayLengthDirective
	root.MaxValue = maxValueDirective
	root.MinValue = minValueDirective

	root.NotEmptyString = notEmptyStringDirective
	root.Authenticated = authenticatedDirective
	root.Email = emailDirective

	root.CatalogItemCode = catalogItemCodeDirective

	return nil
}

func maxLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	fieldLong := len(fieldValue)

	if fieldLong > *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field %d is characters long, expected maximum %d: %w", fieldName, fieldLong, *value, err))
	}

	return next(ctx)
}

func minLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	fieldLong := len(fieldValue)

	if fieldLong < *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field %d is characters long, expected minimum %d: %w", fieldName, fieldLong, *value, err))
	}

	return next(ctx)
}

func maxArrayLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].([]any)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	fieldLong := len(fieldValue)

	if fieldLong > *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field has %d items, expected maximum %d: %w", fieldName, fieldLong, *value, err))
	}

	return next(ctx)
}

func minArrayLengthDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].([]any)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	fieldLong := len(fieldValue)

	if fieldLong < *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field has %d items, expected minimum %d: %w", fieldName, fieldLong, *value, err))
	}

	return next(ctx)
}

func maxValueDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(int)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	if fieldValue > *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field has value = %d, but expected maximum %d: %w", fieldName, fieldValue, *value, err))
	}

	return next(ctx)
}

func minValueDirective(ctx context.Context, obj any, next graphql.Resolver, value *int) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(int)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	if fieldValue < *value {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field has value = %d, but expected minimum %d: %w", fieldName, fieldValue, *value, err))
	}

	return next(ctx)
}

func notEmptyStringDirective(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	trimmedValue := strings.TrimSpace(fieldValue)

	if len(trimmedValue) == 0 {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field can't contains only whitespace or cannot be empty: %w", fieldName, err))
	}

	return next(ctx)
}

func authenticatedDirective(ctx context.Context, _ any, next graphql.Resolver) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	err = resolverUtils.CheckUserAuthFromContext(ctx, *databaseInstance)
	if err != nil {
		return nil, err
	}

	return next(ctx)
}

func emailDirective(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	_, err = mail.ParseAddress(fieldValue)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s field has invalid email pattern: %s: %w", fieldName, fieldValue, err))
	}

	return next(ctx)
}

func catalogItemCodeDirective(ctx context.Context, obj any, next graphql.Resolver, catalog *model.Catalogs) (res any, err error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusInternalServerError, fmt.Errorf("can't get gin context: %w", err))
	}

	fieldName := *graphql.GetPathContext(ctx).Field
	objMap := obj.(map[string]interface{})

	fieldValue, ok := objMap[fieldName].(string)
	if !ok {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("field %s is not found in object or has other data type: %w", fieldName, err))
	}

	catalogItem := model.CatalogItem{}

	catalogItemRow := databaseInstance.QueryRow("select * from Catalog_item ci join Catalog c where ci.code = $1 and c.name = $2", fieldValue, *catalog)
	err = catalogItemRow.Scan(&catalogItem.ID, &catalogItem.Code, &catalogItem.Value, &catalogItem.CatalogID, &catalogItem.IsActive)
	if err != nil {
		return nil, utils.ResponseError(ginContext, http.StatusBadRequest, fmt.Errorf("%s catalog doesn't contains item with code %s: %w", *catalog, fieldValue, err))
	}

	return next(ctx)
}
