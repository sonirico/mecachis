package container

import (
	"testing"
)

type testObject struct {
	value interface{}
}

func newSet(elements []object) *Set {
	set := NewSet()
	for _, element := range elements {
		set.Add(element.value)
	}
	return set
}

func testNodeEquals(t *testing.T, expected testObject, actual object) bool {
	if expected.value != actual.value {
		t.Errorf("element mismatch. want %v, have %v", expected.value, actual.value)
		return false
	}
	return true
}

func testSetState(t *testing.T, s *Set, expected []testObject) bool {
	t.Helper()

	actualLength := len(s.Elements())
	expectedLength := len(expected)
	if expectedLength != actualLength {
		t.Errorf("unexpected set length. want %d, have %d.", expectedLength, actualLength)
		return false
	}

	for position, actualElement := range s.Elements() {
		expectedElement := expected[position]
		testNodeEquals(t, expectedElement, actualElement)
	}
	return true
}

func TestSetAddUpdatesDuplicated_Add(t *testing.T) {
	payload := []object{
		{value: "color"},
		{value: "size"},
		{value: "provider"},
		{value: "theme"},
	}
	set := newSet(payload)
	expected := []testObject{
		{value: "theme"},
		{value: "provider"},
		{value: "size"},
		{value: "color"},
	}
	testSetState(t, set, expected)
}

func TestSet_Remove_Middle(t *testing.T) {
	payload := []object{
		{value: "color"},
		{value: "size"},
		{value: "provider"},
		{value: "theme"},
	}
	set := newSet(payload)
	set.Remove("size")
	expected := []testObject{
		{value: "theme"},
		{value: "provider"},
		{value: "color"},
	}
	testSetState(t, set, expected)
}

func TestSet_Remove_Head(t *testing.T) {
	payload := []object{
		{value: "color"},
		{value: "size"},
		{value: "provider"},
		{value: "theme"},
		{value: "size"},
	}
	set := newSet(payload)
	set.Remove("color")
	expected := []testObject{
		{value: "theme"},
		{value: "provider"},
		{value: "size"},
	}
	testSetState(t, set, expected)
}

func TestSet_Remove_Foot(t *testing.T) {
	payload := []object{
		{value: "color"},
		{value: 1},
		{value: false},
	}
	set := newSet(payload)
	set.Remove(false)
	expected := []testObject{
		{value: 1},
		{value: "color"},
	}
	testSetState(t, set, expected)
}

func TestSet_Remove_All(t *testing.T) {
	payload := []object{
		{value: "color"},
		{value: 1},
		{value: false},
	}
	set := newSet(payload)
	for _, payloadItem := range payload {
		set.Remove(payloadItem.value)
	}
	var expected []testObject
	testSetState(t, set, expected)
}

func TestSet_PopLast(t *testing.T) {
	payload := []object{
		{value: "color"},
		{value: 1},
		{value: false},
	}
	set := newSet(payload)
	expected := []interface{}{
		"color",
		1,
		false,
	}
	for _, expectedObj := range expected {
		actualObject := set.PopFirst()
		if actualObject == nil {
			t.Errorf("unexpected nil pointer. want %v", expectedObj)
		} else {
			if expectedObj != actualObject {
				t.Errorf("value missmatch. want %v, have %v", expectedObj, actualObject)
			}
		}
	}
}
