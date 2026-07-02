package handler_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/google/uuid"
)

func TestAPIHandler_Artists(t *testing.T) {
	h := newTestAPIHandler(t)
	ctx := context.Background()

	siteURL, _ := url.Parse("https://example.com/" + uuid.NewString())
	createRes, err := h.ArtistsPost(ctx, &oas.CreateArtistRequest{Name: "Test Artist", OfficialSiteURL: *siteURL})
	if err != nil {
		t.Fatalf("ArtistsPost: %v", err)
	}
	created, ok := createRes.(*oas.Artist)
	if !ok {
		t.Fatalf("ArtistsPost returned %T, want *oas.Artist", createRes)
	}

	dupRes, err := h.ArtistsPost(ctx, &oas.CreateArtistRequest{Name: "Dup", OfficialSiteURL: *siteURL})
	if err != nil {
		t.Fatalf("ArtistsPost (duplicate url) returned error instead of Conflict: %v", err)
	}
	if _, ok := dupRes.(*oas.ArtistsPostConflict); !ok {
		t.Fatalf("ArtistsPost (duplicate url) = %T, want *oas.ArtistsPostConflict", dupRes)
	}

	getRes, err := h.ArtistsArtistIdGet(ctx, oas.ArtistsArtistIdGetParams{ArtistId: created.ID})
	if err != nil {
		t.Fatalf("ArtistsArtistIdGet: %v", err)
	}
	if got, ok := getRes.(*oas.Artist); !ok || got.Name != created.Name {
		t.Fatalf("ArtistsArtistIdGet = %+v (%T), want name %q", getRes, getRes, created.Name)
	}

	newName := oas.OptString{Value: "Renamed Artist", Set: true}
	putRes, err := h.ArtistsArtistIdPut(ctx, &oas.UpdateArtistRequest{Name: newName}, oas.ArtistsArtistIdPutParams{ArtistId: created.ID})
	if err != nil {
		t.Fatalf("ArtistsArtistIdPut: %v", err)
	}
	updated, ok := putRes.(*oas.Artist)
	if !ok || updated.Name != "Renamed Artist" {
		t.Fatalf("ArtistsArtistIdPut = %+v (%T), want name Renamed Artist", putRes, putRes)
	}
	if updated.OfficialSiteURL.String() != created.OfficialSiteURL.String() {
		t.Fatalf("ArtistsArtistIdPut changed official_site_url without being asked to: %v", updated.OfficialSiteURL)
	}

	listRes, err := h.ArtistsGet(ctx, oas.ArtistsGetParams{Q: oas.OptString{Value: "Renamed", Set: true}})
	if err != nil {
		t.Fatalf("ArtistsGet: %v", err)
	}
	if listRes.Total < 1 {
		t.Fatalf("ArtistsGet total = %d, want >= 1", listRes.Total)
	}

	delRes, err := h.ArtistsArtistIdDelete(ctx, oas.ArtistsArtistIdDeleteParams{ArtistId: created.ID})
	if err != nil {
		t.Fatalf("ArtistsArtistIdDelete: %v", err)
	}
	if _, ok := delRes.(*oas.ArtistsArtistIdDeleteNoContent); !ok {
		t.Fatalf("ArtistsArtistIdDelete = %T, want *oas.ArtistsArtistIdDeleteNoContent", delRes)
	}

	notFoundRes, err := h.ArtistsArtistIdGet(ctx, oas.ArtistsArtistIdGetParams{ArtistId: uuid.New()})
	if err != nil {
		t.Fatalf("ArtistsArtistIdGet (missing) returned error instead of NotFound: %v", err)
	}
	if _, ok := notFoundRes.(*oas.ErrorResponse); !ok {
		t.Fatalf("ArtistsArtistIdGet (missing) = %T, want *oas.ErrorResponse", notFoundRes)
	}
}
