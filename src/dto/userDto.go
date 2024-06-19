package dto

import "gin_demo/src/entity"

type UserDto struct {
	Id        uint
	Name      string
	Telephone string
}

func ToUserDto(user entity.User) *UserDto {
	return &UserDto{
		Id:        user.ID,
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
