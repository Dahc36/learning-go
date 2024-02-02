package person_test

import (
	"testing"

	// You can use this package created by Google for nicer messages and comparisons
	"github.com/google/go-cmp/cmp"

	"github.com/dahc36/learning-go/13-writing-tests/person"
)

func TestCreatePerson(t *testing.T) {
	expected := person.Person{
		Name: "Dennis",
		Age:  87,
	}
	// You can create your own comparer function
	comparer := cmp.Comparer(func(x, y person.Person) bool {
		return x.Name == y.Name && x.Age == y.Age
	})
	result := person.CreatePerson("Dennis", 87)
	if diff := cmp.Diff(expected, result, comparer); diff != "" {
		t.Error(diff)
	}
}
