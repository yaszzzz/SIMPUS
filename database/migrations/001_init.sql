-- SIMPUS Database Schema
-- Sistem Informasi Manajemen Perpustakaan

-- Drop tables if exists (for fresh install)
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS borrowings;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS users;

-- Users table (Admin/Staff)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role ENUM('admin', 'staff') DEFAULT 'staff',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Members table (Mahasiswa/Guru/Karyawan)
CREATE TABLE members (
    id INT PRIMARY KEY AUTO_INCREMENT,
    member_code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    member_type ENUM('mahasiswa', 'guru', 'karyawan') NOT NULL,
    address TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Categories table
CREATE TABLE categories (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Authors table
CREATE TABLE authors (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Books table
CREATE TABLE books (
    id INT PRIMARY KEY AUTO_INCREMENT,
    isbn VARCHAR(20) UNIQUE,
    title VARCHAR(255) NOT NULL,
    category_id INT,
    author_id INT,
    publisher VARCHAR(100),
    publish_year INT,
    stock INT DEFAULT 0,
    available INT DEFAULT 0,
    cover_image VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE SET NULL
);

-- Borrowings table
CREATE TABLE borrowings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    member_id INT NOT NULL,
    book_id INT NOT NULL,
    user_id INT,
    borrow_date DATE NOT NULL,
    due_date DATE NOT NULL,
    return_date DATE,
    status ENUM('dipinjam', 'dikembalikan', 'terlambat') DEFAULT 'dipinjam',
    fine DECIMAL(10, 2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Notifications table
CREATE TABLE notifications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    borrowing_id INT,
    member_id INT NOT NULL,
    type ENUM('keterlambatan', 'pengingat', 'info') DEFAULT 'info',
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (borrowing_id) REFERENCES borrowings(id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_books_category ON books(category_id);
CREATE INDEX idx_books_author ON books(author_id);
CREATE INDEX idx_borrowings_member ON borrowings(member_id);
CREATE INDEX idx_borrowings_book ON borrowings(book_id);
CREATE INDEX idx_borrowings_status ON borrowings(status);
CREATE INDEX idx_borrowings_due_date ON borrowings(due_date);
CREATE INDEX idx_notifications_member ON notifications(member_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);

-- Insert default admin user (password: admin123)
INSERT INTO users (username, email, password, name, role) VALUES
('admin', 'admin@simpus.local', '$2a$10$N9qo8uLOickgx2ZMRZoMye.H9p4FxP7j1FQfXR2X9j5.3MLqImZ4a', 'Administrator', 'admin');

-- Insert sample categories
INSERT INTO categories (name, description) VALUES
('Fiksi', 'Buku-buku cerita fiksi dan novel'),
('Non-Fiksi', 'Buku-buku pengetahuan dan fakta'),
('Pendidikan', 'Buku pelajaran dan akademik'),
('Teknologi', 'Buku tentang teknologi dan komputer'),
('Sejarah', 'Buku tentang sejarah dan budaya'),
('Sains', 'Buku ilmu pengetahuan alam');

-- Insert sample authors
INSERT INTO authors (name, bio) VALUES
('Andrea Hirata', 'Penulis novel Laskar Pelangi'),
('Tere Liye', 'Penulis novel populer Indonesia'),
('Pramoedya Ananta Toer', 'Sastrawan Indonesia terkenal'),
('J.K. Rowling', 'Penulis seri Harry Potter'),
('Robert C. Martin', 'Penulis buku Clean Code');

-- Insert sample books
INSERT INTO books (isbn, title, category_id, author_id, publisher, publish_year, stock, available, description) VALUES
('978-602-03-1234-5', 'Laskar Pelangi', 1, 1, 'Bentang Pustaka', 2005, 5, 5, 'Novel tentang perjuangan anak-anak Belitung dalam menempuh pendidikan'),
('978-602-03-2345-6', 'Bumi', 1, 2, 'Gramedia', 2014, 3, 3, 'Novel fantasi pertama dari serial Bumi'),
('978-602-03-3456-7', 'Bumi Manusia', 1, 3, 'Hasta Mitra', 1980, 4, 4, 'Novel sejarah tentang perjuangan di era kolonial'),
('978-0-13-235088-4', 'Clean Code', 4, 5, 'Prentice Hall', 2008, 2, 2, 'Panduan menulis kode yang bersih dan mudah dipelihara'),
('978-602-03-4567-8', 'Bulan', 1, 2, 'Gramedia', 2015, 3, 3, 'Novel fantasi kedua dari serial Bumi');

-- Insert sample members
INSERT INTO members (member_code, name, email, password, phone, member_type, address) VALUES
('MHS001', 'Budi Santoso', 'budi@student.ac.id', '$2a$10$N9qo8uLOickgx2ZMRZoMye.H9p4FxP7j1FQfXR2X9j5.3MLqImZ4a', '081234567890', 'mahasiswa', 'Jl. Pendidikan No. 1'),
('MHS002', 'Siti Aminah', 'siti@student.ac.id', '$2a$10$N9qo8uLOickgx2ZMRZoMye.H9p4FxP7j1FQfXR2X9j5.3MLqImZ4a', '081234567891', 'mahasiswa', 'Jl. Ilmu No. 2'),
('GRU001', 'Drs. Ahmad Wijaya', 'ahmad@school.ac.id', '$2a$10$N9qo8uLOickgx2ZMRZoMye.H9p4FxP7j1FQfXR2X9j5.3MLqImZ4a', '081234567892', 'guru', 'Jl. Guru No. 3'),
('KRY001', 'Dewi Lestari', 'dewi@office.ac.id', '$2a$10$N9qo8uLOickgx2ZMRZoMye.H9p4FxP7j1FQfXR2X9j5.3MLqImZ4a', '081234567893', 'karyawan', 'Jl. Kantor No. 4');
