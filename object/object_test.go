package object

import "testing"

func TestStringHash(t *testing.T) {
	s1 := &String{Value: "Hello World"}
	s2 := &String{Value: "Hello World"}
	s3 := &String{Value: "My name is PanJinhao"}
	s4 := &String{Value: "My name is PanJinhao"}

	if s1.Hash() != s2.Hash() {
		t.Errorf("strings with same content have different hash keys")
	}

	if s3.Hash() != s4.Hash() {
		t.Errorf("strings with same content have different hash keys")
	}

	if s1.Hash() == s3.Hash() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHash(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.Hash() != true2.Hash() {
		t.Errorf("trues do not have same hash key")
	}

	if false1.Hash() != false2.Hash() {
		t.Errorf("falses do not have same hash key")
	}

	if true1.Hash() == false1.Hash() {
		t.Errorf("true has same hash key as false")
	}
}

func TestIntegerHash(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two1 := &Integer{Value: 2}
	two2 := &Integer{Value: 2}

	if one1.Hash() != one2.Hash() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if two1.Hash() != two2.Hash() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if one1.Hash() == two1.Hash() {
		t.Errorf("integers with twoerent content have same hash keys")
	}
}
