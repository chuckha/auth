CREATE TABLE sessions (
	id text,
	user_id text,
	expires text
);
CREATE TABLE tokens (
	token text,
	user_id text,
	expires text
);
CREATE TABLE users (
	id text
)
