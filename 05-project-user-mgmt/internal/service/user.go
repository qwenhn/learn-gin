package service

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/models"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/repository"
	"github.com/qwenhn/gin-restful-api/05-project-user-mgmt/internal/utils"
)

type UserService interface {
	GetAllUsers(search string, page, limit int) ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetByUUID(uuid string) (models.User, error)
	UpdateUser(uuid string, user models.User) (models.User, error)
	DeleteUser(uuid string) error
}

type userService struct {
	repo repository.UserRepository
}

func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return models.User{}, utils.NewError("Email already exist", utils.ErrCodeConflict)
	}

	user.UUID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "Unable to hash password", utils.ErrCodeInternal)
	}

	user.Password = string(hashedPassword)

	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "Unable to create user", utils.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.WrapError(err, "Unable to delete user", utils.ErrCodeInternal)
	}

	return nil
}

func (us *userService) GetAllUsers(search string, page int, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "Unable to fetch users", utils.ErrCodeInternal)
	}

	var filteredUsers []models.User
	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)

			if strings.Contains(name, search) || strings.Contains(email, search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		filteredUsers = users
	}

	start := (page - 1) * limit
	if start > len(filteredUsers) {
		return []models.User{}, nil
	}

	end := start + limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	return filteredUsers[start:end], nil
}

func (us *userService) GetByUUID(uuid string) (models.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	return user, nil
}

func (us *userService) UpdateUser(uuid string, user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if u, exist := us.repo.FindByEmail(user.Email); exist && u.UUID != uuid {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}

	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	currentUser.Name = user.Name
	currentUser.Email = user.Email
	currentUser.Age = user.Age
	currentUser.Status = user.Status
	currentUser.Level = user.Level

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "Unable to hash password", utils.ErrCodeInternal)
		}

		currentUser.Password = string(hashedPassword)
	}

	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "Unable to update user", utils.ErrCodeInternal)
	}

	return currentUser, nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
