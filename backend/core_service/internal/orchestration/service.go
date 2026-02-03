package orchestration

import (
	"context"

	"core-service/internal/preferences"
)

type Service struct {
	prefRepo *preferences.Repository
}

func NewService(prefRepo *preferences.Repository) *Service {
	return &Service{
		prefRepo: prefRepo,
	}
}

type UserContext struct {
	UserID      string
	Preferences []preferences.Preference
}

func (s *Service) BuildUserContext(
	ctx context.Context,
	userID string,
) (*UserContext, error) {

	prefs, err := s.prefRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserContext{
		UserID:      userID,
		Preferences: prefs,
	}, nil
}
