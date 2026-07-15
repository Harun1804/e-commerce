package user

import "harun1804/e-commerce/modules/access/models"

type UserResponse struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

func NewUserResponseList(user models.User) UserResponse {
	return UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
