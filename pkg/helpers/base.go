package helpers

// IntPtr returns a pointer to an int that is provided.
func IntPtr(number int32) *int32 {
	var intPtr int32
	intPtr = number
	return &intPtr
}
