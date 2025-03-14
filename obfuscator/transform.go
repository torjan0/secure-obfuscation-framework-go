package obfuscator

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// GenerateCrypticName creates a long, cryptic variable name
func GenerateCrypticName() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	length := 16
	name := make([]byte, length)
	for i := range name {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		name[i] = chars[n.Int64()]
	}
	return "x" + string(name) // Ensure it starts with a letter
}

// EncryptString uses AES to encrypt a string
func EncryptString(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// TransformSource applies advanced obfuscation transformations
func TransformSource(dir string, level string) error {
	fset := token.NewFileSet()
	key := make([]byte, 16) // 128-bit AES key
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("generating AES key: %v", err)
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing %s: %v", path, err)
		}

		// 1. Advanced Variable Name Mangling
		if level != "none" {
			ast.Inspect(file, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if !ident.IsExported() && ident.Obj != nil && ident.Obj.Kind == ast.Var {
						ident.Name = GenerateCrypticName()
					}
				}
				return true
			})
		}

		// 2. Dynamic String Encryption
		if level == "medium" || level == "heavy" {
			ast.Inspect(file, func(n ast.Node) bool {
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					plain := strings.Trim(lit.Value, `"`)
					encrypted, err := EncryptString(plain, key)
					if err == nil {
						lit.Value = fmt.Sprintf(`decryptString("%s")`, encrypted)
					}
				}
				return true
			})

			// Add decryption function
			decryptFunc := &ast.FuncDecl{
				Name: ast.NewIdent("decryptString"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Type: ast.NewIdent("string"), Names: []*ast.Ident{ast.NewIdent("s")}},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: ast.NewIdent("string")}},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("decryptAES"),
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun:  ast.NewIdent("base64.StdEncoding.DecodeString"),
											Args: []ast.Expr{ast.NewIdent("s")},
										},
										&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, base64.StdEncoding.EncodeToString(key))},
									},
								},
							},
						},
					},
				},
			}
			file.Decls = append(file.Decls, decryptFunc)

			// Add decryptAES helper
			decryptAESFunc := &ast.FuncDecl{
				Name: ast.NewIdent("decryptAES"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Type: ast.NewIdent("[]byte"), Names: []*ast.Ident{ast.NewIdent("data")}},
							{Type: ast.NewIdent("string"), Names: []*ast.Ident{ast.NewIdent("keyStr")}},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{{Type: ast.NewIdent("string")}},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("key")},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("base64.StdEncoding.DecodeString"), Args: []ast.Expr{ast.NewIdent("keyStr")}}},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("block")},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("aes.NewCipher"), Args: []ast.Expr{ast.NewIdent("key")}}},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("iv")},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.SliceExpr{X: ast.NewIdent("data"), Low: &ast.BasicLit{Kind: token.INT, Value: "0"}, High: &ast.BasicLit{Kind: token.INT, Value: "16"}}},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{ast.NewIdent("plaintext")},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{&ast.CallExpr{Fun: ast.NewIdent("make"), Args: []ast.Expr{ast.NewIdent("[]byte"), &ast.BinaryExpr{X: ast.NewIdent("len"), Op: token.SUB, Y: &ast.BasicLit{Kind: token.INT, Value: "16"}}}}},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun:  &ast.SelectorExpr{X: ast.NewIdent("cipher.NewCFBDecrypter"), Sel: ast.NewIdent("XORKeyStream")},
								Args: []ast.Expr{ast.NewIdent("block"), ast.NewIdent("iv"), ast.NewIdent("plaintext"), &ast.SliceExpr{X: ast.NewIdent("data"), Low: &ast.BasicLit{Kind: token.INT, Value: "16"}}},
							},
						},
						&ast.ReturnStmt{Results: []ast.Expr{ast.NewIdent("string(plaintext)")}},
					},
				},
			}
			file.Decls = append(file.Decls, decryptAESFunc)
		}

		// 3. Control Flow Flattening
		if level == "heavy" {
			for i, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					flattenControlFlow(fn)
				}
				file.Decls[i] = decl
			}
		}

		// 4. Dead Code Injection
		if level == "heavy" {
			for i, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					injectDeadCode(fn)
				}
				file.Decls[i] = decl
			}
		}

		// Write back to file
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("creating %s: %v", path, err)
		}
		defer f.Close()
		if err := printer.Fprint(f, fset, file); err != nil {
			return fmt.Errorf("writing %s: %v", path, err)
		}
		return nil
	})
	return err
}

// flattenControlFlow flattens a function's control flow
func flattenControlFlow(fn *ast.FuncDecl) {
	if len(fn.Body.List) == 0 {
		return
	}
	statements := fn.Body.List
	fn.Body.List = []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("state")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
		},
		&ast.ForStmt{
			Cond: &ast.BinaryExpr{X: ast.NewIdent("state"), Op: token.LSS, Y: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", len(statements))}},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.SwitchStmt{
						Tag: ast.NewIdent("state"),
						Body: &ast.BlockStmt{
							List: func() []ast.Stmt {
								var cases []ast.Stmt
								for i, stmt := range statements {
									cases = append(cases, &ast.CaseClause{
										List: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", i)}},
										Body: append([]ast.Stmt{stmt}, &ast.IncDecStmt{X: ast.NewIdent("state"), Tok: token.INC}),
									})
								}
								return cases
							}(),
						},
					},
				},
			},
		},
	}
}

// injectDeadCode adds meaningless statements
func injectDeadCode(fn *ast.FuncDecl) {
	deadCode := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  ast.NewIdent("fmt.Println"),
			Args: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: `"deadcode"`}},
		},
	}
	fn.Body.List = append([]ast.Stmt{deadCode}, fn.Body.List...)
}
