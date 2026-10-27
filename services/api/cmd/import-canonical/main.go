// import-canonical accepts only the explicitly synthetic v0.1 fixture contract.
package main

import (
	"flag"
	"fmt"
	"observatorio.local/api/internal/historical"
	"os"
	"time"
)

func main() {
	input := flag.String("input", "", "local canonical CSV (SYNTHETIC only)")
	manifest := flag.String("manifest", "", "local provenance manifest JSON")
	output := flag.String("output", "", "new output directory; never overwritten")
	cutoff := flag.String("cutoff", "", "required RFC3339 availability cutoff")
	flag.Parse()
	if *input == "" || *manifest == "" || *output == "" || flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}
	instant, err := time.Parse(time.RFC3339, *cutoff)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid or missing --cutoff (RFC3339 required)")
		os.Exit(1)
	}
	report, err := historical.Run(*input, *manifest, *output, instant)
	if err != nil {
		fmt.Fprintln(os.Stderr, "import failed:", err)
		os.Exit(1)
	}
	fmt.Printf("SYNTHETIC: rows=%d accepted=%d rejected=%d\n", report.InputRows, report.AcceptedRows, report.RejectedRows)
	if !report.QualityPassed {
		os.Exit(2)
	}
}
