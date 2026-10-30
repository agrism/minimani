package db

import (
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

// Text returns a valid pgtype.Text
func Text(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

// NumericFromFloat converts float64 to pgtype.Numeric
func NumericFromFloat(val float64) pgtype.Numeric {
	s := fmt.Sprintf("%.2f", val)
	var num pgtype.Numeric
	_ = num.Scan(s)
	return num
}

// FloatFromNumeric converts pgtype.Numeric to float64
func FloatFromNumeric(num pgtype.Numeric) float64 {
	if !num.Valid {
		return 0
	}
	f, _ := num.Float64Value()
	if f.Valid {
		return f.Float64
	}
	if num.Int != nil {
		bi := new(big.Float).SetInt(num.Int)
		exp := mathPow10(int(num.Exp))
		res := new(big.Float).Quo(bi, big.NewFloat(exp))
		val, _ := res.Float64()
		return val
	}
	return 0
}

func mathPow10(exp int) float64 {
	res := 1.0
	if exp < 0 {
		for i := 0; i < -exp; i++ {
			res *= 10
		}
	}
	return res
}

