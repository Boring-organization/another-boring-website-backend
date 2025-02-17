package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"TestGoLandProject/graph/model"
	"context"
	"fmt"
)

// DeleteDialog is the resolver for the deleteDialog field.
func (r *mutationResolver) DeleteDialog(ctx context.Context, dialogIDHolder model.IDHolder) (*model.DeleteResult, error) {
	panic(fmt.Errorf("not implemented: DeleteDialog - deleteDialog"))
}

// ToggleDialogNotification is the resolver for the toggleDialogNotification field.
func (r *mutationResolver) ToggleDialogNotification(ctx context.Context, stateHolder model.ToggleNotificationData) (bool, error) {
	panic(fmt.Errorf("not implemented: ToggleDialogNotification - toggleDialogNotification"))
}

// GetMyDialogs is the resolver for the getMyDialogs field.
func (r *queryResolver) GetMyDialogs(ctx context.Context, paginationSettings model.ListPaginationSettings) ([]*model.Dialog, error) {
	panic(fmt.Errorf("not implemented: GetMyDialogs - getMyDialogs"))
}
