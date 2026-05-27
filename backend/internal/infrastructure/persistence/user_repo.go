package persistence

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/security"
	"github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc"
)



type userRepository struct {
	db      *PostgresDB
	queries *sqlc.Queries
	encrypt *security.EncryptionService
}

func NewUserRepository(db *PostgresDB, encrypt *security.EncryptionService) repository.UserRepository {
	return &userRepository{
		db:      db,
		queries: db.Queries,
		encrypt: encrypt,
	}
}

func hashEmail(email string) string {
	h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	return fmt.Sprintf("%x", h)
}

func (r *userRepository) encryptUser(user *entity.User) (emailHash string, encryptedEmail, encryptedFullName string, encryptedPhone pgtype.Text, err error) {
	emailHash = hashEmail(user.Email)

	encryptedEmail, err = r.encrypt.Encrypt(user.Email)
	if err != nil {
		return "", "", "", pgtype.Text{}, fmt.Errorf("encrypt email: %w", err)
	}

	encryptedFullName, err = r.encrypt.Encrypt(user.FullName)
	if err != nil {
		return "", "", "", pgtype.Text{}, fmt.Errorf("encrypt full name: %w", err)
	}

	if user.Phone != nil {
		var encPhone string
		encPhone, err = r.encrypt.Encrypt(*user.Phone)
		if err != nil {
			return "", "", "", pgtype.Text{}, fmt.Errorf("encrypt phone: %w", err)
		}
		encryptedPhone = pgtype.Text{String: encPhone, Valid: true}
	}

	return
}

func (r *userRepository) decryptUser(sqlcUser sqlc.User) (*entity.User, error) {
	email, err := r.encrypt.Decrypt(sqlcUser.Email)
	if err != nil {
		return nil, fmt.Errorf("decrypt email: %w", err)
	}

	fullName, err := r.encrypt.Decrypt(sqlcUser.FullName)
	if err != nil {
		return nil, fmt.Errorf("decrypt full name: %w", err)
	}

	var phone *string
	if sqlcUser.Phone.Valid {
		decPhone, err := r.encrypt.Decrypt(sqlcUser.Phone.String)
		if err != nil {
			return nil, fmt.Errorf("decrypt phone: %w", err)
		}
		phone = &decPhone
	}

	return &entity.User{
		ID:           sqlcUser.ID,
		Email:        email,
		PasswordHash: sqlcUser.PasswordHash,
		FullName:     fullName,
		Phone:        phone,
		Role:         sqlcUser.Role,
		Active:       sqlcUser.Active,
		SuperAdmin:   sqlcUser.SuperAdmin,
		CreatedAt:    sqlcUser.CreatedAt,
		UpdatedAt:    sqlcUser.UpdatedAt,
	}, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	emailHash, encryptedEmail, encryptedFullName, encryptedPhone, err := r.encryptUser(user)
	if err != nil {
		return err
	}

	role := user.Role
	if role == "" {
		role = "user"
	}

	sqlcUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        encryptedEmail,
		EmailHash:    emailHash,
		PasswordHash: user.PasswordHash,
		FullName:     encryptedFullName,
		Phone:        encryptedPhone,
		Role:         role,
		SuperAdmin:   user.SuperAdmin,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	decrypted, err := r.decryptUser(sqlcUser)
	if err != nil {
		return err
	}
	*user = *decrypted

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	sqlcUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return r.decryptUser(sqlcUser)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	emailHash := hashEmail(email)

	sqlcUser, err := r.queries.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return r.decryptUser(sqlcUser)
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]entity.User, int, error) {
	total, err := r.queries.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	sqlcUsers, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	users := make([]entity.User, len(sqlcUsers))
	for i, u := range sqlcUsers {
		decrypted, err := r.decryptUser(u)
		if err != nil {
			return nil, 0, fmt.Errorf("decrypt user at index %d: %w", i, err)
		}
		users[i] = *decrypted
	}

	return users, int(total), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	emailHash, encryptedEmail, encryptedFullName, encryptedPhone, err := r.encryptUser(user)
	if err != nil {
		return err
	}

	role := user.Role
	if role == "" {
		role = "user"
	}

	sqlcUser, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        user.ID,
		Email:     encryptedEmail,
		EmailHash: emailHash,
		FullName:  encryptedFullName,
		Phone:     encryptedPhone,
		Role:      role,
		Active:    user.Active,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("update user: %w", err)
	}

	decrypted, err := r.decryptUser(sqlcUser)
	if err != nil {
		return err
	}
	*user = *decrypted

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	_, err := r.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}

func (r *userRepository) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	_, err := r.queries.SetUserActive(ctx, sqlc.SetUserActiveParams{
		ID:     id,
		Active: active,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("set user active: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
