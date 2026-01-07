INSERT INTO genres (name) VALUES 
('Science Fiction'), ('Fantasy'), ('Classic'), ('Non-Fiction'), ('History');

INSERT INTO books (isbn, title, author, daily_rental_fee, total_stock, available_stock) VALUES 
('978-0441172719', 'Dune', 'Frank Herbert', 7500, 5, 5),
('978-0099590088', 'Sapiens', 'Yuval Noah Harari', 9000, 7, 7),
('978-0451524935', '1984', 'George Orwell', 5000, 10, 10);

INSERT INTO book_genres (book_id, genre_id) VALUES (1, 1), (1, 2);
INSERT INTO book_genres (book_id, genre_id) VALUES (2, 4), (2, 5);
INSERT INTO book_genres (book_id, genre_id) VALUES (3, 1), (3, 3);