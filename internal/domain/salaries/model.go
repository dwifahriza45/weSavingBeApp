package salaries

type Salaries struct {
	ID          int64  `db:"id" json:"id"`
	SalaryID    string `db:"salary_id" json:"salary_id"`
	UserID      string `db:"user_id" json:"user_id"`
	Amount      string `db:"amount" json:"amount"`
	Source      string `db:"source" json:"source"`
	Description string `db:"description" json:"description"`
	ReceivedAt  string `db:"received_at" json:"received_at"`
}
