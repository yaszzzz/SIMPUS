package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"simpus/config"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize random seed
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Error loading config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected. Starting seeding...")

	// 1. Seed Categories
	categoryIDs := seedCategories(db)

	// 2. Seed Authors
	authorIDs := seedAuthors(db)

	// 3. Seed Books
	bookIDs := seedBooks(db, categoryIDs, authorIDs)

	// 4. Seed Members
	memberIDs := seedMembers(db)

	// 5. Seed Borrowings
	seedBorrowings(db, memberIDs, bookIDs)

	// 6. Update Admin
	seedAdmin(db)

	log.Println("Seeding completed successfully!")
}

func seedCategories(db *sql.DB) []int {
	categories := []string{
		"Philosophy", "Psychology", "Religion", "Social Sciences",
		"Language", "Science", "Technology", "Arts", "Literature", "History",
	}

	var ids []int
	for _, cat := range categories {
		res, err := db.Exec("INSERT INTO categories (name, description) VALUES (?, ?) ON DUPLICATE KEY UPDATE name=name", cat, "Books about "+cat)
		if err != nil {
			log.Printf("Error inserting category %s: %v", cat, err)
			continue
		}
		id, _ := res.LastInsertId()
		if id == 0 {
			// If duplicate, query the ID
			var existingID int
			err := db.QueryRow("SELECT id FROM categories WHERE name = ?", cat).Scan(&existingID)
			if err == nil {
				ids = append(ids, existingID)
			}
		} else {
			ids = append(ids, int(id))
		}
	}
	log.Printf("Seeded %d categories", len(ids))
	return ids
}

func seedAuthors(db *sql.DB) []int {
	authors := []string{
		"Jane Austen", "Charles Dickens", "Mark Twain", "Leo Tolstoy",
		"Fyodor Dostoevsky", "Ernest Hemingway", "George Orwell", "J.R.R. Tolkien",
		"Agatha Christie", "Stephen King", "Gabriel Garcia Marquez", "Haruki Murakami",
	}

	var ids []int
	for _, auth := range authors {
		res, err := db.Exec("INSERT INTO authors (name, bio) VALUES (?, ?) ON DUPLICATE KEY UPDATE name=name", auth, "Famous author "+auth)
		if err != nil {
			log.Printf("Error inserting author %s: %v", auth, err)
			continue
		}
		id, _ := res.LastInsertId()
		// Since name isn't unique in schema, we assume we always insert or get ID.
		// Actually schema doesn't have unique constraint on author name, so we might duplicate if we run multiple times.
		// Let's just track IDs.
		ids = append(ids, int(id))
	}
	log.Printf("Seeded %d authors", len(ids))
	return ids
}

func seedBooks(db *sql.DB, catIDs, authIDs []int) []int {
	titles := []string{
		"The Great Adventure", "Mystery of the Blue Train", "Coding for Dummies",
		"History of the World", "The Art of War", "Silent Spring", "Cosmos",
		"Sapiens", "Thinking, Fast and Slow", "Dune", "Neuromancer",
		"Snow Crash", "Foundation", "Hyperion", "Starship Troopers",
	}

	var ids []int
	for i, title := range titles {
		isbn := fmt.Sprintf("978-%d-%d-%d-%d", rand.Intn(999), rand.Intn(99), rand.Intn(9999), rand.Intn(9))
		catID := catIDs[rand.Intn(len(catIDs))]
		authID := authIDs[rand.Intn(len(authIDs))]
		stock := rand.Intn(20) + 1

		query := `
			INSERT INTO books (isbn, title, category_id, author_id, publisher, publish_year, stock, available, description)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		res, err := db.Exec(query, isbn, title+fmt.Sprintf(" Vol %d", i), catID, authID, "Random House", 1990+rand.Intn(30), stock, stock, "A very interesting book.")
		if err != nil {
			log.Printf("Error inserting book %s: %v", title, err)
			continue
		}
		id, _ := res.LastInsertId()
		ids = append(ids, int(id))
	}
	log.Printf("Seeded %d books", len(ids))
	return ids
}

func seedMembers(db *sql.DB) []int {
	// Default password hash for 'password123'
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	types := []string{"mahasiswa", "guru", "karyawan"}

	var ids []int
	for i := 1; i <= 20; i++ {
		code := fmt.Sprintf("MEM%03d", i)
		name := fmt.Sprintf("Member %d", i)
		email := fmt.Sprintf("member%d@example.com", i)
		mType := types[rand.Intn(len(types))]

		query := `
			INSERT INTO members (member_code, name, email, password, phone, member_type, address)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=name
		`
		res, err := db.Exec(query, code, name, email, string(pwHash), "08123456789", mType, "Random Address")
		if err != nil {
			log.Printf("Error inserting member %s: %v", name, err)
			continue
		}
		id, _ := res.LastInsertId()
		if id == 0 {
			var existingID int
			db.QueryRow("SELECT id FROM members WHERE member_code = ?", code).Scan(&existingID)
			ids = append(ids, existingID)
		} else {
			ids = append(ids, int(id))
		}
	}
	log.Printf("Seeded %d members", len(ids))
	return ids
}

func seedBorrowings(db *sql.DB, memberIDs, bookIDs []int) {
	statuses := []string{"dipinjam", "dikembalikan", "terlambat"}

	count := 0
	for i := 0; i < 50; i++ {
		memberID := memberIDs[rand.Intn(len(memberIDs))]
		bookID := bookIDs[rand.Intn(len(bookIDs))]
		status := statuses[rand.Intn(len(statuses))]

		borrowDate := time.Now().AddDate(0, 0, -rand.Intn(30)) // Borrowed in last 30 days
		dueDate := borrowDate.AddDate(0, 0, 7)                 // Due in 7 days

		var returnDate *time.Time
		var fine float64 = 0

		if status == "dikembalikan" {
			rd := borrowDate.AddDate(0, 0, rand.Intn(10)) // Returned within 10 days
			returnDate = &rd
			if rd.After(dueDate) {
				daysLate := int(rd.Sub(dueDate).Hours() / 24)
				fine = float64(daysLate * 1000)
			}
		} else if status == "terlambat" {
			// Ensure it's late relative to now
			borrowDate = time.Now().AddDate(0, 0, -20)
			dueDate = borrowDate.AddDate(0, 0, 7)
			fine = float64(int(time.Since(dueDate).Hours()/24) * 1000)
		} else {
			// dipinjam - check if it's already late
			if time.Now().After(dueDate) {
				status = "terlambat" // Auto update status context
			}
		}

		query := `
			INSERT INTO borrowings (member_id, book_id, borrow_date, due_date, return_date, status, fine)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		_, err := db.Exec(query, memberID, bookID, borrowDate, dueDate, returnDate, status, fine)
		if err != nil {
			log.Printf("Error seeding borrowing: %v", err)
		} else {
			count++
		}
	}
	log.Printf("Seeded %d borrowings", count)
}

func seedAdmin(db *sql.DB) {
	// Hash for 'admin'
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	query := `
		INSERT INTO users (username, email, password, name, role)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE password = VALUES(password)
	`
	_, err := db.Exec(query, "admin", "admin@simpus.local", string(pwHash), "Administrator", "admin")
	if err != nil {
		log.Printf("Error updating admin user: %v", err)
	} else {
		log.Println("Admin user updated (username: admin, password: admin)")
	}
}
