package dto

type UpdateProfileInput struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type UpdateProfileOutput struct {
	User UserDTO `json:"user"`
}
