package mongoutils_test

import (
	"testing"

	"github.com/gomig/mongoutils"
)

type typeA struct {
	Name   string
	Family string
}

func (a typeA) toMap() map[string]any {
	return map[string]any{"Name": a.Name, "Family": a.Family}
}

type typeB struct {
	A       string
	Persons []typeA
	Var     int
	Ptr     *int
	Slice   []string
}

func (b typeB) toMap() map[string]any {
	res := map[string]any{}
	res["A"] = b.A
	res["Var"] = b.Var
	res["Ptr"] = b.Ptr
	res["Slice"] = b.Slice
	persons := []any{}
	for _, p := range b.Persons {
		persons = append(persons, p.toMap())
	}
	res["Persons"] = persons
	return res
}

func TestChecksum(t *testing.T) {
	i := 2
	demo := typeB{
		Var:     i,
		Ptr:     &i,
		Slice:   []string{"A", "b", "c"},
		Persons: []typeA{{"John", "Doe"}, {"Jack", "Ma"}},
	}

	checksum := mongoutils.NewChecksum(demo.toMap())
	if n := checksum.Normalize(); n != "A:|Persons.E0.Family:Doe|Persons.E0.Name:John|Persons.E1.Family:Ma|Persons.E1.Name:Jack|Ptr:2|Slice.E0:A|Slice.E1:b|Slice.E2:c|Var:2" {
		t.Log(n)
		t.Fatal("failed expected normalize")
	}

	if n := checksum.MD5(); n != "1c4755dc74daa55c60657667a50a00fb" {
		t.Log(n)
		t.Fatal("failed expected md5")
	}
}
