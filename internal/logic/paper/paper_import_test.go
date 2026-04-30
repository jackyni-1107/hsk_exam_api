package paper

import (
	"testing"

	"exam/internal/consts"
	examentity "exam/internal/model/entity/exam"
)

func TestNormalizeExamImportConflictMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "default empty", input: "", want: consts.ExamImportConflictFail},
		{name: "trimmed overwrite", input: " overwrite ", want: consts.ExamImportConflictOverwrite},
		{name: "case insensitive", input: "NEW", want: consts.ExamImportConflictNew},
		{name: "legacy alias", input: "new_copy", want: consts.ExamImportConflictNew},
		{name: "invalid", input: "replace", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := normalizeExamImportConflictMode(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("normalizeExamImportConflictMode(%q)=%q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIndexJSONURLFromMockResourceURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "zip resource",
			input: "https://cdn.example.com/hsk/level4/paper01.zip?token=1",
			want:  "https://cdn.example.com/hsk/level4/paper01/index.json",
		},
		{
			name:  "folder resource",
			input: "https://cdn.example.com/hsk/level4/paper01/",
			want:  "https://cdn.example.com/hsk/level4/paper01/index.json",
		},
		{
			name:  "existing index",
			input: "https://cdn.example.com/hsk/level4/paper01/index.json#frag",
			want:  "https://cdn.example.com/hsk/level4/paper01/index.json",
		},
		{
			name:    "empty",
			input:   " ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := indexJSONURLFromMockResourceURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("indexJSONURLFromMockResourceURL(%q)=%q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseIndexURL(t *testing.T) {
	t.Parallel()

	baseURL, level, paperID, err := parseIndexURL("https://cdn.example.com/hsk/level4/paper01/index.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if baseURL != "https://cdn.example.com/hsk/level4/paper01/" {
		t.Fatalf("unexpected baseURL: %q", baseURL)
	}
	if level != "level4" {
		t.Fatalf("unexpected level: %q", level)
	}
	if paperID != "paper01" {
		t.Fatalf("unexpected paperID: %q", paperID)
	}

	if _, _, _, err := parseIndexURL("https://cdn.example.com/hsk/index.json"); err == nil {
		t.Fatal("expected error for short path")
	}
	if _, _, _, err := parseIndexURL("https://cdn.example.com/hsk/level4/paper01/topic.json"); err == nil {
		t.Fatal("expected error for non-index path")
	}
}

func TestResolveImportAudioHlsPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		requestPrefix string
		plan          importConflictPlan
		want          string
	}{
		{
			name:          "create uses request prefix",
			requestPrefix: "/hls/new/",
			plan:          importConflictPlan{},
			want:          "hls/new",
		},
		{
			name:          "overwrite keeps existing when request empty",
			requestPrefix: "",
			plan: importConflictPlan{
				overwritePaperID: 1,
				existing: []examentity.ExamPaper{
					{Id: 1, AudioHlsPrefix: "/hls/existing/"},
				},
			},
			want: "hls/existing",
		},
		{
			name:          "overwrite uses explicit request prefix",
			requestPrefix: "hls/override",
			plan: importConflictPlan{
				overwritePaperID: 1,
				existing: []examentity.ExamPaper{
					{Id: 1, AudioHlsPrefix: "/hls/existing/"},
				},
			},
			want: "hls/override",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := resolveImportAudioHlsPrefix(tt.requestPrefix, tt.plan)
			if got != tt.want {
				t.Fatalf("resolveImportAudioHlsPrefix(%q)=%q, want %q", tt.requestPrefix, got, tt.want)
			}
		})
	}
}
