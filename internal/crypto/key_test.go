package crypto

import (
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tt := []struct {
		s int
		l int
	}{
		{32, 64},
		{10, 20},
		{12, 24},
		{128, 256},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d", tc.s), func(t *testing.T) {
			k, err := GenerateKey(tc.s)
			if err != nil {
				t.Errorf("incorrect error, got %v, want nil", err)
				t.FailNow()
			}
			if len(k) != tc.l {
				t.Errorf("incorrect length, got %d, want %d", len(k), tc.l)
			}
		})
	}

}
