CREATE TABLE IF NOT EXISTS status_pages (
    id SERIAL PRIMARY KEY,
    hostname TEXT UNIQUE NOT NULL,
    page_data_url TEXT NOT NULL
);

