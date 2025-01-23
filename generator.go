package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Config struct {
	ModelPath  string
	OutputPath string
}

func Generate(cfg Config) error {
	// 1. 解析model目录下的所有结构体
	structs, err := ParseModelDir(cfg.ModelPath)
	if err != nil {
		return fmt.Errorf("解析model目录失败: %w", err)
	}

	// 2. 为每个结构体生成代码
	for _, st := range structs {
		if err := generateForStruct(cfg, st); err != nil {
			return fmt.Errorf("生成结构体 %s 的代码失败: %w", st.Name, err)
		}
	}

	return nil
}

type Field struct {
	Name     string
	Type     string
	Tag      string
	Comment  string
	JsonName string // 从json tag中解析
}

type Struct struct {
	Name    string
	Fields  []*Field
	Comment string
}

func ParseModelDir(dir string) ([]*Struct, error) {
	var structs []*Struct

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// 解析Go文件
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("解析文件失败 %s: %w", path, err)
		}

		// 遍历AST，查找结构体定义
		ast.Inspect(node, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}

			// 创建结构体对象
			st := &Struct{
				Name:   typeSpec.Name.Name,
				Fields: make([]*Field, 0),
			}

			// 获取结构体注释
			if typeSpec.Doc != nil {
				st.Comment = typeSpec.Doc.Text()
			}

			// 解析字段
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				f := &Field{
					Name: field.Names[0].Name,
					Type: typeToString(field.Type),
				}

				// 解析字段注释
				if field.Doc != nil {
					f.Comment = field.Doc.Text()
				}

				// 解析字段标签
				if field.Tag != nil {
					f.Tag = field.Tag.Value
					// 解析json tag
					if jsonTag := getTagValue(f.Tag, "json"); jsonTag != "" {
						f.JsonName = strings.Split(jsonTag, ",")[0]
					}
				}

				st.Fields = append(st.Fields, f)
			}

			structs = append(structs, st)
			return true
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return structs, nil
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// 获取结构体标签中指定key的值
func getTagValue(tag string, key string) string {
	tag = strings.Trim(tag, "`")
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		parts := strings.Split(t, ":")
		if len(parts) == 2 && parts[0] == key {
			return strings.Trim(parts[1], "\"")
		}
	}
	return ""
}

func generateForStruct(cfg Config, st *Struct) error {
	// 创建必要的目录
	dirs := []string{
		filepath.Join(cfg.OutputPath, "internal/repository"),
		filepath.Join(cfg.OutputPath, "internal/service"),
		filepath.Join(cfg.OutputPath, "internal/handler"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 生成各层代码
	generators := []generator{
		{templateName: "repository", outputDir: "repository", template: repositoryTemplate},
		{templateName: "service", outputDir: "service", template: serviceTemplate},
		{templateName: "handler", outputDir: "handler", template: handlerTemplate},
	}

	for _, g := range generators {
		if err := g.generate(st); err != nil {
			return fmt.Errorf("生成%s代码失败: %w", g.templateName, err)
		}
	}

	return nil
}

type generator struct {
	templateName string
	outputDir    string
	template     string
}

func (g *generator) generate(st *Struct) error {
	outputFile := filepath.Join("internal", g.outputDir, fmt.Sprintf("%s_%s.go", strings.ToLower(st.Name), g.outputDir))
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建%s文件失败: %w", g.templateName, err)
	}
	defer f.Close()

	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
	}

	tmpl, err := template.New(g.templateName).
		Funcs(funcMap).
		Parse(g.template)
	if err != nil {
		return fmt.Errorf("解析%s模板失败: %w", g.templateName, err)
	}

	if err := tmpl.Execute(f, st); err != nil {
		return fmt.Errorf("生成%s代码失败: %w", g.templateName, err)
	}

	return nil
}
