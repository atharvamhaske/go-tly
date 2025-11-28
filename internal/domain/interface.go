package domain

import (
	"context"

	"github.com/atharvamhaske/go-tly/internal/domain/models"
)

// repository interfaces
type URLRepo interface {
	Save(ctx context.Context, url *models.URL) error
	FindByKey(ctx context.Context, key string) (*models.URL, error)
	UpdateByKey(ctx context.Context, key string, updates map[string]any) error
	DeleteByKey(ctx context.Context, key string) error
}

type UserRepo interface {
	Create (ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

type AnalyticsRepo interface {
	InsertClick(ctx context.Context, h *models.ClickEvents) error
	GetAggregated(ctx context.Context, url string) (*models.AnalyticsSummary, error)
}

type DomainHealthRepo interface {
	Save(ctx context.Context, h *models.DomainHealth) error
	Get(ctx context.Context, url string) (*models.DomainHealth, error)
}

//cache interface
type URLCache interface {
	SetKey(ctx context.Context, key string, url string) error
	GetKey(ctx context.Context, key string) (string, error)
	DeleteKey(ctx context.Context, key string) (string, error)	
}

//service interfaces

type URLService interface {
	Shorten(ctx context.Context, req models.ShortenRequest) (*models.ShortenResponse, error)
	Edit(ctx context.Context, key string, req models.EditURLRequest) error
	Delete(ctx context.Context, key string) error
	Resolve(ctx context.Context, key string) (string, error)
}

type UserService interface {}

type AnalyticsService interface {}