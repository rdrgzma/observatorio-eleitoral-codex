package historical

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const validCSV = "election_year,round,office,uf,municipality_code,zone,section,electorate,turnout,valid_votes\n2024,1,PREFEITO,ZZ,00001,0001,0001,100,80,70\n"

func fixture(data string) Manifest {
	digest := sha256.Sum256([]byte(data))
	instant := time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)
	return Manifest{SourceURL: "fixture://section-totals", RetrievedAt: instant, AvailableAt: instant, AvailabilityEvidence: "synthetic fixture created for tests; not historical evidence", Cutoff: instant, SHA256: hex.EncodeToString(digest[:]), DatasetVersion: "synthetic/1", SchemaVersion: SchemaVersion, ParserVersion: ParserVersion, License: "TEST_FIXTURE_ONLY", Classification: "HYPOTHESIS", DataKind: "SYNTHETIC"}
}

func TestImport(t *testing.T) {
	tests := []struct{ name, from, to, reason string }{
		{"valid", "", "", ""},
		{"zero_counts", ",100,80,70", ",0,0,0", ""},
		{"boundary_counts", ",100,80,70", ",100,100,100", ""},
		{"negative_votes", ",100,80,70", ",100,80,-1", "electoral_invariant"},
		{"negative_turnout", ",100,80,70", ",100,-1,0", "electoral_invariant"},
		{"negative_electorate", ",100,80,70", ",-1,0,0", "electoral_invariant"},
		{"votes_exceed_turnout", ",100,80,70", ",100,80,81", "electoral_invariant"},
		{"turnout_exceeds_electorate", ",100,80,70", ",100,101,70", "electoral_invariant"},
		{"integer_overflow", ",100,80,70", ",9223372036854775808,80,70", "invalid_integer"},
		{"fractional_vote", ",100,80,70", ",100,80,1.5", "invalid_integer"},
		{"empty_value", ",100,80,70", ",100,,70", "invalid_integer"},
		{"wrong_office", ",PREFEITO,ZZ", ",SENADOR,ZZ", "unsupported_election"},
		{"wrong_round", "2024,1,PREFEITO", "2024,3,PREFEITO", "unsupported_election"},
		{"invalid_identity", ",0001,0001,100", ",section,0001,100", "invalid_identity"},
		{"wrong_field_count", ",100,80,70", ",100,80", "field_count"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := validCSV
			if test.from != "" {
				data = strings.Replace(data, test.from, test.to, 1)
			}
			m := fixture(data)
			records, report, err := Import([]byte(data), m, m.Cutoff)
			if err != nil {
				t.Fatal(err)
			}
			if test.reason != "" {
				if report.QualityPassed || len(records) != 0 || report.RejectedRows != 1 || report.Rejections[0].Reason != test.reason {
					t.Fatalf("unexpected rejection: %+v", report)
				}
			} else {
				if !report.QualityPassed || len(records) != 1 || report.AcceptedRows != 1 {
					t.Fatalf("unexpected result %+v", report)
				}
				if records[0].Section != "0001" || records[0].Classification != "HYPOTHESIS" || records[0].DataKind != "SYNTHETIC" {
					t.Fatal("identity spelling or synthetic label lost")
				}
			}
		})
	}
}

func TestDuplicateKeysAndElectionScope(t *testing.T) {
	row := strings.Split(validCSV, "\n")[1] + "\n"
	for _, spelling := range []string{row, strings.Replace(row, "2024,1,", "02024,01,", 1)} {
		data := validCSV + spelling
		m := fixture(data)
		rows, report, err := Import([]byte(data), m, m.Cutoff)
		if err != nil || rows != nil || report.RejectedRows != 1 || report.Rejections[0].Reason != "duplicate_key" {
			t.Fatalf("duplicate not blocked: %+v %v", report, err)
		}
	}
	data := validCSV + strings.Replace(row, "2024,1,", "2020,1,", 1)
	m := fixture(data)
	rows, report, err := Import([]byte(data), m, m.Cutoff)
	if err != nil || len(rows) != 2 || !report.QualityPassed || report.HistoricalMatching != "NOT_PERFORMED" {
		t.Fatalf("different elections must not imply historical matches: %+v %v", report, err)
	}
}

func TestManifestIntegrityAndTemporalBoundaries(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Manifest)
	}{
		{"hash", func(m *Manifest) { m.SHA256 = strings.Repeat("0", 64) }},
		{"missing_hash", func(m *Manifest) { m.SHA256 = "" }},
		{"unknown_schema", func(m *Manifest) { m.SchemaVersion = "future" }},
		{"unknown_parser", func(m *Manifest) { m.ParserVersion = "future" }},
		{"missing_license", func(m *Manifest) { m.License = "" }},
		{"missing_evidence", func(m *Manifest) { m.AvailabilityEvidence = "" }},
		{"missing_version", func(m *Manifest) { m.DatasetVersion = "" }},
		{"real_data", func(m *Manifest) { m.DataKind = "OFFICIAL"; m.Classification = "DATA_FACT" }},
		{"network_source", func(m *Manifest) { m.SourceURL = "https://example.com/input.csv" }},
		{"after_cutoff", func(m *Manifest) { m.AvailableAt = m.Cutoff.Add(time.Nanosecond); m.RetrievedAt = m.AvailableAt }},
		{"after_retrieval", func(m *Manifest) { m.RetrievedAt = m.AvailableAt.Add(-time.Nanosecond) }},
		{"missing_availability", func(m *Manifest) { m.AvailableAt = time.Time{} }},
		{"missing_retrieval", func(m *Manifest) { m.RetrievedAt = time.Time{} }},
		{"missing_cutoff", func(m *Manifest) { m.Cutoff = time.Time{} }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := fixture(validCSV)
			cutoff := m.Cutoff
			test.mutate(&m)
			records, _, err := Import([]byte(validCSV), m, cutoff)
			if err == nil || records != nil {
				t.Fatal("invalid manifest released data")
			}
		})
	}
	m := fixture(validCSV)
	m.RetrievedAt = m.Cutoff.Add(time.Hour)
	if _, _, err := Import([]byte(validCSV), m, m.Cutoff); err != nil {
		t.Fatalf("documented availability may precede later retrieval: %v", err)
	}
	if _, _, err := Import([]byte(validCSV), m, m.Cutoff.Add(time.Second)); err == nil {
		t.Fatal("different execution cutoff accepted")
	}
	if _, _, err := Import([]byte(validCSV), m, time.Time{}); err == nil {
		t.Fatal("missing execution cutoff accepted")
	}
}

func TestMalformedAndEmptyInputs(t *testing.T) {
	for _, data := range []string{"", strings.Split(validCSV, "\n")[0] + "\n", strings.Replace(validCSV, "valid_votes", "votes", 1), validCSV + "\"unterminated\n", strings.Repeat("x", MaxBytes+1)} {
		m := fixture(data)
		if rows, _, err := Import([]byte(data), m, m.Cutoff); err == nil || rows != nil {
			t.Fatal("invalid input accepted")
		}
	}
	for _, data := range []string{`{"unexpected":true}`, `{} {}`, `{} trailing`} {
		if _, err := ReadManifest(strings.NewReader(data)); err == nil {
			t.Fatal("ambiguous manifest accepted")
		}
	}
}

func TestRunArtifactsAndNoOverwrite(t *testing.T) {
	for _, reject := range []bool{false, true} {
		t.Run(map[bool]string{false: "valid", true: "rejected"}[reject], func(t *testing.T) {
			root := t.TempDir()
			input := filepath.Join(root, "input.csv")
			manifest := filepath.Join(root, "manifest.json")
			output := filepath.Join(root, "run")
			data := validCSV
			if reject {
				data = strings.Replace(data, ",100,80,70", ",100,80,81", 1)
			}
			m := fixture(data)
			payload, _ := json.Marshal(m)
			if err := os.WriteFile(input, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(manifest, payload, 0600); err != nil {
				t.Fatal(err)
			}
			report, err := Run(input, manifest, output, m.Cutoff)
			if err != nil {
				t.Fatal(err)
			}
			raw := filepath.Join(output, "raw", m.SHA256+".csv")
			archived, err := os.ReadFile(raw)
			if err != nil || string(archived) != data {
				t.Fatal("raw bytes changed")
			}
			original, _ := os.ReadFile(input)
			if string(original) != data {
				t.Fatal("input overwritten")
			}
			info, _ := os.Stat(raw)
			if info.Mode().Perm()&0222 != 0 {
				t.Fatal("archive must be read-only")
			}
			normalized, err := os.ReadFile(filepath.Join(output, "normalized.jsonl"))
			if reject {
				if !os.IsNotExist(err) || report.QualityPassed {
					t.Fatal("rejected dataset released normalized artifact")
				}
			} else {
				if err != nil || !strings.Contains(string(normalized), `"classification":"HYPOTHESIS"`) {
					t.Fatal("synthetic labels missing")
				}
				var record Record
				if err = json.Unmarshal(bytesTrim(normalized), &record); err != nil {
					t.Fatal(err)
				}
			}
			for _, name := range []string{"manifest.json", "quality.json", "artifacts.json", "complete.json"} {
				if _, err := os.Stat(filepath.Join(output, name)); err != nil {
					t.Fatal(err)
				}
			}
			before, _ := os.ReadFile(filepath.Join(output, "quality.json"))
			if _, err = Run(input, manifest, output, m.Cutoff); err == nil {
				t.Fatal("existing run overwritten")
			}
			after, _ := os.ReadFile(filepath.Join(output, "quality.json"))
			if string(before) != string(after) {
				t.Fatal("existing artifacts changed")
			}
		})
	}
}

func bytesTrim(data []byte) []byte { return []byte(strings.TrimSpace(string(data))) }

func TestBadHashProducesNoArtifacts(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.csv")
	manifest := filepath.Join(root, "manifest.json")
	output := filepath.Join(root, "run")
	m := fixture(validCSV)
	m.SHA256 = strings.Repeat("0", 64)
	payload, _ := json.Marshal(m)
	os.WriteFile(input, []byte(validCSV), 0600)
	os.WriteFile(manifest, payload, 0600)
	if _, err := Run(input, manifest, output, m.Cutoff); err == nil {
		t.Fatal("bad hash accepted")
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatal("bad input created output artifacts")
	}
}
