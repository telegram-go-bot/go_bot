package vk

import "testing"

func TestInitArrayOfIndexes(t *testing.T) {
	const size = 5
	arr := InitArrayOfIndexes(size)
	if len(arr) != size {
		t.Errorf("size of Array failed. Expected: %d, Actual: %d", size, len(arr))
	}
}

func TestInitArrayOfIndexesIsRandomShuffled(t *testing.T) {
	const size = 10
	arr := InitArrayOfIndexes(size)
	// if array is sorted, or filled with similar vals - fail
	sortedAsc := true
	sortedDesc := true
	prevVal := 0
	for idx, item := range arr {
		if idx == 0 {
			prevVal = item
			continue
		}

		if prevVal <= item {
			sortedAsc = false
		} else if prevVal >= item {
			sortedDesc = false
		}
	}

	if sortedAsc || sortedDesc {
		t.Errorf("Error. Got not shuffled array: %v", arr)
	}
}

func TestRandomShuffle(t *testing.T) {
	const size = 10
	vals := []int{1, 2, 3, 4, 5}
	resVals := RandomShuffle(vals)

	if len(resVals) != len(vals) {
		t.Errorf("RandomShuffle returned array with invalid size. Expected: %d, Actial: %d", len(vals), len(resVals))
	}

	same := true
	for idx, val := range resVals {
		if val != vals[idx] {
			same = false
		}
	}

	if same {
		t.Errorf("After shuffling array stayed unchanged")
	}
}
