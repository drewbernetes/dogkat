package hack

import "testing"

// TestIntPtr tests the IntPtr function to ensure a pointer with a value is returned.
func TestIntPtr(t *testing.T) {
	ptr := IntPtr(int32(10))
	if *ptr != 10 {
		t.Errorf("The Int Pointer function failed to create a pointer to the number.")
	}
}
