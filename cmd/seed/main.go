package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Default admin credentials
	email := "admin@example.com"
	password := "admin123"
	role := "admin"

	// Allow override via args
	if len(os.Args) > 1 {
		email = os.Args[1]
	}
	if len(os.Args) > 2 {
		password = os.Args[2]
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Check if admin already exists
	var exists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check existing admin: %v", err)
	}
	if exists {
		fmt.Printf("⚠️  Admin user '%s' already exists, skipping.\n", email)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Insert admin user
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, role, is_active, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, true, NOW(), NOW())`,
		email, string(hashedPassword), role,
	)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Printf("   Email:    %s\n", email)
	fmt.Printf("   Password: %s\n", password)
	fmt.Printf("   Role:     %s\n", role)
}
