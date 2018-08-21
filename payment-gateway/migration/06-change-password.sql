CREATE TABLE password_change_requests (
	id          SERIAL PRIMARY KEY,
	user_id     INT REFERENCES account(id) NOT NULL,
	token 		TEXT NOT NULL UNIQUE,
	create_at   TIMESTAMP NOT NULL DEFAULT NOW()
);