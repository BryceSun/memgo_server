package databean

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Id       int64  `json:"id"`
}
