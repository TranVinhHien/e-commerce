package db

import services "new-project/services/entity"

func (u *Users) Convert() services.Users {
	return services.Users{
		Username: u.Username,
		Password: u.Password,
		FullName: u.FullName,
		IsActive: u.IsActive,
		CreateAt: u.CreateAt,
	}
}
