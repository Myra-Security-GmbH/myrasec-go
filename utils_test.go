package myrasec

import (
	"testing"
)

func TestIntInSlice(t *testing.T) {
	var result bool

	result = intInSlice(1, []int{1})
	if !result {
		t.Errorf("Expected to find the passed needle in the haystack")
	}

	result = intInSlice(1, []int{5})
	if result {
		t.Errorf("Expected not to find the passed needle in the haystack")
	}

	result = intInSlice(1, []int{1})
	if !result {
		t.Errorf("Expected to find the passed needle in the haystack")
	}

	result = intInSlice(200, []int{1, 2, 3, 4, 200, 201, 204, 403, 404, 500})
	if !result {
		t.Errorf("Expected to find the passed needle in the haystack")
	}

	result = intInSlice(200, []int{})
	if result {
		t.Errorf("Expected not to find the passed needle in the haystack")
	}

	result = intInSlice(0, []int{})
	if result {
		t.Errorf("Expected not to find the passed needle in the haystack")
	}

	result = intInSlice(0, []int{1})
	if result {
		t.Errorf("Expected not to find the passed needle in the haystack")
	}

	result = intInSlice(0, []int{0})
	if !result {
		t.Errorf("Expected to find the passed needle in the haystack")
	}
}

func TestBuildSHA256(t *testing.T) {
	var result string

	result = BuildSHA256("")
	if result != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Errorf("Expected to get [%s] but got [%s]", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", result)
	}

	result = BuildSHA256("https://apiv2.myracloud.com/domains")
	if result != "b67371960a647a3638f99d28bfd15ba147c18196142efad47d9a2743f8a43515" {
		t.Errorf("Expected to get [%s] but got [%s]", "b67371960a647a3638f99d28bfd15ba147c18196142efad47d9a2743f8a43515", result)
	}
}
