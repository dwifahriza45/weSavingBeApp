package users

type User struct {
	ID       int64  `db:"id" json:"id"`
	USER_ID  string `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Fullname string `db:"fullname" json:"fullname"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password_hash" json:"password"`
}
