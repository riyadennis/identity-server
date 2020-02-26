Use identity_db;

Create table identity_users (
    id varchar(40) PRIMARY KEY,
    first_name varchar(20),
    last_name varchar(20),
    email varchar(20),
    password varchar(20),
    company varchar(20),
    post_code varchar(20),
    terms int);
