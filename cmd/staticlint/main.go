package main

//go:generate go build -o=../../bin/staticlint

import (
	exitcheck "github.com/vladislaoramos/alemetric/internal/analyzer"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

var analyzers = []*analysis.Analyzer{
	// reports mismatches between assembly files and Go declaration
	asmdecl.Analyzer,
	// detects useless assignments
	assign.Analyzer,
	// checks for common mistakes using the sync/atomic package
	atomic.Analyzer,
	// checks for non-64-bit-aligned arguments to sync/atomic functions
	atomicalign.Analyzer,
	// checks for common mistakes involving boolean operators
	bools.Analyzer,
	// checks build tags
	buildtag.Analyzer,
	// checks for unkeyed composite literals
	composite.Analyzer,
	// checks for locks erroneously passed by value
	copylock.Analyzer,
	// checks for the use of reflect.DeepEqual with error values
	deepequalerrors.Analyzer,
	// checks that the second argument to errors.As is a pointer to a type implementing error
	errorsas.Analyzer,
	// checks for mistakes using HTTP responses
	httpresponse.Analyzer,
	// checks for references to enclosing loop variables from within nested functions
	loopclosure.Analyzer,
	// checks for failure to call a context cancellation function
	lostcancel.Analyzer,
	// checks for useless comparisons against nil
	nilfunc.Analyzer,
	// inspects the control-flow graph of an SSA function
	// and reports errors such as nil pointer dereferences
	// and degenerate nil pointer comparisons
	nilness.Analyzer,
	// checks consistency of Printf format strings and arguments
	printf.Analyzer,
	// checks for shifts that exceed the width of an integer
	shift.Analyzer,
	// checks for misspellings in the signatures of methods similar to well-known interfaces
	stdmethods.Analyzer,
	// checks struct field tags are well-formed
	structtag.Analyzer,
	// checks for common mistaken usages of tests and examples
	tests.Analyzer,
	// checks for passing non-pointer
	// or non-interface types to unmarshal and decode functions
	unmarshal.Analyzer,
	// checks for unreachable code
	unreachable.Analyzer,
	// checks for invalid conversions of uintptr to unsafe.Pointer
	unsafeptr.Analyzer,
	// checks for unused results of calls to certain pure functions
	unusedresult.Analyzer,
	// checks for unused writes to the elements of a struct or array object
	unusedwrite.Analyzer,
	// flags type conversions from integers to strings
	stringintconv.Analyzer,
	// flags impossible interface-interface type assertions
	ifaceassert.Analyzer,

	// checks os.Exit() calls inside the main.main()
	exitcheck.Analyzer,
}

var additionalAnalyzers = map[string]struct{}{
	// The documentation of an exported function should start with the function’s name
	"ST1020": {},
	// Omit redundant type from variable declaration
	"QF1011": {},
}

func main() {
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		} else if _, ok := additionalAnalyzers[v.Analyzer.Name]; ok {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	multichecker.Main(
		analyzers...,
	)
}
