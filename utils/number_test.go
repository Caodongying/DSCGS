package utils

import (
	"testing"
)

func TestDecimalTo62(t *testing.T) { // 测试用例用Test开头
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "a"},
		{61, "Z"},
		{62, "10"},
		{12345, "3d7"},
		{3844, "100"},
	}

	for _, test := range tests {
		result := DecimalTo62(test.input)
		if result != test.expected {
			t.Errorf("DecimalTo62(%d) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
