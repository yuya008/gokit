package dep

import "testing"

func TestDepScanner_Scan(t *testing.T) {
	scanner := &DepScanner{
		RootDir: "/Users/yuya/testimport",
		Deep: true,
		HasRoot: true,
	}
	p, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	for _, pk := range p {
		t.Log(pk)
	}

}
