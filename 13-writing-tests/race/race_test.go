package race

import "testing"

func Test_getCounter(t *testing.T) {
	counter := GetCounter()
	if counter != 5000 {
		t.Error("unexpected counter:", counter)
	}
}
