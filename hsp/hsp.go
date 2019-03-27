package main

// hsp is a code generation tool for
// creating methods to serialize and de-serialize
// Go data structures to and from MessagePack.
//
// This package is targeted at the `go generate` tool.
// To use it, include the following directive in a
// go source file with types requiring source generation:
//
//     //go:generate hsp
//
// The go generate tool should set the proper environment variables for
// the generator to execute without any command-line flags. However, the
// following options are supported, if you need them:
//
//  -o = output file name (default is {input}_gen.go)
//  -file = input file name (or directory; default is $GOFILE, which is set by the `go generate` command)
//  -tests = generate tests and benchmarks (default is true)
//
// For more information, please read README.md, and the wiki at github.com/CovenantSQL/HashStablePack
//

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttacon/chalk"

	"github.com/CovenantSQL/HashStablePack/gen"
	"github.com/CovenantSQL/HashStablePack/parse"
	"github.com/CovenantSQL/HashStablePack/printer"
)

var (
	out        = flag.String("o", "", "output file")
	file       = flag.String("file", "", "input file")
	tests      = flag.Bool("tests", true, "create tests and benchmarks")
	unexported = flag.Bool("unexported", false, "also process unexported types")
)

func main() {
	flag.Parse()

	// GOFILE is set by go generate
	if *file == "" {
		*file = os.Getenv("GOFILE")
		if *file == "" {
			fmt.Println(chalk.Red.Color("No file to parse."))
			os.Exit(1)
		}
	}
	flag.Visit(func(f *flag.Flag) {
		fmt.Printf("args %#v : %s", f.Name, f.Value)
	})
	fmt.Printf("file: %s\n", *file)

	var mode gen.Method
	mode |= gen.Marshal | gen.Size
	if *tests {
		mode |= gen.Test
	}

	if mode&^gen.Test == 0 {
		fmt.Println(chalk.Red.Color("No methods to generate; -io=false && -marshal=false"))
		os.Exit(1)
	}

	if err := Run(*file, mode, *unexported); err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
}

// Run writes all methods using the associated file or path, e.g.
//
//	err := hsp.Run("path/to/myfile.go", gen.Size|gen.Marshal|gen.Unmarshal|gen.Test, false)
//
func Run(gofile string, mode gen.Method, unexported bool) error {
	if mode&^gen.Test == 0 {
		return nil
	}
	fmt.Println(chalk.Magenta.Color("======== HashStablePack Code Generator ======="))
	fmt.Printf(chalk.Magenta.Color(">>> Input: \"%s\"\n"), gofile)
	fs, err := parse.File(gofile, unexported)
	if err != nil {
		return err
	}

	if len(fs.Identities) == 0 {
		fmt.Println(chalk.Magenta.Color("No types requiring code generation were found!"))
		return nil
	}

	var versionTypes []*gen.Struct

	for _, el := range fs.Identities {
		if st, ok := el.(*gen.Struct); ok {
			if st.Versioning {
				versionTypes = append(versionTypes, st)
			}
		}
	}

	genFileName := newFilename(gofile, fs.Package)

	if len(versionTypes) > 0 {
		// should parse existing _gen.go for old version data
		if err := parse.ParseOldGenFile(genFileName, versionTypes); err != nil {
			return err
		}

		// set numeric versions
		for _, st := range versionTypes {
			found := false
			for i, v := range st.VersionList {
				if v == st.CurrentVersion {
					found = true
					st.CurrentNumericVersion = i
					break
				}
			}
			if !found {
				st.CurrentNumericVersion = len(st.VersionList)
				st.VersionList = append(st.VersionList, st.CurrentVersion)
			}

			// print version type files
			if err := printer.PrintVersionFile(genFileName, fs, st, mode); err != nil {
				return err
			}

			if st.OldMarshalBody != "" && st.OldMsgSizeBody != "" {
				if err := printer.PrintOldVersionFile(genFileName, fs, st, mode); err != nil {
					return err
				}
			}
		}
	}

	return printer.PrintFile(genFileName, fs, mode)
}

// picks a new file name based on input flags and input filename(s).
func newFilename(old string, pkg string) string {
	if *out != "" {
		if pre := strings.TrimPrefix(*out, old); len(pre) > 0 &&
			!strings.HasSuffix(*out, ".go") {
			return filepath.Join(old, *out)
		}
		return *out
	}

	if fi, err := os.Stat(old); err == nil && fi.IsDir() {
		old = filepath.Join(old, pkg)
	}
	// new file name is old file name + _gen.go
	return strings.TrimSuffix(old, ".go") + "_gen.go"
}
