package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/myinquisitor/backend/internal/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func hashEmail(email string) string {
	h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	return fmt.Sprintf("%x", h)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found")
	}

	encKey := os.Getenv("ENCRYPTION_KEY")
	if encKey == "" {
		log.Fatal("ENCRYPTION_KEY is required in .env")
	}

	encSvc, err := security.NewEncryptionService(encKey)
	if err != nil {
		log.Fatalf("invalid encryption key: %v", err)
	}

	email := "admin@myinquisitor.app"
	password := "admin123"
	fullName := "Super Admin"

	if len(os.Args) >= 2 {
		email = os.Args[1]
	}
	if len(os.Args) >= 3 {
		password = os.Args[2]
	}
	if len(os.Args) >= 4 {
		fullName = os.Args[3]
	}

	emailSHA := hashEmail(email)

	encEmail, err := encSvc.Encrypt(email)
	if err != nil {
		log.Fatalf("encrypt email: %v", err)
	}

	encFullName, err := encSvc.Encrypt(fullName)
	if err != nil {
		log.Fatalf("encrypt full name: %v", err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	fmt.Printf("INSERT INTO users (id, email, email_hash, password_hash, full_name, role, active, super_admin, created_at, updated_at) \n")
	fmt.Printf("VALUES (\n")
	fmt.Printf("  gen_random_uuid(),\n")
	fmt.Printf("  '%s',\n", encEmail)
	fmt.Printf("  '%s',\n", emailSHA)
	fmt.Printf("  '%s',\n", string(passHash))
	fmt.Printf("  '%s',\n", encFullName)
	fmt.Printf("  'super_admin',\n")
	fmt.Printf("  true,\n")
	fmt.Printf("  true,\n")
	fmt.Printf("  now(),\n")
	fmt.Printf("  now()\n")
	fmt.Printf(");\n")
}
