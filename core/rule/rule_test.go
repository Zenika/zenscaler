package rule

import "testing"

func TestConfigDecode(t *testing.T) {
	decodeHelper(t, 10.0, "< 25.23", true)
	decodeHelper(t, 0.20, "< .23", true)
	decodeHelper(t, 26, "< 25", false)
	decodeHelper(t, 55, "> 56.", false)
}

func decodeHelper(t *testing.T, f float64, s string, expected bool) {
	fun, err := Decode(s)
	if err != nil {
		t.Error(err)
	}
	if fun(f) != expected {
		t.Fail()
	}
}
