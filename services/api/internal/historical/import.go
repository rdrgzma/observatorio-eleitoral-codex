// Package historical imports the project's synthetic canonical contract.
// It does not parse TSE files or certify the availability of official data.
package historical

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const ParserVersion = "canonical-csv/0.1.0"
const SchemaVersion = "section-totals/0.1.0"
const MaxBytes = 8 * 1024 * 1024

var header = []string{"election_year", "round", "office", "uf", "municipality_code", "zone", "section", "electorate", "turnout", "valid_votes"}
var digits = regexp.MustCompile(`^[0-9]+$`)
var stateCode = regexp.MustCompile(`^[A-Z]{2}$`)
var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// Manifest describes evidence supplied by the caller; timestamps are not inferred
// from the election year or the time of download.
type Manifest struct {
	SourceURL            string    `json:"source_url"`
	RetrievedAt          time.Time `json:"retrieved_at"`
	AvailableAt          time.Time `json:"available_at"`
	AvailabilityEvidence string    `json:"availability_evidence"`
	Cutoff               time.Time `json:"cutoff"`
	SHA256               string    `json:"sha256"`
	DatasetVersion       string    `json:"dataset_version"`
	SchemaVersion        string    `json:"schema_version"`
	ParserVersion        string    `json:"parser_version"`
	License              string    `json:"license"`
	Classification       string    `json:"classification"`
	DataKind             string    `json:"data_kind"`
}

func ReadManifest(r io.Reader) (Manifest, error) {
	var m Manifest
	decoder := json.NewDecoder(io.LimitReader(r, MaxBytes+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&m); err != nil {
		return m, fmt.Errorf("manifest: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return m, errors.New("manifest must contain exactly one JSON object")
	}
	return m, nil
}

func (m Manifest) Validate(cutoff time.Time) error {
	if m.DataKind != "SYNTHETIC" || m.Classification != "HYPOTHESIS" {
		return errors.New("only explicitly labelled SYNTHETIC/HYPOTHESIS input is supported; official adapter is pending")
	}
	if !strings.HasPrefix(m.SourceURL, "fixture://") || len(m.SourceURL) <= len("fixture://") {
		return errors.New("synthetic source_url must identify a fixture:// resource")
	}
	if !digestPattern.MatchString(m.SHA256) || m.DatasetVersion == "" || m.License == "" || m.AvailabilityEvidence == "" {
		return errors.New("missing provenance, version, license or invalid SHA-256")
	}
	if m.SchemaVersion != SchemaVersion || m.ParserVersion != ParserVersion {
		return errors.New("unsupported schema or parser version")
	}
	if cutoff.IsZero() || m.Cutoff.IsZero() || !m.Cutoff.Equal(cutoff) {
		return errors.New("manifest cutoff must match the explicit execution cutoff")
	}
	if m.AvailableAt.IsZero() || m.RetrievedAt.IsZero() || m.AvailableAt.After(cutoff) || m.AvailableAt.After(m.RetrievedAt) {
		return errors.New("invalid availability: require available_at <= cutoff and available_at <= retrieved_at")
	}
	return nil
}

// Record is one section/cargo/turno summary, not a candidate-vote row. v0.1 is
// limited to prefeito; its valid_votes <= turnout invariant is not applicable
// without review to contests where a voter can cast more than one valid vote.
type Record struct {
	Classification   string `json:"classification"`
	DataKind         string `json:"data_kind"`
	ElectionYear     int    `json:"election_year"`
	Round            int    `json:"round"`
	Office           string `json:"office"`
	UF               string `json:"uf"`
	MunicipalityCode string `json:"municipality_code"`
	Zone             string `json:"zone"`
	Section          string `json:"section"`
	Electorate       int64  `json:"electorate"`
	Turnout          int64  `json:"turnout"`
	ValidVotes       int64  `json:"valid_votes"`
}

type Rejection struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

type Report struct {
	Classification     string      `json:"classification"`
	DataKind           string      `json:"data_kind"`
	SchemaVersion      string      `json:"schema_version"`
	ParserVersion      string      `json:"parser_version"`
	Cutoff             time.Time   `json:"cutoff"`
	InputRows          int         `json:"input_rows"`
	AcceptedRows       int         `json:"accepted_rows"`
	RejectedRows       int         `json:"rejected_rows"`
	Rejections         []Rejection `json:"rejections"`
	QualityPassed      bool        `json:"quality_passed"`
	HistoricalMatching string      `json:"historical_matching"`
}

func parse(row []string) (Record, error) {
	var r Record
	if len(row) != len(header) {
		return r, errors.New("field_count")
	}
	year, e1 := strconv.Atoi(row[0])
	round, e2 := strconv.Atoi(row[1])
	electorate, e3 := strconv.ParseInt(row[7], 10, 64)
	turnout, e4 := strconv.ParseInt(row[8], 10, 64)
	votes, e5 := strconv.ParseInt(row[9], 10, 64)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil {
		return r, errors.New("invalid_integer")
	}
	r = Record{Classification: "HYPOTHESIS", DataKind: "SYNTHETIC", ElectionYear: year, Round: round, Office: row[2], UF: row[3], MunicipalityCode: row[4], Zone: row[5], Section: row[6], Electorate: electorate, Turnout: turnout, ValidVotes: votes}
	if year < 1900 || year > 9999 || (round != 1 && round != 2) || r.Office != "PREFEITO" {
		return r, errors.New("unsupported_election")
	}
	if !stateCode.MatchString(r.UF) || !digits.MatchString(r.MunicipalityCode) || !digits.MatchString(r.Zone) || !digits.MatchString(r.Section) {
		return r, errors.New("invalid_identity")
	}
	if votes < 0 || turnout < 0 || electorate < 0 || votes > turnout || turnout > electorate {
		return r, errors.New("electoral_invariant")
	}
	return r, nil
}

// Import validates the entire bounded input before releasing normalized rows.
// If any record is rejected, no normalized rows are returned.
func Import(data []byte, m Manifest, cutoff time.Time) ([]Record, Report, error) {
	report := Report{Classification: "METHOD_DOCUMENTED", DataKind: "SYNTHETIC", SchemaVersion: SchemaVersion, ParserVersion: ParserVersion, Cutoff: cutoff, Rejections: []Rejection{}, HistoricalMatching: "NOT_PERFORMED"}
	if err := m.Validate(cutoff); err != nil {
		return nil, report, err
	}
	if len(data) > MaxBytes {
		return nil, report, errors.New("input exceeds 8 MiB limit")
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != m.SHA256 {
		return nil, report, errors.New("input SHA-256 mismatch")
	}
	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1
	columns, err := reader.Read()
	if err != nil {
		return nil, report, fmt.Errorf("CSV header: %w", err)
	}
	if strings.Join(columns, "\x00") != strings.Join(header, "\x00") {
		return nil, report, errors.New("unexpected canonical CSV header (not a TSE layout)")
	}
	records := []Record{}
	seen := map[string]bool{}
	for {
		values, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, report, errors.New("malformed CSV; normalized output withheld")
		}
		report.InputRows++
		line, _ := reader.FieldPos(0)
		r, parseErr := parse(values)
		key := ""
		if parseErr == nil {
			// Numeric identifiers retain their original spelling. Canonical input must
			// use fixed identifier spellings; no cross-election correspondence is implied.
			key = strings.Join([]string{strconv.Itoa(r.ElectionYear), strconv.Itoa(r.Round), r.Office, r.UF, r.MunicipalityCode, r.Zone, r.Section}, "\x00")
			if seen[key] {
				parseErr = errors.New("duplicate_key")
			}
		}
		if parseErr != nil {
			report.RejectedRows++
			report.Rejections = append(report.Rejections, Rejection{line, parseErr.Error()})
			continue
		}
		seen[key] = true
		report.AcceptedRows++
		records = append(records, r)
	}
	if report.InputRows == 0 {
		return nil, report, errors.New("empty dataset does not validate ingestion")
	}
	report.QualityPassed = report.RejectedRows == 0
	if !report.QualityPassed {
		records = nil
	}
	return records, report, nil
}

func writeJSON(path string, value any) error {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return writeExclusive(path, append(content, '\n'))
}

func writeExclusive(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0444)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	return errors.Join(writeErr, closeErr)
}

// Run archives original bytes and provenance in a new run directory. Existing
// directories are never replaced, including on reruns. A completion marker is
// written last, so partial I/O failures cannot masquerade as completed runs.
func Run(input, manifestPath, output string, cutoff time.Time) (Report, error) {
	var report Report
	mf, err := os.Open(manifestPath)
	if err != nil {
		return report, err
	}
	m, err := ReadManifest(mf)
	closeErr := mf.Close()
	if err != nil || closeErr != nil {
		return report, errors.Join(err, closeErr)
	}
	file, err := os.Open(input)
	if err != nil {
		return report, err
	}
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		file.Close()
		return report, errors.New("input must be a regular file")
	}
	data, err := io.ReadAll(io.LimitReader(file, MaxBytes+1))
	closeErr = file.Close()
	if err != nil || closeErr != nil {
		return report, errors.Join(err, closeErr)
	}
	records, report, err := Import(data, m, cutoff)
	if err != nil {
		return report, err
	}
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return report, err
	}
	if err = os.Mkdir(output, 0755); err != nil {
		return report, err
	}
	if err = os.Mkdir(filepath.Join(output, "raw"), 0755); err != nil {
		return report, err
	}
	if err = writeExclusive(filepath.Join(output, "raw", m.SHA256+".csv"), data); err != nil {
		return report, err
	}
	if err = writeJSON(filepath.Join(output, "manifest.json"), m); err != nil {
		return report, err
	}
	if err = writeJSON(filepath.Join(output, "quality.json"), report); err != nil {
		return report, err
	}
	artifacts := map[string]string{"raw": m.SHA256}
	if report.QualityPassed {
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		for _, r := range records {
			if err = encoder.Encode(r); err != nil {
				return report, err
			}
		}
		if err = writeExclusive(filepath.Join(output, "normalized.jsonl"), buffer.Bytes()); err != nil {
			return report, err
		}
		digest := sha256.Sum256(buffer.Bytes())
		artifacts["normalized"] = hex.EncodeToString(digest[:])
	}
	if err = writeJSON(filepath.Join(output, "artifacts.json"), artifacts); err != nil {
		return report, err
	}
	err = writeJSON(filepath.Join(output, "complete.json"), map[string]any{"completed": true, "quality_passed": report.QualityPassed})
	return report, err
}
