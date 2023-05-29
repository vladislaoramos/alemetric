package exitcheck

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestExitCheck(t *testing.T) {
	testdata := []struct {
		name     string
		filename string
	}{
		{
			name:     "ExitInMain",
			filename: "testdata/pkg1/pkg1.go",
		},
		{
			name:     "NoExitInMain",
			filename: "testdata/pkg2/pkg2.go",
		},
	}

	analyzer := Analyzer

	for _, tc := range testdata {
		t.Run(tc.name, func(t *testing.T) {
			file, err := parseFile(t, tc.filename)
			require.NoError(t, err, err)

			pass := &analysis.Pass{
				Analyzer: analyzer,
				Fset:     token.NewFileSet(),
				Files:    []*ast.File{file},
				Report:   func(analysis.Diagnostic) {},
			}

			_, err = analyzer.Run(pass)
			require.NoError(t, err, "analyzer.Run() error: %v", err)
		})
	}
}

func parseFile(t *testing.T, filename string) (*ast.File, error) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("error parsing file %s: %v", filename, err)
	}

	return file, nil
}
