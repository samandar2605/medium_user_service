package utils

import "database/sql"

func NullString(s string) (ns sql.NullString) {
	if s != "" {
		ns.String = s
		ns.Valid = true
	}
	return ns
}

func NullFloat64(s float64) (nf sql.NullFloat64) {
	if s != 0 {
		nf.Float64 = s
		nf.Valid = true
	}
	return nf
}

func NullInt32(v int32) (ni sql.NullInt32) {
	if v != 0 {
		ni.Int32 = v
		ni.Valid = true
	}
	return ni
}

func FormatNullTime(nt sql.NullTime, format string) string {
	if nt.Valid {
		return nt.Time.Format(format)
	}
	return ""
}
