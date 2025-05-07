-- 000001_schema.up.sql: DVD Rental Database Schema (Simplified)

-- Stores table (with direct address fields)
CREATE TABLE store
(
    store_id    SERIAL PRIMARY KEY,
    store_name  VARCHAR(100) NOT NULL,
    address     VARCHAR(100) NOT NULL,
    address2    VARCHAR(100),
    district    VARCHAR(50)  NOT NULL,
    city        VARCHAR(50)  NOT NULL,
    country     VARCHAR(50)  NOT NULL,
    postal_code VARCHAR(20)  NOT NULL,
    phone       VARCHAR(20)  NOT NULL
);

-- Staff table (employees working at stores, with direct address fields)
CREATE TABLE staff
(
    staff_id    SERIAL PRIMARY KEY,
    store_id    INT          NOT NULL REFERENCES store (store_id),
    first_name  VARCHAR(50)  NOT NULL,
    last_name   VARCHAR(50)  NOT NULL,
    email       VARCHAR(100) NOT NULL,
    username    VARCHAR(50)  NOT NULL,
    address     VARCHAR(100) NOT NULL,
    address2    VARCHAR(100),
    district    VARCHAR(50)  NOT NULL,
    city        VARCHAR(50)  NOT NULL,
    country     VARCHAR(50)  NOT NULL,
    postal_code VARCHAR(20)  NOT NULL,
    phone       VARCHAR(20)  NOT NULL,
    active      BOOLEAN      NOT NULL DEFAULT TRUE,
    last_update TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email),
    UNIQUE (username)
);

-- Customers table (store customers, with direct address fields)
CREATE TABLE customer
(
    customer_id SERIAL PRIMARY KEY,
    store_id    INT          NOT NULL REFERENCES store (store_id),
    first_name  VARCHAR(50)  NOT NULL,
    last_name   VARCHAR(50)  NOT NULL,
    email       VARCHAR(100) NOT NULL,
    address     VARCHAR(100) NOT NULL,
    address2    VARCHAR(100),
    district    VARCHAR(50)  NOT NULL,
    city        VARCHAR(50)  NOT NULL,
    country     VARCHAR(50)  NOT NULL,
    postal_code VARCHAR(20)  NOT NULL,
    phone       VARCHAR(20)  NOT NULL,
    active      BOOLEAN      NOT NULL DEFAULT TRUE,
    create_date DATE         NOT NULL DEFAULT CURRENT_DATE,
    UNIQUE (email)
);

-- Categories table (film genres)
CREATE TABLE category
(
    category_id SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE
);

-- Films table (movies available for rental)
CREATE TABLE film
(
    film_id      SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    release_year SMALLINT     NOT NULL,
    length       SMALLINT     NOT NULL,
    category_id  INT          NOT NULL REFERENCES category (category_id)
);

-- Actors table
CREATE TABLE actor
(
    actor_id   SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name  VARCHAR(50) NOT NULL
);

-- Film actors join table (many-to-many relationship between films and actors)
CREATE TABLE film_actors
(
    film_id  INT NOT NULL REFERENCES film (film_id),
    actor_id INT NOT NULL REFERENCES actor (actor_id),
    PRIMARY KEY (film_id, actor_id)
);

-- Inventory table (copies of films in each store)
CREATE TABLE inventory
(
    inventory_id SERIAL PRIMARY KEY,
    film_id      INT NOT NULL REFERENCES film (film_id),
    store_id     INT NOT NULL REFERENCES store (store_id)
    -- multiple copies of the same film at a store are allowed
);

-- Rentals table (film rental transactions)
CREATE TABLE rental
(
    rental_id    SERIAL PRIMARY KEY,
    rental_date  TIMESTAMP NOT NULL,
    inventory_id INT       NOT NULL REFERENCES inventory (inventory_id),
    customer_id  INT       NOT NULL REFERENCES customer (customer_id),
    return_date  TIMESTAMP,
    staff_id     INT       NOT NULL REFERENCES staff (staff_id)
    -- You could add a constraint to ensure an inventory item is not double-booked at the same time
);

-- Payments table (payments for rentals)
CREATE TABLE payment
(
    payment_id   SERIAL PRIMARY KEY,
    customer_id  INT           NOT NULL REFERENCES customer (customer_id),
    staff_id     INT           NOT NULL REFERENCES staff (staff_id),
    rental_id    INT           NOT NULL REFERENCES rental (rental_id),
    amount       NUMERIC(5, 2) NOT NULL,
    payment_date TIMESTAMP     NOT NULL,
    UNIQUE (rental_id)
    -- Assuming one payment per rental
);

-- Add indexes to optimize queries
CREATE INDEX idx_customer_last_name ON customer (last_name);
CREATE INDEX idx_actor_last_name ON actor (last_name);
CREATE INDEX idx_film_title ON film (title);
CREATE INDEX idx_inventory_film_id ON inventory (film_id);
CREATE INDEX idx_rental_customer_id ON rental (customer_id);
CREATE INDEX idx_rental_inventory_id ON rental (inventory_id);
CREATE INDEX idx_payment_customer_id ON payment (customer_id);