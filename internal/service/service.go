package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"tinylynx/internal/model"
	"tinylynx/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jxskiss/base62"
)

// Private functions
func fetchLink(ctx context.Context, query string, args ...any) (*model.Link, error) {
	var l model.Link
	
	err := storage.GetPool().QueryRow(ctx, query, args...).Scan(
		&l.ID, 
		&l.OriginalLink, 
		&l.ShortCode,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch scan: %w", err)
	}
	
	return &l, nil
}

func checkExistence(ctx context.Context, query string, args ...any) (bool, error) {
	var exists bool

	err := storage.GetPool().QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists scan: %w", err)
	}
	
	return exists, nil
}

func generateShortCode(link string) string {
	h := sha256.Sum256([]byte(link))
	return base62.EncodeToString(h[:6])
}

// Public functions

// Working with original_link
func GetByOriginalLink(ctx context.Context, link string) (*model.Link, error) {
	queryExistance := `SELECT EXISTS(SELECT 1 FROM links
	WHERE original_link = $1)`

	ok, err := checkExistence(ctx, queryExistance, link)
	if err != nil {
		return nil, fmt.Errorf("couldn't check existance: %w", err)
	}

	// Create entry if link is not found
	if !ok {
		queryInsert := `INSERT INTO links (original_link, short_code) VALUES ($1, $2)`
		_, err := storage.GetPool().Exec(ctx, queryInsert, link, generateShortCode(link))
		if err != nil {
			return nil, fmt.Errorf("insesrt values: %w", err)
		}
	}

	queryFind := `SELECT id, original_link, short_code FROM links WHERE original_link = $1`
	

	l, err := fetchLink(ctx, queryFind, link)
	if err != nil {
		return nil, fmt.Errorf("couldn't find link: %w", err)
	}

	return l, nil
}

func ExistsOriginalLink(ctx context.Context, link string) (bool, error) {
	query := `
	SELECT EXISTS(SELECT 1 FROM links WHERE original_link = $1 LIMIT 1)`
	return checkExistence(ctx, query, link)
}

// Working with short_code
func FindByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	query := `
	SELECT id, original_link, short_code FROM links
	WHERE short_code = $1`
	return fetchLink(ctx, query, shortCode)
}

func ExistsShortCode(ctx context.Context, shortCode string) (bool, error) {
	query := `
	SELECT EXISTS(SELECT 1 FROM links WHERE short_code = $1 LIMIT 1)`
	return checkExistence(ctx, query, shortCode)
}

// Analytics functions
func RecordAnalytics(ctx context.Context, linkID int, ip, userAgent, referrer, country, device, browser, platform string) error {
	query := `
	INSERT INTO link_analytics (link_id, ip_address, user_agent, referrer, country, device, browser, platform)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := storage.GetPool().Exec(ctx, query, linkID, ip, userAgent, referrer, country, device, browser, platform)
	if err != nil {
		return fmt.Errorf("record analytics: %w", err)
	}

	return nil
}

func GetLinkAnalytics(ctx context.Context, linkID int, limit int) ([]model.Analytics, error) {
	query := `
	SELECT id, link_id, created_at, ip_address, user_agent, referrer, country, device, browser, platform
	FROM link_analytics
	WHERE link_id = $1
	ORDER BY created_at DESC
	LIMIT $2`

	rows, err := storage.GetPool().Query(ctx, query, linkID, limit)
	if err != nil {
		return nil, fmt.Errorf("get analytics: %w", err)
	}
	defer rows.Close()

	var analytics []model.Analytics
	for rows.Next() {
		var a model.Analytics
		err := rows.Scan(
			&a.ID,
			&a.LinkID,
			&a.CreatedAt,
			&a.IPAddress,
			&a.UserAgent,
			&a.Referrer,
			&a.Country,
			&a.Device,
			&a.Browser,
			&a.Platform,
		)
		if err != nil {
			return nil, fmt.Errorf("scan analytics: %w", err)
		}
		analytics = append(analytics, a)
	}

	return analytics, nil
}

func GetLinkStats(ctx context.Context, linkID int) (map[string]interface{}, error) {
	query := `
	SELECT
		COUNT(*) as total_clicks,
		COUNT(DISTINCT ip_address) as unique_visitors,
		COUNT(DISTINCT country) as countries,
		COUNT(DISTINCT device) as devices,
		COUNT(DISTINCT browser) as browsers
	FROM link_analytics
	WHERE link_id = $1`

	var stats struct {
		TotalClicks  int `json:"total_clicks"`
		UniqueVisitors int `json:"unique_visitors"`
		Countries int `json:"countries"`
		Devices int `json:"devices"`
		Browsers int `json:"browsers"`
	}

	err := storage.GetPool().QueryRow(ctx, query, linkID).Scan(
		&stats.TotalClicks,
		&stats.UniqueVisitors,
		&stats.Countries,
		&stats.Devices,
		&stats.Browsers,
	)

	if err != nil {
		return nil, fmt.Errorf("get link stats: %w", err)
	}

	return map[string]interface{}{
		"link_id": linkID,
		"total_clicks": stats.TotalClicks,
		"unique_visitors": stats.UniqueVisitors,
		"countries": stats.Countries,
		"devices": stats.Devices,
		"browsers": stats.Browsers,
	}, nil
}
