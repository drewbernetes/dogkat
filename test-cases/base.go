package test_cases

func intPtr(number int32) *int32 {
	var intPtr int32
	intPtr = number
	return &intPtr
}
