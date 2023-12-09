package ptr

func String(v string) *string {
	return &v
}

func Uint64(i uint64) *uint64 {
	return &i
}

func Int(i int) *int {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}
