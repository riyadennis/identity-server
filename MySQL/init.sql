Use identity_db;

Create table identity_users (
    id varchar(40) PRIMARY KEY,
    first_name varchar(20),
    last_name varchar(20),
    email varchar(20),
    password varchar(100),
    company varchar(20),
    post_code varchar(20),
    terms int);


INSERT INTO identity_users
(id, first_name, last_name,password,email, company, post_code, terms)
 VALUES ('testId', 'John', 'Doe', 'password','john@doe.com', 'test','code', 1)
