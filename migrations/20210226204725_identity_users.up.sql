CREATE TABLE IF NOT EXISTS 
	identity_users (
    id VARCHAR(64) PRIMARY KEY,
	first_name  VARCHAR(100),
	last_name VARCHAR(100),
	email VARCHAR(120),
	password VARCHAR(120),
	company VARCHAR(64),
	post_code VARCHAR(64),
    terms INTEGER)
