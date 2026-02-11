package service

import (
	"context"
	"testing"
	"time"

	"github.com/yang/go-learning-backend/gin-backend/internal/auth"
	"github.com/yang/go-learning-backend/gin-backend/internal/model"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
)

type fakeUserRepo struct {
	usersByEmail map[string]model.User
	usersByID    map[string]model.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		usersByEmail: map[string]model.User{},
		usersByID:    map[string]model.User{},
	}
}

func (r *fakeUserRepo) Create(_ context.Context, user model.User) (model.User, error) {
	if _, ok := r.usersByEmail[user.Email]; ok {
		return model.User{}, repository.ErrConflict
	}
	r.usersByEmail[user.Email] = user
	r.usersByID[user.ID] = user
	return user, nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (model.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return model.User{}, repository.ErrNotFound
	}
	return user, nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, id string) (model.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return model.User{}, repository.ErrNotFound
	}
	return user, nil
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	repo := newFakeUserRepo()
	jwt := auth.NewJWTManager("test-secret-123", "test", 10*time.Minute)
	svc := NewAuthService(repo, jwt, 8)

	user, err := svc.Register(context.Background(), "u1@example.com", "password123")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.ID == "" {
		t.Fatal("expected user id")
	}

	token, gotUser, err := svc.Login(context.Background(), "u1@example.com", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if token == "" || gotUser.Email != "u1@example.com" {
		t.Fatalf("unexpected login result token=%q user=%+v", token, gotUser)
	}
}

func TestAuthService_LoginInvalidPassword(t *testing.T) {
	repo := newFakeUserRepo()
	jwt := auth.NewJWTManager("test-secret-123", "test", 10*time.Minute)
	svc := NewAuthService(repo, jwt, 8)

	if _, err := svc.Register(context.Background(), "u1@example.com", "password123"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if _, _, err := svc.Login(context.Background(), "u1@example.com", "bad"); err == nil {
		t.Fatal("expected login error for invalid password")
	}
}
