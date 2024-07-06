package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"roob.re/m3u"
)

func main() {
	allowDupes := flag.Bool("allow-duplicates", false, "allow duplicates in output")
	flag.Parse()

	globalDirectivePrefixes, _ := os.LookupEnv("M3UCAT_GLOBAL_DIRECTIVE_PREFIXES")
	decoder := m3u.Decoder{IsGlobalDirective: globalDirectivePrefixFromEnv(globalDirectivePrefixes)}

	all := m3u.Playlist{}

	for _, path := range flag.Args() {
		func() {
			file, err := openCLI(path)
			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()

			playlist, err := decoder.Decode(file)
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

// globalDirectivePrefixFromEnv returns a function to be used on an m3u decoder, which itself returns true if a
// directive matches any of the comma-separated prefixes on env.
func globalDirectivePrefixFromEnv(env string) func(string) bool {
	if env == "" {
		return func(_ string) bool { return false }
	}

	list := strings.Split(env, ",")

	return func(s string) bool {
		for _, prefix := range list {
			if strings.HasPrefix(s, prefix) {
				return true
			}
		}

		return false
	}
}
