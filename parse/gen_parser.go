package parse

import (
	"bytes"
	"github.com/CovenantSQL/HashStablePack/gen"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func ParseOldGenFile(f string, versionTypes []*gen.Struct) (err error) {
	pushstate(f)
	defer popstate()

	fset := token.NewFileSet()
	st, err := os.Stat(f)
	if err != nil || st.IsDir() {
		err = nil
		return
	}

	fl, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
	if err != nil {
		return
	}

	// build version type map
	versionTypeMap := map[string]*gen.Struct{}
	for _, st := range versionTypes {
		versionTypeMap[st.TypeName()] = st
	}

	// search for version types methods
	for i := range fl.Decls {
		if g, ok := fl.Decls[i].(*ast.GenDecl); ok {
			for _, s := range g.Specs {
				if vs, ok := s.(*ast.ValueSpec); ok {
					for i := range vs.Names {
						varName := vs.Names[i].String()
						if strings.HasPrefix(varName, "hspVersions") {
							versionType := varName[len("hspVersions"):]

							var genType *gen.Struct
							if genType, ok = versionTypeMap[versionType]; !ok {
								continue
							}

							if cpt, ok := vs.Values[i].(*ast.CompositeLit); ok {
								if _, ok := cpt.Type.(*ast.ArrayType); !ok {
									continue
								}

								for _, e := range cpt.Elts {
									if sl, ok := e.(*ast.BasicLit); ok && sl.Kind == token.STRING {
										genType.VersionList = append(genType.VersionList, strings.Trim(sl.Value, "\""))
									}
								}
							}
						}
					}
				}
			}
		}
	}

	for i := range fl.Decls {
		if fk, ok := fl.Decls[i].(*ast.FuncDecl); ok {
			fn := fk.Name.String()

			if fn != "MarshalHash" && fn != "Msgsize" {
				continue
			}

			tp := fk.Recv.List[0].Type

			for {
				if sexpr, ok := tp.(*ast.StarExpr); ok {
					tp = sexpr.X
				} else {
					break
				}
			}

			if tid, ok := tp.(*ast.Ident); ok {
				var genType *gen.Struct
				if genType, ok = versionTypeMap[tid.String()]; !ok {
					continue
				}

				if len(genType.VersionList) > 0 {
					// already converted to new version
					continue
				}

				var fBytes bytes.Buffer
				_ = printer.Fprint(&fBytes, fset, fk.Body)

				if fn == "MarshalHash" {
					genType.OldMarshalBody = fBytes.String()
				} else if fn == "Msgsize" {
					genType.OldMsgSizeBody = fBytes.String()
				}

				if genType.OldMarshalBody != "" && genType.OldMsgSizeBody != "" {
					genType.VersionList = append(genType.VersionList, "oldver")
				}
			}
		}
	}

	return
}
