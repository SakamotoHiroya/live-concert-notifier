package handler_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/google/uuid"
)

func TestAPIHandler_Follows(t *testing.T) {
	h := newTestAPIHandler(t)
	ctx := context.Background()

	userRes, err := h.UsersPost(ctx, &oas.CreateUserRequest{Email: uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("UsersPost: %v", err)
	}
	user := userRes.(*oas.User)

	siteURL, _ := url.Parse("https://example.com/" + uuid.NewString())
	artistRes, err := h.ArtistsPost(ctx, &oas.CreateArtistRequest{Name: "Follow Target", OfficialSiteURL: *siteURL})
	if err != nil {
		t.Fatalf("ArtistsPost: %v", err)
	}
	artist := artistRes.(*oas.Artist)

	followRes, err := h.UsersUserIdFollowsPost(ctx, &oas.FollowArtistRequest{ArtistID: artist.ID}, oas.UsersUserIdFollowsPostParams{UserId: user.ID})
	if err != nil {
		t.Fatalf("UsersUserIdFollowsPost: %v", err)
	}
	if _, ok := followRes.(*oas.UsersUserIdFollowsPostCreated); !ok {
		t.Fatalf("UsersUserIdFollowsPost = %T, want *oas.UsersUserIdFollowsPostCreated", followRes)
	}

	dupRes, err := h.UsersUserIdFollowsPost(ctx, &oas.FollowArtistRequest{ArtistID: artist.ID}, oas.UsersUserIdFollowsPostParams{UserId: user.ID})
	if err != nil {
		t.Fatalf("UsersUserIdFollowsPost (dup) returned error instead of Conflict: %v", err)
	}
	if _, ok := dupRes.(*oas.UsersUserIdFollowsPostConflict); !ok {
		t.Fatalf("UsersUserIdFollowsPost (dup) = %T, want *oas.UsersUserIdFollowsPostConflict", dupRes)
	}

	missingArtistRes, err := h.UsersUserIdFollowsPost(ctx, &oas.FollowArtistRequest{ArtistID: uuid.New()}, oas.UsersUserIdFollowsPostParams{UserId: user.ID})
	if err != nil {
		t.Fatalf("UsersUserIdFollowsPost (missing artist) returned error instead of NotFound: %v", err)
	}
	if _, ok := missingArtistRes.(*oas.UsersUserIdFollowsPostNotFound); !ok {
		t.Fatalf("UsersUserIdFollowsPost (missing artist) = %T, want *oas.UsersUserIdFollowsPostNotFound", missingArtistRes)
	}

	listRes, err := h.UsersUserIdFollowsGet(ctx, oas.UsersUserIdFollowsGetParams{UserId: user.ID})
	if err != nil {
		t.Fatalf("UsersUserIdFollowsGet: %v", err)
	}
	list, ok := listRes.(*oas.ArtistList)
	if !ok || list.Total != 1 {
		t.Fatalf("UsersUserIdFollowsGet = %+v (%T), want 1 followed artist", listRes, listRes)
	}

	delRes, err := h.UsersUserIdFollowsArtistIdDelete(ctx, oas.UsersUserIdFollowsArtistIdDeleteParams{UserId: user.ID, ArtistId: artist.ID})
	if err != nil {
		t.Fatalf("UsersUserIdFollowsArtistIdDelete: %v", err)
	}
	if _, ok := delRes.(*oas.UsersUserIdFollowsArtistIdDeleteNoContent); !ok {
		t.Fatalf("UsersUserIdFollowsArtistIdDelete = %T, want NoContent", delRes)
	}

	redoDelRes, err := h.UsersUserIdFollowsArtistIdDelete(ctx, oas.UsersUserIdFollowsArtistIdDeleteParams{UserId: user.ID, ArtistId: artist.ID})
	if err != nil {
		t.Fatalf("UsersUserIdFollowsArtistIdDelete (again) returned error instead of NotFound: %v", err)
	}
	if _, ok := redoDelRes.(*oas.ErrorResponse); !ok {
		t.Fatalf("UsersUserIdFollowsArtistIdDelete (again) = %T, want *oas.ErrorResponse", redoDelRes)
	}
}
