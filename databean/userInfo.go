package databean

type Loginer interface {
	GetName() string
	GetPassword() string
}

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Id       int64  `json:"id"`
}

func (user *UserInfo) GetName() string {
	username, email, mobile := user.Username, user.Email, user.Mobile
	if len(username) != 0 {
		return username
	}
	if len(mobile) != 0 {
		return mobile
	}
	if len(email) != 0 {
		return email
	}
	return ""
}

func (user *UserInfo) GetPassword() string {
	return user.Password
}
