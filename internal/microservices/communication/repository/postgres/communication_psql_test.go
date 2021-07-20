package postgres

import "testing"

func TestCommunicationPSQL_CreateCommunication(t *testing.T) {
	testCases := map[string]struct {
		testValue bool
		want      bool
	}{
		"Expected true":  {true, true},
		"Expected false": {false, false},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if testCase.testValue != testCase.want {
				t.Fatal("Boolean variables differ")
			}
		})
	}
}
