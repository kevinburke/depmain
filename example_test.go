package depmain_test

import (
	"flag"
	"io"
	"os"

	"github.com/onemedical/depmain"
)

func _main(ext *depmain.Ext) int {
	flag.CommandLine.Parse(ext.Args)
	io.WriteString(ext.Stdout, "Running program")
	return 0
}

func Example() {
	os.Exit(_main(depmain.New()))
}
