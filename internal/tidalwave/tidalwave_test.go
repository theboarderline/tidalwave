package tidalwave

import "testing"

func TestNameFormatter(t *testing.T) {
	for _, v := range []string{"my-controlplane", "my"} {
		name := NameFormatter(v)
		if name != "my" {
			t.Errorf("NameFormatter(\"%s\") = %s want my", v, name)
		}
	}
}
