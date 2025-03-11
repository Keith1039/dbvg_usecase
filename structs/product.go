package structs

import "github.com/jackc/pgx/v5/pgtype"

type Product struct {
	ID          pgtype.UUID `db:"id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	Price       float64     `db:"price"`
}
