package main

import (
	"flag"
	"log"
	"os"

	"roob.re/m3u"
)

func main() {
	allowDupes := flag.Bool("allow-duplicates", false, "allow duplicates in output")
	flag.Parse()

	all := m3u.Playlist{}

	for _, path := range flag.Args() {
		func() {
			file, err := openCLI(path)
			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()

			playlist, err := m3u.Decode(file)
			if err != nil {
				log.Printf("Error parsing %q, ignoring: %v", path, err)
			}

			all = all.Join(playlist)
		}()
	}

	if !*allowDupes {
		all = all.Deduplicate()
	}

	_ = all.Encode(os.Stdout)
}

func openCLI(path string) (*os.File, error) {
	if path == "-" {
		return os.Stdin, nil
	}

	return os.Open(path)
}
