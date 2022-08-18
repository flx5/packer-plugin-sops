package main

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"os"
	scaffoldingData "packer-plugin-sops/datasource/file"
)

func main() {
	pps := plugin.NewSet()
	// TODO Set Version

	pps.RegisterDatasource("file", new(scaffoldingData.DataSource))
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
