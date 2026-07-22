package handler

import (
	"context"
	"errors"
	"net/url"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
)

const defaultLimit = 20

func toOASArtist(a domain.Artist) (*oas.Artist, error) {
	u, err := url.Parse(a.OfficialSiteURL)
	if err != nil {
		return nil, err
	}
	return &oas.Artist{ID: a.ID, Name: a.Name, OfficialSiteURL: *u, CreatedAt: a.CreatedAt}, nil
}

// ArtistsGet implements GET /artists.
func (h *APIHandler) ArtistsGet(ctx context.Context, params oas.ArtistsGetParams) (*oas.ArtistList, error) {
	q := ""
	if params.Q.IsSet() {
		q = params.Q.Value
	}
	limit, offset := int32(defaultLimit), int32(0)
	if params.Limit.IsSet() {
		limit = int32(params.Limit.Value)
	}
	if params.Offset.IsSet() {
		offset = int32(params.Offset.Value)
	}

	artists, total, err := h.artists.ListArtists(ctx, q, limit, offset)
	if err != nil {
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
	return &oas.ArtistList{Items: items, Total: int(total)}, nil
}

// ArtistsPost implements POST /artists.
func (h *APIHandler) ArtistsPost(ctx context.Context, req *oas.CreateArtistRequest) (oas.ArtistsPostRes, error) {
	artist, err := h.artists.CreateArtist(ctx, req.Name, req.OfficialSiteURL.String())
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return &oas.ArtistsPostConflict{Code: "CONFLICT", Message: "official_site_url already registered"}, nil
		}
		return nil, err
	}
	return toOASArtist(artist)
}

// ArtistsArtistIdGet implements GET /artists/{artistId}.
func (h *APIHandler) ArtistsArtistIdGet(ctx context.Context, params oas.ArtistsArtistIdGetParams) (oas.ArtistsArtistIdGetRes, error) {
	artist, err := h.artists.GetArtist(ctx, params.ArtistId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "artist not found"}, nil
		}
		return nil, err
	}
	return toOASArtist(artist)
}

// ArtistsArtistIdPut implements PUT /artists/{artistId}.
func (h *APIHandler) ArtistsArtistIdPut(ctx context.Context, req *oas.UpdateArtistRequest, params oas.ArtistsArtistIdPutParams) (oas.ArtistsArtistIdPutRes, error) {
	var name, officialSiteURL *string
	if req.Name.IsSet() {
		v := req.Name.Value
		name = &v
	}
	if req.OfficialSiteURL.IsSet() {
		v := req.OfficialSiteURL.Value.String()
		officialSiteURL = &v
	}

	artist, err := h.artists.UpdateArtist(ctx, params.ArtistId, name, officialSiteURL)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ArtistsArtistIdPutNotFound{Code: "NOT_FOUND", Message: "artist not found"}, nil
		}
		return nil, err
	}
	return toOASArtist(artist)
}

// ArtistsArtistIdDelete implements DELETE /artists/{artistId}.
func (h *APIHandler) ArtistsArtistIdDelete(ctx context.Context, params oas.ArtistsArtistIdDeleteParams) (oas.ArtistsArtistIdDeleteRes, error) {
	if err := h.artists.DeleteArtist(ctx, params.ArtistId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "artist not found"}, nil
		}
		return nil, err
	}
	return &oas.ArtistsArtistIdDeleteNoContent{}, nil
}
