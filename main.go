/*
Example command line usage

    $ ./timeat 'los angeles'
    Los Angeles, CA, USA: Fri May 23, 2014 22:33
    $ ./timeat paris
    Paris, France: Sat May 24, 2014 07:37
    Paris, TX, USA: Sat May 24, 2014 00:37
    Paris, TN 38242, USA: Sat May 24, 2014 00:37
    Paris, IL 61944, USA: Sat May 24, 2014 00:37
    Paris, KY 40361, USA: Sat May 24, 2014 01:37
    $
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/tebeka/timeat/lib"
)

// die prints error message and aborts the program
func die(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

func main() {
	version := flag.Bool("version", false, "show version")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s ADDRESS\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if *version {
		fmt.Println(timeat.Version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		die("wrong number of arguments")
	}

	address := flag.Arg(0)
	timesat, err := timeat.TimeAt(address)
	if err != nil {
		die("error: can't get times for %s - %s", address, err)
	}

	if len(timesat) == 0 {
		die("error: no locations found matching %s", address)
	}

	for _, ta := range timesat {
		fmt.Printf("%s: %s\n", ta.Address, ta.Time.Format("Mon Jan 2, 2006 15:04"))
	}
}
