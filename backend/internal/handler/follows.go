package handler

import (
	"context"
	"errors"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
)

// UsersUserIdFollowsGet implements GET /users/{userId}/follows.
func (h *APIHandler) UsersUserIdFollowsGet(ctx context.Context, params oas.UsersUserIdFollowsGetParams) (oas.UsersUserIdFollowsGetRes, error) {
	artists, err := h.follows.ListFollowed(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "user not found"}, nil
		}
		return nil, err
	}
	items := make([]oas.Artist, 0, len(artists))
	for _, a := range artists {
		oa, err := toOASArtist(a)
		if err != nil {
			return nil, err
		}
		items = append(items, *oa)
	}
	return &oas.ArtistList{Items: items, Total: len(items)}, nil
}

// UsersUserIdFollowsPost implements POST /users/{userId}/follows.
func (h *APIHandler) UsersUserIdFollowsPost(ctx context.Context, req *oas.FollowArtistRequest, params oas.UsersUserIdFollowsPostParams) (oas.UsersUserIdFollowsPostRes, error) {
	err := h.follows.Follow(ctx, params.UserId, req.ArtistID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.UsersUserIdFollowsPostNotFound{Code: "NOT_FOUND", Message: "user or artist not found"}, nil
		}
		if errors.Is(err, repository.ErrConflict) {
			return &oas.UsersUserIdFollowsPostConflict{Code: "CONFLICT", Message: "already following"}, nil
		}
		return nil, err
	}
	return &oas.UsersUserIdFollowsPostCreated{}, nil
}

// UsersUserIdFollowsArtistIdDelete implements DELETE /users/{userId}/follows/{artistId}.
func (h *APIHandler) UsersUserIdFollowsArtistIdDelete(ctx context.Context, params oas.UsersUserIdFollowsArtistIdDeleteParams) (oas.UsersUserIdFollowsArtistIdDeleteRes, error) {
	if err := h.follows.Unfollow(ctx, params.UserId, params.ArtistId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "follow relationship not found"}, nil
		}
		return nil, err
	}
	return &oas.UsersUserIdFollowsArtistIdDeleteNoContent{}, nil
}
