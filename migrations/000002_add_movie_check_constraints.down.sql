ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_genres_length_check;

ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_runtime_check;

ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_year_check;
