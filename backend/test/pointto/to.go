package pointto

func Any(a interface{}) interface{} {
	return &a
}

func String(s string) *string {
	return &s
}
