CREATE TABLE IF NOT EXISTS comments
(
    id SERIAL PRIMARY KEY,
    parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    content_tsv TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', content)) STORED
);

CREATE INDEX IF NOT EXISTS idx_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_created_at ON comments(created_at);
CREATE INDEX IF NOT EXISTS idx_content_tsv ON comments USING GIN(content_tsv);