
 package service

import (
	"log"
	"pro1/dto"
	"pro1/entity"
	"pro1/repository"

	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (s *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed to %v: ", err)
	}

	updateUser := s.userRepository.UpdateUser(userToUpdate)
	return updateUser
}

func (s *userService) Profile(userID string) entity.User {
	return s.userRepository.ProfileUser(userID)
}
