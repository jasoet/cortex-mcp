-- 000002_seed_data.down.sql: Remove all seed data in reverse order

-- Delete all payment data
DELETE FROM payment;

-- Delete all rental data
DELETE FROM rental;

-- Delete all inventory data
DELETE FROM inventory;

-- Delete all film-actor relationships
DELETE FROM film_actors;

-- Delete all actor data
DELETE FROM actor;

-- Delete all film data
DELETE FROM film;

-- Delete all category data
DELETE FROM category;

-- Delete all customer data
DELETE FROM customer;

-- Delete all staff data
DELETE FROM staff;

-- Delete all store data
DELETE FROM store;