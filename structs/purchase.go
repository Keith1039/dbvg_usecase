package structs

type Purchase struct {
	Username  string  `db:"username"`
	ItemName  string  `db:"name"`
	ItemDesc  string  `db:"description"`
	ItemPrice float64 `db:"price"`
}
