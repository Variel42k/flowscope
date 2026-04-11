package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/flowscope/flowscope/internal/model"
)

func (r *Repository) ListSavedViews(ctx context.Context, username string, scope string) ([]model.SavedView, error) {
	username = strings.TrimSpace(username)
	scope = strings.ToLower(strings.TrimSpace(scope))
	rows, err := r.db.QueryContext(ctx, `
SELECT view_id, name, description, scope, owner_user, is_shared, filters_json, created_at, updated_at
FROM (
	SELECT view_id,
	       argMax(name, updated_at) AS name,
	       argMax(description, updated_at) AS description,
	       argMax(scope, updated_at) AS scope,
	       argMax(owner_user, updated_at) AS owner_user,
	       argMax(is_shared, updated_at) AS is_shared,
	       argMax(filters_json, updated_at) AS filters_json,
	       min(created_at) AS created_at,
	       max(updated_at) AS updated_at,
	       argMax(deleted, updated_at) AS deleted
	FROM saved_views
	GROUP BY view_id
)
WHERE deleted = 0
  AND (owner_user = ? OR is_shared = 1)
  AND (? = '' OR scope = ?)
ORDER BY updated_at DESC`, username, scope, scope)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.SavedView, 0, 32)
	for rows.Next() {
		var item model.SavedView
		var rawFilters string
		if err := rows.Scan(
			&item.ViewID,
			&item.Name,
			&item.Description,
			&item.Scope,
			&item.OwnerUser,
			&item.IsShared,
			&rawFilters,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Filters = parseFiltersJSON(rawFilters)
		out = append(out, item)
	}
	return out, nil
}

func (r *Repository) CreateSavedView(ctx context.Context, view model.SavedView) (model.SavedView, error) {
	now := time.Now().UTC()
	if strings.TrimSpace(view.ViewID) == "" {
		view.ViewID = uuid.NewString()
	}
	if view.CreatedAt.IsZero() {
		view.CreatedAt = now
	}
	view.UpdatedAt = now
	if view.Filters == nil {
		view.Filters = map[string]string{}
	}
	rawFilters, _ := json.Marshal(view.Filters)
	_, err := r.db.ExecContext(ctx, `
INSERT INTO saved_views (
	view_id, name, description, scope, owner_user, is_shared, filters_json, created_at, updated_at, deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		view.ViewID,
		view.Name,
		view.Description,
		view.Scope,
		view.OwnerUser,
		view.IsShared,
		string(rawFilters),
		view.CreatedAt,
		view.UpdatedAt,
		false,
	)
	return view, err
}

func (r *Repository) GetSavedViewByID(ctx context.Context, viewID string) (model.SavedView, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT view_id, name, description, scope, owner_user, is_shared, filters_json, created_at, updated_at, deleted
FROM (
	SELECT view_id,
	       argMax(name, updated_at) AS name,
	       argMax(description, updated_at) AS description,
	       argMax(scope, updated_at) AS scope,
	       argMax(owner_user, updated_at) AS owner_user,
	       argMax(is_shared, updated_at) AS is_shared,
	       argMax(filters_json, updated_at) AS filters_json,
	       min(created_at) AS created_at,
	       max(updated_at) AS updated_at,
	       argMax(deleted, updated_at) AS deleted
	FROM saved_views
	WHERE view_id = ?
	GROUP BY view_id
)
LIMIT 1`, viewID)
	var item model.SavedView
	var rawFilters string
	var deleted bool
	if err := row.Scan(
		&item.ViewID,
		&item.Name,
		&item.Description,
		&item.Scope,
		&item.OwnerUser,
		&item.IsShared,
		&rawFilters,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deleted,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.SavedView{}, errors.New("saved view not found")
		}
		return model.SavedView{}, err
	}
	if deleted {
		return model.SavedView{}, errors.New("saved view not found")
	}
	item.Filters = parseFiltersJSON(rawFilters)
	return item, nil
}

func (r *Repository) UpdateSavedView(ctx context.Context, view model.SavedView) (model.SavedView, error) {
	if strings.TrimSpace(view.ViewID) == "" {
		return model.SavedView{}, errors.New("missing view id")
	}
	if view.Filters == nil {
		view.Filters = map[string]string{}
	}
	if view.CreatedAt.IsZero() {
		existing, err := r.GetSavedViewByID(ctx, view.ViewID)
		if err != nil {
			return model.SavedView{}, err
		}
		view.CreatedAt = existing.CreatedAt
	}
	view.UpdatedAt = time.Now().UTC()
	rawFilters, _ := json.Marshal(view.Filters)
	_, err := r.db.ExecContext(ctx, `
INSERT INTO saved_views (
	view_id, name, description, scope, owner_user, is_shared, filters_json, created_at, updated_at, deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		view.ViewID,
		view.Name,
		view.Description,
		view.Scope,
		view.OwnerUser,
		view.IsShared,
		string(rawFilters),
		view.CreatedAt,
		view.UpdatedAt,
		false,
	)
	return view, err
}

func (r *Repository) DeleteSavedView(ctx context.Context, viewID string) error {
	view, err := r.GetSavedViewByID(ctx, viewID)
	if err != nil {
		return err
	}
	rawFilters, _ := json.Marshal(view.Filters)
	_, err = r.db.ExecContext(ctx, `
INSERT INTO saved_views (
	view_id, name, description, scope, owner_user, is_shared, filters_json, created_at, updated_at, deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		view.ViewID,
		view.Name,
		view.Description,
		view.Scope,
		view.OwnerUser,
		view.IsShared,
		string(rawFilters),
		view.CreatedAt,
		time.Now().UTC(),
		true,
	)
	return err
}

func parseFiltersJSON(raw string) map[string]string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}
	}
	out := make(map[string]string)
	if err := json.Unmarshal([]byte(raw), &out); err == nil {
		return out
	}
	return map[string]string{}
}
