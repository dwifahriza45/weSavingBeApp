package categories

type Categories struct {
	ID          int64  `db:"id" json:"id"`
	CATEGORY_ID string `db:"category_id" json:"category_id"`
	USER_ID     string `db:"user_id" json:"user_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}
