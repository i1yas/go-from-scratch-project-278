-- +goose Up
CREATE TABLE links (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    original_url VARCHAR(200) NOT NULL,
    shortcode VARCHAR(32) NOT NULL,

    CONSTRAINT unique_shortcode UNIQUE (shortcode)
);

-- +goose Down
DROP TABLE links;
