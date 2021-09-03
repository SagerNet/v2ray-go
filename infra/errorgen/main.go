package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func generateError(path string) {
	pkg := filepath.Base(path)
	if pkg == "v2ray-go" {
		pkg = "core"
	}

	file, err := os.OpenFile(path+"/errors.generated.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		fmt.Printf("Failed to generate errors.generated.go: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintf(file, `package %s

import "github.com/v2fly/v2ray-core/v4/common/errors"

type errPathObjHolder struct{}

func newError(values ...interface{}) *errors.Error {
	return errors.New(values...).WithPathObj(errPathObjHolder{})
}
`, pkg)
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("can not get current working directory")
		os.Exit(1)
	}

	generateError(pwd)

	genPkgs := []string{"app", "infra", "common", "features", "main", "proxy", "transport"}
	noGenPkgs := []string{
		"common/errors",
		"common/log",
		"common/platform",
		"common/serial",
		"common/signal",
		"common/signal/done",
		"common/signal/semaphore",
		"infra/errorgen",
		"infra/vprotogen"}

	for _, c := range genPkgs {
		walkErr := filepath.Walk(c, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}

			if !info.IsDir() {
				return nil
			}

			for _, noGen := range noGenPkgs {
				if path == noGen {
					return nil
				}
			}

			println(path)
			generateError("./" + "/" + path)

			return nil
		})

		if walkErr != nil {
			fmt.Println(walkErr)
			os.Exit(1)
		}
	}
}
