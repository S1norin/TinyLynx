-- +goose Up
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_link TEXT NOT NULL,
    short_code TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE link_analytics (
    id SERIAL PRIMARY KEY,
    link_id INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT,
    country TEXT,
    device TEXT,
    browser TEXT,
    platform TEXT,
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);

CREATE INDEX idx_links_short_code ON links(short_code);
CREATE INDEX idx_links_original_link ON links(original_link);
CREATE INDEX idx_analytics_link_id ON link_analytics(link_id);
CREATE INDEX idx_analytics_created_at ON link_analytics(created_at);

-- +goose Down
DROP TABLE IF EXISTS link_analytics;
DROP TABLE IF EXISTS links;