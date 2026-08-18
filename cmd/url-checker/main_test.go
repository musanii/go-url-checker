package main

import "testing"

func TestRunRequiresURLs(t *testing.T){
	err := run([]string{})

	if err ==nil{
		t.Fatal("expected URL argument error, got nil")
	}

	expected := "at least one URL is required"

	if err.Error() != expected{
		t.Fatalf(
			"expected error %q, got %q",
			expected,
			err.Error(),
		)

	}
}