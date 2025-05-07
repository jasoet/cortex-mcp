-- 000002_seed_data.up.sql: Populate DVD Rental sample data (Simplified)

-- Store data with direct address fields
INSERT INTO store (store_name, address, address2, district, city, country, postal_code, phone)
VALUES ('Store 1 - Fortaleza', '1825 Main Ln', 'Suite 755', 'SP', 'Fortaleza', 'Brazil', '8196001', '3389083863'),
       ('Store 2 - Melbourne', '107 Elm Way', 'Suite 221', 'VIC', 'Melbourne', 'Australia', '16155', '9407816184'),
       ('Store 3 - Salvador', '5926 2nd Rd', NULL, 'SP', 'Salvador', 'Brazil', '3413164', '7525534192'),
       ('Store 4 - Glasgow', '7516 Oak Dr', 'Floor 983', 'Western Cape', 'Glasgow', 'United Kingdom', '775445', '048-449-4143'),
       ('Store 5 - Cape Town', '6197 Broadway Blvd', 'Suite 283', 'Ile-de-France', 'Cape Town', 'South Africa', '460307', '373-429-5950'),
       ('Store 6 - Marseille', '5158 Oak Ln', 'Apt. 863', 'Ile-de-France', 'Marseille', 'France', '13022', '9886563435'),
       ('Store 7 - Bangalore', '7518 Cedar Way', 'Suite 296', 'Karnataka', 'Bangalore', 'India', '3627', '3884482131'),
       ('Store 8 - Ottawa', '6691 Park Ave', 'Floor 986', 'ON', 'Ottawa', 'Canada', '7068', '1216855787'),
       ('Store 9 - Phoenix', '8839 2nd Blvd', 'Suite 897', 'AZ', 'Phoenix', 'United States', '48586', '7232690602'),
       ('Store 10 - Munich', '2139 Highland St', 'Unit 478', 'Bavaria', 'Munich', 'Germany', '571647', '3323864591');

-- Staff data with direct address fields
INSERT INTO staff (store_id, first_name, last_name, email, username, address, address2, district, city, country, postal_code, phone, active)
VALUES (1, 'Patricia', 'Anderson', 'patricia.anderson1@dvdrental.com', 'panderson1', '3864 Oak Rd', 'Floor 634', 'Delhi', 'New Delhi', 'India', '476383', '3324278458', TRUE),
       (1, 'Charles', 'Green', 'charles.green2@dvdrental.com', 'cgreen2', '6348 Cherry Ln', 'Apt. 802', 'Sichuan', 'Chengdu', 'China', '8770', '9760234479', TRUE),
       (1, 'Jennifer', 'Lopez', 'jennifer.lopez3@dvdrental.com', 'jlopez3', '8278 5th Way', 'Apt. 385', 'Maharashtra', 'Mumbai', 'India', '9062169', '2156497431', TRUE),
       (1, 'Melissa', 'Garcia', 'melissa.garcia4@dvdrental.com', 'mgarcia4', '3621 4th Blvd', 'Unit 751', 'Hamburg', 'Hamburg', 'Germany', '65844', '1269352169', TRUE),
       (1, 'Mary', 'Robinson', 'mary.robinson5@dvdrental.com', 'mrobinson5', '8500 Washington Rd', 'Apt. 991', 'Delhi', 'New Delhi', 'India', '36038', '913-296-1755', TRUE),
       (1, 'William', 'Johnson', 'william.johnson6@dvdrental.com', 'wjohnson6', '2229 3rd Ave', 'Unit 218', 'Quebec', 'Montreal', 'Canada', '72752', '3689668800', TRUE),
       (1, 'Sandra', 'Perez', 'sandra.perez7@dvdrental.com', 'sperez7', '8086 Hill Rd', 'Apt. 553', 'Delhi', 'New Delhi', 'India', '4170299', '828-526-4402', TRUE),
       (1, 'Sarah', 'Lee', 'sarah.lee8@dvdrental.com', 'slee8', '9080 Highland Ave', 'Suite 510', 'Brittany', 'Rennes', 'France', '61018', '5954847269', TRUE),
       (1, 'Ronald', 'Turner', 'ronald.turner9@dvdrental.com', 'rturner9', '6629 Maple Ave', 'Unit 636', 'Sichuan', 'Chengdu', 'China', '3079', '543-472-4678', TRUE),
       (2, 'Betty', 'Walker', 'betty.walker10@dvdrental.com', 'bwalker10', '8492 Hill Way', 'Unit 795', 'Gauteng', 'Johannesburg', 'South Africa', '16710', '107-957-4514', TRUE),
       (2, 'Josef', 'Baker', 'josef.baker11@dvdrental.com', 'jbaker11', '6932 Broadway Ln', 'Unit 843', 'NSW', 'Sydney', 'Australia', '15364', '354-892-6428', TRUE),
       (2, 'Charles', 'Clark', 'charles.clark12@dvdrental.com', 'cclark12', '7379 Park Dr', 'Unit 20', 'Gujarat', 'Ahmedabad', 'India', '92898', '7525821571', TRUE),
       (2, 'Susan', 'Adams', 'susan.adams13@dvdrental.com', 'sadams13', '3661 Cherry St', 'Floor 69', 'Guangdong', 'Guangzhou', 'China', '25435', '6038537534', TRUE),
       (2, 'Deborah', 'Rodriguez', 'deborah.rodriguez14@dvdrental.com', 'drodriguez14', '9760 Main Ln', 'Apt. 830', 'Northern Ireland', 'Belfast', 'United Kingdom', '161184', '4995500281', TRUE),
       (2, 'Anthony', 'Evans', 'anthony.evans15@dvdrental.com', 'aevans15', '5226 Maple Blvd', 'Floor 601', 'Bavaria', 'Munich', 'Germany', '91984', '1934951876', TRUE),
       (2, 'Joseph', 'Perez', 'joseph.perez16@dvdrental.com', 'jperez16', '9950 Chestnut St', 'Floor 583', 'Buenos Aires', 'Buenos Aires', 'Argentina', '96366', '783-697-6285', TRUE),
       (2, 'Dorothy', 'Jackson', 'dorothy.jackson17@dvdrental.com', 'djackson17', '8731 Cherry Ln', 'Floor 308', 'Bavaria', 'Munich', 'Germany', '78768', '324-832-3791', TRUE),
       (2, 'Andrew', 'Harris', 'andrew.harris18@dvdrental.com', 'aharris18', '3660 Oak Ln', 'Floor 467', 'Brandenburg', 'Potsdam', 'Germany', '5868', '2595470585', TRUE),
       (3, 'Jessica', 'Jones', 'jessica.jones19@dvdrental.com', 'jjones19', '5964 River Dr', 'Suite 991', 'Texas', 'Houston', 'United States', '48723', '942-698-2706', TRUE),
       (3, 'Donald', 'Moore', 'donald.moore20@dvdrental.com', 'dmoore20', '5524 Maple Blvd', 'Floor 473', 'Andalusia', 'Seville', 'Spain', '20032', '576-131-1982', TRUE);

-- More staff data (truncated for brevity)
-- Add more staff data as needed

-- Customer data with direct address fields
INSERT INTO customer (store_id, first_name, last_name, email, address, address2, district, city, country, postal_code, phone, active, create_date)
VALUES (7, 'Ashley', 'Martinez', 'ashley.martinez1@example.com', '1476 Highland Way', 'Suite 806', 'Karnataka', 'Bangalore', 'India', '8035', '647-542-4065', TRUE, '2022-12-29'),
       (9, 'Mary', 'Hill', 'mary.hill2@example.com', '5714 Cherry Way', 'Unit 496', 'NSW', 'Sydney', 'Australia', '1107360', '5605285906', TRUE, '2023-07-17'),
       (6, 'Betty', 'Green', 'betty.green3@example.com', '9372 Cherry Ln', 'Apt. 374', 'Bavaria', 'Munich', 'Germany', '58645', '7124674551', TRUE, '2023-05-10'),
       (1, 'Donna', 'Hall', 'donna.hall4@example.com', '1270 5th Rd', 'Floor 671', 'Delhi', 'New Delhi', 'India', '21571', '7043629497', TRUE, '2023-12-28'),
       (5, 'Jessica', 'Scott', 'jessica.scott5@example.com', '6650 Broadway Blvd', 'Suite 488', 'Bavaria', 'Munich', 'Germany', '1367', '1854733995', TRUE, '2022-07-17'),
       (7, 'Karen', 'Hernandez', 'karen.hernandez6@example.com', '3524 Pine Way', 'Floor 908', 'Delhi', 'New Delhi', 'India', '3017044', '484-264-2864', TRUE, '2022-06-29'),
       (2, 'James', 'Allen', 'james.allen7@example.com', '1226 Main Blvd', 'Apt. 37', 'Rio de Janeiro', 'Rio de Janeiro', 'Brazil', '3294', '1063747264', TRUE, '2024-12-31'),
       (5, 'Edward', 'Brown', 'edward.brown8@example.com', '810 4th St', 'Unit 609', 'São Paulo', 'São Paulo', 'Brazil', '251587', '6674525389', TRUE, '2024-03-13'),
       (10, 'Joseph', 'Taylor', 'joseph.taylor9@example.com', '1071 Chestnut Terrace', 'Floor 885', 'Hamburg', 'Hamburg', 'Germany', '513643', '854-614-1370', TRUE, '2023-07-26'),
       (1, 'Dorothy', 'Young', 'dorothy.young10@example.com', '5777 Main Blvd', 'Apt. 469', 'Karnataka', 'Bangalore', 'India', '9621927', '8543136674', TRUE, '2023-06-09');

-- More customer data (truncated for brevity)
-- Add more customer data as needed

-- Category data
INSERT INTO category (name)
VALUES ('Action'),
       ('Comedy'),
       ('Drama'),
       ('Horror'),
       ('Romance'),
       ('Sci-Fi'),
       ('Documentary'),
       ('Thriller'),
       ('Animation'),
       ('Fantasy');

-- Film data
INSERT INTO film (title, release_year, length, category_id)
VALUES ('Red Dream', 1983, 101, 3),
       ('Return of the Child', 1951, 147, 5),
       ('The Silent Sky', 1962, 150, 3),
       ('Return of the Deadly Revenge', 1976, 158, 1),
       ('Return of the World', 2012, 142, 1),
       ('Return of the Fear', 1988, 170, 4),
       ('The Hero of Fear', 2001, 129, 1),
       ('The Green Fire', 1963, 158, 9),
       ('The Journey of Blade', 1979, 140, 10),
       ('The Star of Time', 1954, 103, 6);

-- More film data (truncated for brevity)
-- Add more film data as needed

-- Actor data
INSERT INTO actor (first_name, last_name)
VALUES ('Melissa', 'Evans'),
       ('Margaret', 'Harris'),
       ('Ronald', 'Scott'),
       ('Jessica', 'Johnson'),
       ('Donald', 'Walker'),
       ('Mary', 'Lee'),
       ('Margaret', 'Ramirez'),
       ('James', 'Taylor'),
       ('Paul', 'White'),
       ('William', 'Lopez');

-- More actor data (truncated for brevity)
-- Add more actor data as needed

-- Film-Actor relationships
INSERT INTO film_actors (film_id, actor_id)
VALUES (1, 1),
       (1, 2),
       (2, 3),
       (2, 4),
       (2, 5),
       (3, 6),
       (3, 7),
       (3, 8),
       (4, 9),
       (4, 10);

-- More film-actor relationships (truncated for brevity)
-- Add more film-actor relationships as needed

-- Inventory data
INSERT INTO inventory (film_id, store_id)
VALUES (1, 1),
       (2, 1),
       (3, 1),
       (4, 2),
       (5, 2),
       (6, 2),
       (7, 3),
       (8, 3),
       (9, 3),
       (10, 4);

-- More inventory data (truncated for brevity)
-- Add more inventory data as needed

-- Rental transactions (truncated for brevity)
-- Add rental transactions as needed

-- Payment transactions (truncated for brevity)
-- Add payment transactions as needed