package utils

import (
	g "balabanovds/go-stats/globs"
	"testing"
)

func TestSplitSlice(t *testing.T) {
	els := []g.Element{
		{IP: "1"},
		{IP: "2"},
		{IP: "3"},
	}
	divided := SplitSlice(els, 2)

	if len(divided) != 2 {
		t.Errorf("Divided len expected 2 but got %d", len(divided))
	}

	if len(divided[0]) != 2 {
		t.Errorf("Divided[0] len expected 2 but got %d", len(divided[0]))
	}

	if len(divided[1]) != 1 {
		t.Errorf("Divided[1] len expected 1 but got %d", len(divided[1]))
	}
}

func TestFilter(t *testing.T) {
	els := initEls()
	mpr := GetMPRElements(els)
	if len(mpr) != 1 {
		t.Errorf("MPR len want %v, got %v, full %+v", 1, len(mpr), mpr)
	}

	mpre := GetMPRAndMPReElements(els)
	if len(mpre) != 2 {
		t.Errorf("MPRe len want %v, got %v, full %+v", 2, len(mpr), mpre)
	}

	men := GetMENElements(els)
	if len(men) != 1 {
		t.Errorf("MEN len want %v, got %v, full %+v", 1, len(men), men)
	}

}

func TestFindDuplicates(t *testing.T) {
	els := initEls()
	els = append(els, els[0])

	dupl := FindDuplicates(els)

	if len(dupl) != 1 {
		t.Errorf("Dupl want %v, got %v", 1, len(dupl))
	}
}

func TestToServerMap(t *testing.T) {
	els := initEls()
	m := toServersMap(els)
	if len(m["1"]) != 2 {
		t.Errorf("Want %d, got %d\nfull: %+v", 2, len(m["1"]), m)
	}
}

func initEls() []g.Element {
	return []g.Element{
		{IP: "1", ElementType: 18, Server: "1"},
		{IP: "2", ElementType: 1, Server: "1"},
		{IP: "3", ElementType: 35, Server: "2"},
	}
}
