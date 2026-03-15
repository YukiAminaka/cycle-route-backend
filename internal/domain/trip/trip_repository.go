package trip

import "context"

type ITripRepository interface {
	GetTripsByUserID(ctx context.Context, userID string) ([]*Trip, error)
	CountTripsByUserID(ctx context.Context, userID string) (int64, error)
	GetTripByID(ctx context.Context, id string) (*Trip, error)
	GetTripByKratosID(ctx context.Context, kratosID string) ([]*Trip, error)
	SaveTrip(ctx context.Context, trip *Trip) error
	DeleteTrip(ctx context.Context, id string) error
	UpdateTrip(ctx context.Context, trip *Trip) error
}