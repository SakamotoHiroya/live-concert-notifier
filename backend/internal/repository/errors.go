package repository

import "errors"

// ErrNotFound is returned when a lookup finds no matching row.
var ErrNotFound = errors.New("repository: not found")

// ErrConflict is returned when a write would violate a uniqueness
// constraint (e.g. duplicate email, official_site_url, or concert).
var ErrConflict = errors.New("repository: conflict")
