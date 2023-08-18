package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func numericToFloat64(value *pgtype.Numeric) (float64, error) {
	val, err := value.Float64Value()
	if err != nil || !val.Valid {
		return 0, err
	}

	return val.Float64, nil
}

func numericToInt64(value *pgtype.Numeric) (int64, error) {
	val, err := value.Int64Value()
	if err != nil || !val.Valid {
		return 0, err
	}

	return val.Int64, nil
}
