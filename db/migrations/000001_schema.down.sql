-- 000001_schema.down.sql: Drop all tables created in the schema migration

-- Drop indexes first
DROP INDEX IF EXISTS idx_payment_customer_id;
DROP INDEX IF EXISTS idx_rental_inventory_id;
DROP INDEX IF EXISTS idx_rental_customer_id;
DROP INDEX IF EXISTS idx_inventory_film_id;
DROP INDEX IF EXISTS idx_film_title;
DROP INDEX IF EXISTS idx_actor_last_name;
DROP INDEX IF EXISTS idx_customer_last_name;

-- Drop tables in reverse order of creation to handle dependencies
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS rental;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS film_actors;
DROP TABLE IF EXISTS actor;
DROP TABLE IF EXISTS film;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS store;