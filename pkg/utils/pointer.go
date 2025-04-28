package utils

func PtrString(s string) *string { return &s }
func PtrInt(i int) *int           { return &i }
func PtrBool(b bool) *bool        { return &b }
func PtrFloat64(f float64) *float64 { return &f }
