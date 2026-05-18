ALTER TABLE identity_users ADD COLUMN created_by VARCHAR(64) DEFAULT '' AFTER terms;
