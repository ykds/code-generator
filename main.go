package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "codegen",
		Usage: "代码生成工具",
		Commands: []*cli.Command{
			generate(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generate() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "生成代码",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "model",
				Aliases:  []string{"m"},
				Usage:    "model目录路径",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "输出目录",
				Value:   ".",
			},
		},
		Action: func(c *cli.Context) error {
			return Generate(Config{
				ModelPath:  c.String("model"),
				OutputPath: c.String("output"),
			})
		},
	}
}
