package service

import (
	"log"
	"pro1/dto"
	"pro1/entity"
	"pro1/repository"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredental(email, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}

func (s *authService) VerifyCredental(email, password string) interface{} {
	res := s.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		compare_password := comparePassword(v.Password, []byte(password))
		if v.Email == email && compare_password {
			return res
		}
		return false
	}

	return false
}

func (s *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	res := s.userRepository.InsertUser(userToCreate)
	return res
}

func (s *authService) FindByEmail(email string) entity.User {
	return s.userRepository.FindByEmail(email)
}

func (s *authService) IsDuplicateEmail(email string) bool {
	res := s.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashed_password string, password []byte) bool {
	byteHash := []byte(hashed_password)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
