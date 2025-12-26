package database

func GetMigrationSchema() string {
	query := `
		-- migrations/schema.sql

		-- Libraries Table
		CREATE TABLE IF NOT EXISTS libraries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			address TEXT,
			city TEXT,
			state TEXT,
			zip_code TEXT,
			country TEXT,
			phone TEXT,
			email TEXT,
			website TEXT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);

		-- Users table
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			dni TEXT UNIQUE NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE,
			phone TEXT,
			address TEXT,
			user_type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'active',
			registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Authors table
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			biography TEXT,
			nationality TEXT,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Publishers table
		CREATE TABLE IF NOT EXISTS publishers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			country TEXT,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Categories table
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			description TEXT,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Library zones table
		CREATE TABLE IF NOT EXISTS library_zones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			floor INTEGER NOT NULL,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Shelves table
		CREATE TABLE IF NOT EXISTS shelves (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL,
			zone_id INTEGER NOT NULL,
			description TEXT,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (zone_id) REFERENCES library_zones(id),
			FOREIGN KEY (library_id) REFERENCES libraries(id),
			UNIQUE(code, zone_id)
		);

		-- Books table
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			isbn TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			subtitle TEXT,
			edition TEXT,
			publication_year INTEGER,
			language TEXT DEFAULT 'Spanish',
			pages INTEGER,
			synopsis TEXT,
			publisher_id INTEGER,
			shelf_id INTEGER,
			status TEXT NOT NULL DEFAULT 'available',
			registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (publisher_id) REFERENCES publishers(id),
			FOREIGN KEY (shelf_id) REFERENCES shelves(id),
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Book authors junction table
		CREATE TABLE IF NOT EXISTS book_authors (
			book_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			position INTEGER DEFAULT 1,
			library_id INTEGER NOT NULL,
			PRIMARY KEY (book_id, author_id),
			FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Book categories junction table
		CREATE TABLE IF NOT EXISTS book_categories (
			book_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			library_id INTEGER NOT NULL,
			PRIMARY KEY (book_id, category_id),
			FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Copies table
		CREATE TABLE IF NOT EXISTS copies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			book_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'available',
			condition TEXT DEFAULT 'good',
			acquisition_date TIMESTAMP,
			purchase_price REAL,
			notes TEXT,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Loans table
		CREATE TABLE IF NOT EXISTS loans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			loan_code TEXT UNIQUE NOT NULL,
			user_id INTEGER NOT NULL,
			copy_id INTEGER NOT NULL,
			loan_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			due_date TIMESTAMP NOT NULL,
			return_date TIMESTAMP,
			status TEXT NOT NULL DEFAULT 'active',
			loan_days INTEGER NOT NULL DEFAULT 15,
			renewals INTEGER DEFAULT 0,
			notes TEXT,
			librarian_id INTEGER,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (copy_id) REFERENCES copies(id),
			FOREIGN KEY (librarian_id) REFERENCES users(id),
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Reservations table
		CREATE TABLE IF NOT EXISTS reservations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			book_id INTEGER NOT NULL,
			reservation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expiration_date TIMESTAMP NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			priority INTEGER DEFAULT 1,
			notified BOOLEAN DEFAULT 0,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (book_id) REFERENCES books(id),
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Fines table
		CREATE TABLE IF NOT EXISTS fines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			loan_id INTEGER,
			reason TEXT NOT NULL,
			amount REAL NOT NULL,
			generated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			payment_date TIMESTAMP,
			status TEXT NOT NULL DEFAULT 'pending',
			notes TEXT,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (loan_id) REFERENCES loans(id),
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Configuration table
		CREATE TABLE IF NOT EXISTS configuration (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student_loan_days INTEGER DEFAULT 15,
			teacher_loan_days INTEGER DEFAULT 30,
			max_renewals INTEGER DEFAULT 2,
			max_books_per_loan INTEGER DEFAULT 5,
			fine_per_day REAL DEFAULT 0.50,
			reservation_days INTEGER DEFAULT 3,
			grace_days INTEGER DEFAULT 2,
			library_id INTEGER NOT NULL,
			FOREIGN KEY (library_id) REFERENCES libraries(id)
		);

		-- Create indexes for better performance
		CREATE INDEX IF NOT EXISTS idx_libraries_name ON libraries(name);
		CREATE INDEX IF NOT EXISTS idx_libraries_username ON libraries(username);
		CREATE INDEX IF NOT EXISTS idx_users_code ON users(code);
		CREATE INDEX IF NOT EXISTS idx_users_dni ON users(dni);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
		CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
		CREATE INDEX IF NOT EXISTS idx_books_status ON books(status);
		CREATE INDEX IF NOT EXISTS idx_copies_code ON copies(code);
		CREATE INDEX IF NOT EXISTS idx_copies_status ON copies(status);
		CREATE INDEX IF NOT EXISTS idx_loans_user_id ON loans(user_id);
		CREATE INDEX IF NOT EXISTS idx_loans_copy_id ON loans(copy_id);
		CREATE INDEX IF NOT EXISTS idx_loans_status ON loans(status);
		CREATE INDEX IF NOT EXISTS idx_loans_due_date ON loans(due_date);
		CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations(user_id);
		CREATE INDEX IF NOT EXISTS idx_reservations_book_id ON reservations(book_id);
		CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
		CREATE INDEX IF NOT EXISTS idx_fines_user_id ON fines(user_id);
		CREATE INDEX IF NOT EXISTS idx_fines_status ON fines(status);
	`

	return query
}
