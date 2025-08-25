CREATE TABLE IF NOT EXISTS raw_ocr_data (
    publish_id UUID REFERENCES publishes(id) ON DELETE CASCADE,
    raw JSONB NOT NULL
);