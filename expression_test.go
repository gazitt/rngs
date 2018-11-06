package main

import (
	"reflect"
	"testing"
)

func TestExpressionSplit(t *testing.T) {
	data := []struct {
		expression      string
		expect          []string
		shouldBeAnError bool
	}{
		{
			expression:      "/foo/bar",
			expect:          []string{"", "foo", "bar"},
			shouldBeAnError: false,
		},
		{
			expression:      "i/foo/bar",
			expect:          []string{"i", "foo", "bar"},
			shouldBeAnError: false,
		},
		{
			expression:      "/\\/foo/bar",
			expect:          []string{"", "\\/foo", "bar"},
			shouldBeAnError: false,
		},
		{
			expression:      "/foo/",
			expect:          []string{"", "foo", ""},
			shouldBeAnError: false,
		},
		{
			expression:      "/",
			shouldBeAnError: true,
		},
		{
			expression:      "",
			shouldBeAnError: true,
		},
	}

	for i, v := range data {
		actual, err := split(v.expression)
		if v.shouldBeAnError {
			if err == nil {
				t.Errorf("%d: should be an error: %s\n", i, v.expression)
			}
			continue
		}
		if err != nil {
			t.Errorf("%d: error: %s: %s\n", i, v.expression, err)
			continue
		}
		if !reflect.DeepEqual(actual, v.expect) {
			t.Errorf("%d:\ngot : %s, want: %s\n", i, actual, v.expect)
		}
	}
}

func TestExpressionToScript(t *testing.T) {
	e := Expression{}
	if err := e.Parse("//"); err != nil {
		t.Errorf("1: Parse error: %s\n", err)
	}
	if err := e.Parse("//$index"); err != nil {
		t.Errorf("2: Parse error: %s\n", err)
	}

	expect := "$func.join($dirname, $basename.replace(//, '').replace(//, $index));"
	actual := e.ToScript()
	if expect != actual {
		t.Errorf("\ngot : %s\nwant: %s\n", actual, expect)
	}
}
