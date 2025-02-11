package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"TestGoLandProject/graph/model"
	"context"
	"fmt"
)

// CreateGallery is the resolver for the createGallery field.
func (r *mutationResolver) CreateGallery(ctx context.Context, newGallery model.NewGallery) (*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: CreateGallery - createGallery"))
}

// AddImagesFromGallery is the resolver for the addImagesFromGallery field.
func (r *mutationResolver) AddImagesFromGallery(ctx context.Context, imageIds model.ImageIds) (*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: AddImagesFromGallery - addImagesFromGallery"))
}

// DeleteImagesFromGallery is the resolver for the deleteImagesFromGallery field.
func (r *mutationResolver) DeleteImagesFromGallery(ctx context.Context, imageIds model.ImageIds) (*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: DeleteImagesFromGallery - deleteImagesFromGallery"))
}

// UpdateGallery is the resolver for the updateGallery field.
func (r *mutationResolver) UpdateGallery(ctx context.Context, galleryIDHolder model.IDHolder, newGallery model.NewGallery) (*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: UpdateGallery - updateGallery"))
}

// DeleteGallery is the resolver for the deleteGallery field.
func (r *mutationResolver) DeleteGallery(ctx context.Context, galleryIDHolder model.IDHolder) (*model.DeleteResult, error) {
	panic(fmt.Errorf("not implemented: DeleteGallery - deleteGallery"))
}

// GetGallery is the resolver for the getGallery field.
func (r *queryResolver) GetGallery(ctx context.Context, galleryIDHolder model.IDHolder) (*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: GetGallery - getGallery"))
}

// GetMyGalleries is the resolver for the getMyGalleries field.
func (r *queryResolver) GetMyGalleries(ctx context.Context, paginationSettings model.ListPaginationSettings) ([]*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: GetMyGalleries - getMyGalleries"))
}

// GetUserGalleries is the resolver for the getUserGalleries field.
func (r *queryResolver) GetUserGalleries(ctx context.Context, paginationSettings model.ListPaginationSettings) ([]*model.Gallery, error) {
	panic(fmt.Errorf("not implemented: GetUserGalleries - getUserGalleries"))
}
