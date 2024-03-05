package entity

type User struct {
	UserID   int64
	Name     string
	Password string
	Email    string
	Mobile   string
	Avatar   string
	Account  string

	ActiveInfo *ActiveInfo
}

func NewUser(userID int64, name, password, email, mobile, avatar string) *User {
	return &User{
		UserID:   userID,
		Name:     name,
		Password: password,
		Email:    email,
		Mobile:   mobile,
		Avatar:   avatar,
	}
}
