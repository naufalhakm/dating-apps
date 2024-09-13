DROP INDEX IF EXISTS idx_users_age;
DROP INDEX IF EXISTS idx_users_gender;
DROP INDEX IF EXISTS idx_users_location;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS postgis;