package verdeps

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"sync"
)

// goFileParser parses an AST from an go file.
type goFileParser func(
	fset *token.FileSet,
	filename string,
	src interface{},
	mode parser.Mode,
) (f *ast.File, err error)

// parseGoFileArgs is the arguments struct for parseGoFileArgs.
type parseGoFileArgs struct {
	parse           goFileParser
	errors          *syncedErrors
	filePath        string
	waitGroup       *sync.WaitGroup
	importCounts    *syncedImportCounts
	vendorContext   *vendorContext
	importSpecChan  chan *importSpec
	packageSpecChan chan *packageSpec
}

// parseGoFile uses the golang compiler's AST parser to parse out the import
// specs and package spec of a file written in Go. The specs are then returned
// via the import spec and package spec channels.
func parseGoFile(args parseGoFileArgs) {
	var (
		f     *ast.File
		err   error
		specs []*importSpec
	)

	// Signal that the wait group should continue when this function exits.
	defer args.waitGroup.Done()

	// Parse the imports of the file.
	if f, err = args.parse(
		token.NewFileSet(),
		args.filePath,
		nil,
		parser.ImportsOnly,
	); err != nil {
		args.errors.add(err)
		return
	}

	// Filter the specs.
	for _, spec := range f.Imports {
		// Ignore the surrounding quotes.
		importString := strings.Trim(spec.Path.Value, "\"")

		// Only pursue a dependency if it has a github prefix and is not vendored.
		if !args.vendorContext.contains(importString) &&
			strings.HasPrefix(importString, "github.com/") {
			// Both conditions were met, so add this import spec to the list.
			specs = append(specs, &importSpec{
				imports:  spec,
				filePath: args.filePath,
			})
		}
	}

	// Set the import count for this file path.
	if args.importCounts != nil {
		// TODO(skeswa): figure out if this is necessary since it also appears in
		// the buffer_vendorables_test.go.
		args.importCounts.setImportCount(args.filePath, len(specs))
	}

	// Throw all the newly discovered specs into the mix.
	args.packageSpecChan <- &packageSpec{
		filePath:   args.filePath,
		startIndex: int(f.Package),
	}
	for _, spec := range specs {
		args.importSpecChan <- spec
	}

	return
}
