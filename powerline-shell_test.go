package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_addCwd_root(t *testing.T) {
	segments := [][]string{}

	dir := "/"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(segments, []string{"250", "237", "/"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_one(t *testing.T) {
	segments := [][]string{}

	dir := "/Go"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"250", "237", "Go"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_two(t *testing.T) {
	segments := [][]string{}

	dir := "/Go/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"250", "237", "Go", ">", "244"},
		[]string{"250", "237", "src"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_root_three(t *testing.T) {
	segments := [][]string{}

	dir := "/Go/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"250", "237", "Go", ">", "244"},
		[]string{"250", "237", "...", ">", "244"},
		[]string{"250", "237", "github.com"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home(t *testing.T) {
	segments := [][]string{}

	dir := "~"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(segments, []string{"015", "031", "~"})

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_one(t *testing.T) {
	segments := [][]string{}

	dir := "~/Go"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "Go"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_two(t *testing.T) {
	segments := [][]string{}

	dir := "~/Go/src"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "Go", ">", "244"},
		[]string{"250", "237", "src"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_three(t *testing.T) {
	segments := [][]string{}

	dir := "~/Go/src/github.com"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "Go", ">", "244"},
		[]string{"250", "237", "...", ">", "244"},
		[]string{"250", "237", "github.com"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}

func Test_addCwd_home_five(t *testing.T) {
	segments := [][]string{}

	dir := "~/Go/src/github.com/sivel/powerline-shell-go"
	parts := strings.Split(dir, "/")

	rootSegments := addCwd(parts, "...", ">")
	want := append(
		segments,
		[]string{"015", "031", "~"},
		[]string{"250", "237", "Go", ">", "244"},
		[]string{"250", "237", "...", ">", "244"},
		[]string{"250", "237", "powerline-shell-go"},
	)

	if !reflect.DeepEqual(rootSegments, want) {
		t.Errorf("addCwd returned %+v, not %+v", rootSegments, want)
	}
}
