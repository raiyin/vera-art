package dto

import (
	"reflect"
	"testing"
)

func TestSplitImages(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty string", "", nil},
		{"single image", "img1.jpg", []string{"img1.jpg"}},
		{"multiple images", "img1.jpg;img2.jpg;img3.jpg", []string{"img1.jpg", "img2.jpg", "img3.jpg"}},
		{"with trailing semicolon", "a.jpg;b.jpg;", []string{"a.jpg", "b.jpg", ""}},
		{"only semicolon", ";", []string{"", ""}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SplitImages(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("SplitImages(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestJoinImages(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"nil slice", nil, ""},
		{"empty slice", []string{}, ""},
		{"single image", []string{"img1.jpg"}, "img1.jpg"},
		{"multiple images", []string{"img1.jpg", "img2.jpg"}, "img1.jpg;img2.jpg"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := JoinImages(tc.input)
			if got != tc.want {
				t.Errorf("JoinImages(%v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestSplitJoinRoundTrip(t *testing.T) {
	original := "a.jpg;b.jpg;c.jpg"
	split := SplitImages(original)
	joined := JoinImages(split)
	if joined != original {
		t.Errorf("round trip failed: %q -> %v -> %q", original, split, joined)
	}
}

func TestJoinSplitRoundTrip(t *testing.T) {
	original := []string{"x.jpg", "y.jpg", "z.jpg"}
	joined := JoinImages(original)
	split := SplitImages(joined)
	if !reflect.DeepEqual(split, original) {
		t.Errorf("round trip failed: %v -> %q -> %v", original, joined, split)
	}
}

func TestSplitEmptyRoundTrip(t *testing.T) {
	if got := SplitImages(""); got != nil {
		t.Errorf("SplitImages(\"\") should return nil, got %v", got)
	}
}

func TestJoinEmptyRoundTrip(t *testing.T) {
	if got := JoinImages(nil); got != "" {
		t.Errorf("JoinImages(nil) should return \"\", got %q", got)
	}
	if got := JoinImages([]string{}); got != "" {
		t.Errorf("JoinImages([]string{}) should return \"\", got %q", got)
	}
}

func TestWorkResponseFields(t *testing.T) {
	resp := WorkResponse{
		ID:       1,
		StrID:    "test",
		Width:    100,
		Height:   200,
		Year:     2024,
		NameRu:   "Работа",
		NameEn:   "Work",
		BaseID:   3,
		DescrRu:  "Описание",
		DescrEn:  "Description",
		WorkPath: "test/",
		Images:   []string{"a.jpg"},
	}

	if resp.ID != 1 || resp.StrID != "test" || resp.Width != 100 || resp.Height != 200 {
		t.Error("WorkResponse basic fields mismatch")
	}
	if resp.Year != 2024 || resp.NameRu != "Работа" || resp.NameEn != "Work" {
		t.Error("WorkResponse name fields mismatch")
	}
	if resp.BaseID != 3 || resp.DescrRu != "Описание" || resp.DescrEn != "Description" {
		t.Error("WorkResponse description fields mismatch")
	}
	if resp.WorkPath != "test/" {
		t.Errorf("WorkPath = %q, want %q", resp.WorkPath, "test/")
	}
	if len(resp.Images) != 1 || resp.Images[0] != "a.jpg" {
		t.Errorf("Images = %v, want [a.jpg]", resp.Images)
	}
}

func TestCreateWorkRequestRequiredFields(t *testing.T) {
	req := CreateWorkRequest{
		NameRu: "Работа",
		NameEn: "Work",
	}
	if req.NameRu != "Работа" {
		t.Errorf("NameRu = %q, want %q", req.NameRu, "Работа")
	}
	if req.NameEn != "Work" {
		t.Errorf("NameEn = %q, want %q", req.NameEn, "Work")
	}
}

func TestUpdateWorkRequest(t *testing.T) {
	req := UpdateWorkRequest{
		Width:  150,
		Height: 250,
		Year:   2025,
	}
	if req.Width != 150 || req.Height != 250 || req.Year != 2025 {
		t.Error("UpdateWorkRequest basic fields mismatch")
	}
}
