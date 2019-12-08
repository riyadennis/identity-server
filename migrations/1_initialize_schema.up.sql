-- noinspection SqlDialectInspectionForFile
-- noinspection SqlNoDataSourceInspectionForFile
CREATE TABLE IF NOT EXISTS identity_users(id varchar(100) NOT NULL PRIMARY KEY,first_name  varchar(100),last_name varchar(100),email varchar(100),company varchar(100),post_code varchar(100),terms int, created_datetime DATETIME);