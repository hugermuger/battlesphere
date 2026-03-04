package dbutils

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func ToNullInt32(i *int) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}

func ToNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

func ToNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func ToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: *u, Valid: true}
}

func ToNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func StringToNullFloat64(s *string) (sql.NullFloat64, error) {
	if s == nil || *s == "" {
		return sql.NullFloat64{Valid: false}, nil
	}
	f, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		fmt.Println("Error parsing string:", err)
		return sql.NullFloat64{Valid: false}, err
	}
	return sql.NullFloat64{Float64: f, Valid: true}, nil
}

func SafeSlice(s *[]string) []string {
	if s == nil {
		return []string{}
	}
	return *s
}
