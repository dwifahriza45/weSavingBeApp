package cetegoriesbudget

type CategoriesBudget struct {
	ID              int64  `db:"id" json:"id"`
	BUDGET_ID       string `db:"budget_id" json:"budget_id"`
	USER_ID         string `db:"user_id" json:"user_id"`
	CATEGORY_ID     string `db:"cateogry_id" json:"cateogry_id"`
	AllocatedAmount string `db:"allocated_amount" json:"allocated_amount"`
	UsedAmount      string `db:"used_amount" json:"used_amount"`
	Period          string `db:"period" json:"period"`
}
