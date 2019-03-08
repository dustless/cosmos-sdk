package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	progName string
	// sections name-title map
	sections = map[string]string{
		"breaking":     "Breaking Changes",
		"features":     "New features",
		"improvements": "Improvements",
		"bugfixes":     "Bugfixes",
	}
	// stanzas name-title map
	stanzas = map[string]string{
		"gaia":       "Gaia",
		"gaiacli":    "Gaia CLI",
		"gaiarest":   "Gaia REST API",
		"sdk":        "SDK",
		"tendermint": "Tendermint",
	}
)

func init() {
	progName = filepath.Base(os.Args[0])
	flag.Usage = printUsage
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", filepath.Base(progName)))
	flag.Parse()

	if flag.NArg() < 2 {
		errInsufficientArgs()
	}

	changesDir := flag.Arg(0)
	cmd := flag.Arg(1)
	switch cmd {

	case "init":
		initChangesDir(changesDir)

	case "add":
		if flag.NArg() < 4 {
			errInsufficientArgs()
		}
		if flag.NArg() > 4 {
			errTooManyArgs()
		}
		addEntryFile(changesDir, flag.Arg(2), flag.Arg(3))

	case "generate":
		version := "UNRELEASED"
		if flag.NArg() > 2 {
			version = strings.Join(flag.Args()[2:], " ")
		}
		generateChangelog(changesDir, version)

	default:
		unknownCommand(cmd)
	}
}

func initChangesDir(changesDir string) {
	if _, err := os.Stat(changesDir); os.IsNotExist(err) {
		if err := os.Mkdir(changesDir, 0755); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("directory %q already exists", changesDir)
	}

	for sectionDir := range sections {
		for stanzaDir := range stanzas {
			path := filepath.Join(changesDir, sectionDir, stanzaDir)
			if err := os.MkdirAll(path, 0755); err != nil {
				log.Fatal(err)
			}
			// create stamp file to allow git commit of the dir structure
			if _, err := os.Create(filepath.Join(path, ".stamp")); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// nolint: errcheck
func addEntryFile(changesDir, sectionDir, stanzaDir string) {
	if _, ok := sections[sectionDir]; !ok {
		log.Fatalf("invalid section -- %s", sectionDir)
	}
	if _, ok := stanzas[stanzaDir]; !ok {
		log.Fatalf("invalid stanza -- %s", stanzaDir)
	}

	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	filename := generateFileName(strings.Split(string(bs), "\n")[0])
	filepath := filepath.Join(changesDir, sectionDir, stanzaDir, filename)
	outFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if _, err := outFile.Write(bs); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "Written to %s\n", filepath)
}

var filenameInvalidChars = regexp.MustCompile(`[^a-zA-Z0-9-_]`)

func generateFileName(line string) string {
	var chunks []string
	subsWithInvalidCharsRemoved := strings.Split(filenameInvalidChars.ReplaceAllString(line, " "), " ")
	for _, sub := range subsWithInvalidCharsRemoved {
		sub = strings.TrimSpace(sub)
		if len(sub) != 0 {
			chunks = append(chunks, sub)
		}
	}

	return strings.Join(chunks, "-")
}

func generateChangelog(changesDir, version string) {
	fmt.Printf("# %s\n\n", version)
	for sectionDir, sectionTitle := range sections {

		fmt.Printf("## %s\n\n", sectionTitle)
		for stanzaDir, stanzaTitle := range stanzas {
			fmt.Printf("### %s\n\n", stanzaTitle)
			path := filepath.Join(changesDir, sectionDir, stanzaDir)
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				if f.Name()[0] == '.' {
					continue
				}
				if err := indentAndPrintFile(filepath.Join(path, f.Name())); err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println()
		}
	}
}

// nolint: errcheck
func indentAndPrintFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if firstLine {
			fmt.Printf("* %s\n", scanner.Text())
			firstLine = false
			continue
		}

		fmt.Printf("  %s\n", scanner.Text())
	}

	return scanner.Err()
}

func printUsage() {
	usageText := fmt.Sprintf(`usage: %s DIRECTORY COMMAND

Maintain unreleased changelog entries in a modular fashion.

Commands:
    init
    add SECTION STANZA            Add an entry file.
                                  Read from stdin until it
                                  encounters EOF.
    generate [VERSION]            Generate a changelog in
                                  Markdown format and print it
                                  to stdout. VERSION defaults
                                  to UNRELEASED.

    Sections             Stanzas
         ---                 ---
    breaking                gaia
    features             gaiacli
improvements            gaiarest
    bugfixes                 sdk
                      tendermint
`, progName)
	fmt.Fprintf(os.Stderr, "%s\n", usageText)
	//flag.PrintDefaults()
}

func errInsufficientArgs() {
	log.Fatalf("insufficient arguments\nTry '%s -help' for more information.", progName)
}

func errTooManyArgs() {
	log.Fatalf("too many arguments\nTry '%s -help' for more information.", progName)
}

func unknownCommand(cmd string) {
	log.Fatalf("unknown command -- '%s'\nTry '%s -help' for more information.", cmd, progName)
}

// DONTCOVER
