package printer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ttacon/chalk"
	"golang.org/x/tools/imports"

	"github.com/CovenantSQL/HashStablePack/gen"
	"github.com/CovenantSQL/HashStablePack/parse"
)

func infof(s string, v ...interface{}) {
	fmt.Printf(chalk.Magenta.Color(s), v...)
}

// PrintFile prints the methods for the provided list
// of elements to the given file name and canonical
// package path.
func PrintFile(file string, f *parse.FileSet, mode gen.Method) error {
	out, tests, err := generate(f, mode)
	if err != nil {
		return err
	}

	// we'll run goimports on the main file
	// in another goroutine, and run it here
	// for the test file. empirically, this
	// takes about the same amount of time as
	// doing them in serial when GOMAXPROCS=1,
	// and faster otherwise.
	res := goformat(file, out.Bytes())
	if tests != nil {
		testfile := strings.TrimSuffix(file, ".go") + "_test.go"
		err = format(testfile, tests.Bytes())
		if err != nil {
			return err
		}
		infof(">>> Wrote and formatted \"%s\"\n", testfile)
	}
	err = <-res
	if err != nil {
		return err
	}
	return nil
}

// PrintVersionFile prints the method for the provide versioned type.
func PrintVersionFile(file string, f *parse.FileSet, s *gen.Struct, mode gen.Method) error {
	out, tests, err := generateVersion(f, s, mode)
	if err != nil {
		return err
	}

	genFileName := strings.TrimSuffix(file, "_gen.go") + "_" +
		strings.ToLower(s.TypeName()) + "_" + s.CurrentVersion + "_gen.go"
	res := goformat(genFileName, out.Bytes())
	if tests != nil {
		testfile := strings.TrimSuffix(genFileName, ".go") + "_test.go"
		err = format(testfile, tests.Bytes())
		if err != nil {
			return err
		}
		infof(">>> Wrote and formatted \"%s\"\n", testfile)
	}
	err = <-res
	if err != nil {
		return err
	}
	return nil
}

// PrintOldVersionFile prints the method for the provide versioned type.
func PrintOldVersionFile(file string, f *parse.FileSet, s *gen.Struct, mode gen.Method) error {
	out, tests, err := generateOldVersion(f, s, mode)
	if err != nil {
		return err
	}

	genFileName := strings.TrimSuffix(file, "_gen.go") + "_" +
		strings.ToLower(s.TypeName()) + "_oldver_gen.go"
	res := goformat(genFileName, out.Bytes())
	if tests != nil {
		testfile := strings.TrimSuffix(genFileName, ".go") + "_test.go"
		err = format(testfile, tests.Bytes())
		if err != nil {
			return err
		}
		infof(">>> Wrote and formatted \"%s\"\n", testfile)
	}
	err = <-res
	if err != nil {
		return err
	}
	return nil
}

func format(file string, data []byte) error {
	out, err := imports.Process(file, data, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, out, 0600)
}

func goformat(file string, data []byte) <-chan error {
	out := make(chan error, 1)
	go func(file string, data []byte, end chan error) {
		end <- format(file, data)
		infof(">>> Wrote and formatted \"%s\"\n", file)
	}(file, data, out)
	return out
}

func dedupImports(imp []string) []string {
	m := make(map[string]struct{})
	for i := range imp {
		m[imp[i]] = struct{}{}
	}
	r := []string{}
	for k := range m {
		r = append(r, k)
	}
	return r
}

func generate(f *parse.FileSet, mode gen.Method) (*bytes.Buffer, *bytes.Buffer, error) {
	outbuf := bytes.NewBuffer(make([]byte, 0, 4096))
	writePkgHeader(outbuf, f.Package)

	myImports := []string{}
	myImports = append(myImports, `hsp "github.com/CovenantSQL/HashStablePack/marshalhash"`)

	for _, imp := range f.Imports {
		if imp.Name != nil {
			// have an alias, include it.
			myImports = append(myImports, imp.Name.Name+` `+imp.Path.Value)
		} else {
			myImports = append(myImports, imp.Path.Value)
		}
	}
	dedup := dedupImports(myImports)
	writeImportHeader(outbuf, dedup...)

	var testbuf *bytes.Buffer
	var testwr io.Writer
	if mode&gen.Test == gen.Test {
		testbuf = bytes.NewBuffer(make([]byte, 0, 4096))
		writePkgHeader(testbuf, f.Package)
		writeImportHeader(testbuf, "bytes", "crypto/rand", "encoding/binary", `hsp "github.com/CovenantSQL/HashStablePack/marshalhash"`, "testing")
		testwr = testbuf
	}
	return outbuf, testbuf, f.PrintTo(gen.NewPrinter(mode, outbuf, testwr, ""))
}

func generateOldVersion(f *parse.FileSet, s *gen.Struct, mode gen.Method) (*bytes.Buffer, *bytes.Buffer, error) {
	outbuf := bytes.NewBuffer(make([]byte, 0, 4096))
	writePkgHeader(outbuf, f.Package)

	myImports := []string{}
	myImports = append(myImports, `hsp "github.com/CovenantSQL/HashStablePack/marshalhash"`)

	for _, imp := range f.Imports {
		if imp.Name != nil {
			// have an alias, include it.
			myImports = append(myImports, imp.Name.Name+` `+imp.Path.Value)
		} else {
			myImports = append(myImports, imp.Path.Value)
		}
	}
	dedup := dedupImports(myImports)
	writeImportHeader(outbuf, dedup...)

	var testbuf *bytes.Buffer
	var testwr io.Writer
	if mode&gen.Test == gen.Test {
		testbuf = bytes.NewBuffer(make([]byte, 0, 4096))
		writePkgHeader(testbuf, f.Package)
		writeImportHeader(testbuf, "bytes", "crypto/rand", "encoding/binary", `hsp "github.com/CovenantSQL/HashStablePack/marshalhash"`, "testing")
		testwr = testbuf
	}
	return outbuf, testbuf, f.PrintVersion(s, gen.NewPrinter(mode, outbuf, testwr, "oldver"), "oldver")
}

func generateVersion(f *parse.FileSet, s *gen.Struct, mode gen.Method) (*bytes.Buffer, *bytes.Buffer, error) {
	outbuf := bytes.NewBuffer(make([]byte, 0, 4096))
	writePkgHeader(outbuf, f.Package)

	myImports := []string{}
	myImports = append(myImports, `hsp "github.com/CovenantSQL/HashStablePack/marshalhash"`)

	for _, imp := range f.Imports {
		if imp.Name != nil {
			// have an alias, include it.
			myImports = append(myImports, imp.Name.Name+` `+imp.Path.Value)
		} else {
			myImports = append(myImports, imp.Path.Value)
		}
	}
	dedup := dedupImports(myImports)
	writeImportHeader(outbuf, dedup...)

	var testbuf *bytes.Buffer
	var testwr io.Writer
	if mode&gen.Test == gen.Test {
		testbuf = bytes.NewBuffer(make([]byte, 0, 4096))
		writePkgHeader(testbuf, f.Package)
		writeImportHeader(testbuf, "bytes", "crypto/rand", "encoding/binary", `hsp "github.com/CovenantSQL/HashStablePack/marshalhash"`, "testing")
		testwr = testbuf
	}
	return outbuf, testbuf, f.PrintVersion(s, gen.NewPrinter(mode, outbuf, testwr, s.CurrentVersion), s.CurrentVersion)
}

func writePkgHeader(b *bytes.Buffer, name string) {
	b.WriteString("package ")
	b.WriteString(name)
	b.WriteByte('\n')
	// write generated code marker
	// https://github.com/tinylib/hsp/issues/229
	// https://golang.org/s/generatedcode
	b.WriteString("// Code generated by github.com/CovenantSQL/HashStablePack DO NOT EDIT.\n\n")
}

func writeImportHeader(b *bytes.Buffer, imports ...string) {
	b.WriteString("import (\n")
	for _, im := range imports {
		if im[len(im)-1] == '"' {
			// support aliased imports
			fmt.Fprintf(b, "\t%s\n", im)
		} else {
			fmt.Fprintf(b, "\t%q\n", im)
		}
	}
	b.WriteString(")\n\n")
}
