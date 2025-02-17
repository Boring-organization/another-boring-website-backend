package commonUtils

import (
	"TestGoLandProject/core/utils/database_utils"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/net/context"
	"time"
)

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		return nil, fmt.Errorf("could not retrieve gin.Context")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("gin.Context has wrong type")
	}

	return gc, nil
}

func ResponseError(ginContext *gin.Context, status int, err error) *gqlerror.Error {
	ginContext.Status(status)
	return gqlerror.Errorf(err.Error())
}

func PeriodicImageCleaner(ctx context.Context, database sq.StatementBuilderType) {
	ticker := time.NewTicker(24 * time.Hour)

	for {
		select {
		case <-ticker.C:
			err := databaseUtils.ClearTemporaryImageTable(database)
			if err != nil {
				fmt.Printf("Unhandled error catched: %w", err)
				return
			}

		case <-ctx.Done():
			ticker.Stop()
		}
	}
}
