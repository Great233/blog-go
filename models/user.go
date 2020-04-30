package models

type User struct {
	Id        uint   `json:"-"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	LastLogin int    `json:"last_login"`
}

func GetUserByUsername(username string) (*User, error) {
	var user User

	err := Query().Where("username=?", username).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}
