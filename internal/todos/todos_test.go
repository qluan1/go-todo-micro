package todos

import "testing"

func TestChecksValidation(t *testing.T) {
	td := &Todo{}
	err := td.Validate()

	if  err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
