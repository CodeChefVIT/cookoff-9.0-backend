CREATE TABLE users (
	id uuid PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL
);