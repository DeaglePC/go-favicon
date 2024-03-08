// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-09

// Command favicon finds favicons for a URL and prints their format, size and URL.
// It can search HTML tags, manifest files and common locations on the server.
// Results are printed in a pretty table by default, but can be dumped as
// JSON, CSV or TSV.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"strings"

	"github.com/thanhpk/go-favicon"
)

var (
	// set by Makefile
	version   = "undefined"
	buildDate = "undefined"
)

var (
	fs          = flag.NewFlagSet("favicon", flag.ExitOnError)
	flagHelp    = fs.Bool("h", false, "show this message and exit")
	flagJSON    = fs.Bool("json", false, "output favicon list as JSON")
	flagVersion = fs.Bool("version", false, "show version number and exit")

	log  *stdlog.Logger
	opts []favicon.Option
)

func usage() {
	fmt.Fprint(fs.Output(),
		`usage: favicon <url>

retrieve, favicons for URL

`)
	fs.PrintDefaults()
}

func init() {
	fs.Usage = usage
	log = stdlog.New(os.Stderr, "", 0)
}

func main() {
	checkErr(fs.Parse(os.Args[1:]))

	if *flagVersion {
		fmt.Printf("favicon %s\n", version)
		fmt.Printf("built: %s\n", buildDate)
		return
	}

	if *flagHelp || fs.NArg() == 0 {
		usage()
		return
	}

	u := fs.Arg(0)
	s := strings.ToLower(u)
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		log.Fatalf("invalid URL: %q", s)
	}

	f := favicon.New(opts...)
	icons, err := f.Find(u)
	checkErr(err)
	// log.Printf("%d icon(s) found for %q", len(icons), u)

	if *flagJSON {
		data, err := json.MarshalIndent(icons, "", "  ")
		checkErr(err)
		fmt.Println(string(data))
		return
	}

	for i, icon := range icons {
		fmt.Printf("%3d: %s [%s]\n", i+1, icon.URL, icon.MimeType)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
