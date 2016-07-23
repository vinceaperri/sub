package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/perriv/go-tasker"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
)

var version = "0.1.0"

func is_visible_dir(fi os.FileInfo) bool {
	return fi.Mode().IsDir() && !strings.HasPrefix(fi.Name(), ".")
}

func list_visible_dirs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0)
	for _, fi := range fis {
		if is_visible_dir(fi) {
			dirs = append(dirs, fi.Name())
		}
	}
	dirs = dirs[:len(dirs)]
	sort.Strings(dirs)
	return dirs, nil
}

// Read a file line-by-line with \n (linux) line endings. Returns a list of
// lines with line endings removed.
func read_lines(path string) ([]string, error) {
	nl := byte(10)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)

	// Read file line-by-line until EOF is reached. Ignore empty lines.
	dirs := make([]string, 0)
	for {
		line, err := r.ReadString(nl)

		// Strip the line ending.
		if len(line) > 0 && line[len(line)-1] == nl {
			line = line[:len(line)-1]
		}

		// Ignore empty lines.
		if len(line) > 0 {
			dirs = append(dirs, line)
		}

		// Done reading file.
		if err == io.EOF {
			break
		}

		// Unexpected error.
		if err != nil {
			return nil, err
		}

	}
	return dirs, nil
}

// -d flag value, which accumulates multiple instances -d. It implements the
// flag.Value interface.
type directory_flag_value []string

func (d *directory_flag_value) String() string {
	return fmt.Sprint(*d)
}

func (d *directory_flag_value) Set(value string) error {
	*d = append(*d, value)
	return nil
}

// Print a status message in bold blue to stdout
func print_status(message string) {
	fmt.Println("\x1B[1;34m" + message + "\x1B[0m")
}

// Print an error message in bold red to stderr
func print_error(message string) {
	fmt.Fprintln(os.Stderr, "\x1B[1;31m"+message+"\x1B[0m")
}

// Run name as a subprocess with args as args inside each dir in dirs in
// parallel.
// j is the number of simultaneous subprocesses.
func run(name string, args []string, dirs []string, j int) error {
	mux := &sync.Mutex{}

	tr, err := tasker.NewTasker(j)
	if err != nil {
		fmt.Println(err)
	}
	for _, dir := range dirs {
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		tr.Add(dir, nil, func() error {
			out, err := cmd.CombinedOutput()

			// Write command output one at a time.
			mux.Lock()
			defer mux.Unlock()

			if err == nil {
				// Write done message in bold blue.
				print_status(cmd.Dir + ": done")
			} else {
				// Write failed message in bold red.
				print_error(cmd.Dir + ": failed: " + err.Error())
			}
			if _, oe := os.Stdout.Write(out); oe != nil {
				print_error(cmd.Dir + ": failed to write output: " + err.Error())
			}
			return err
		})
	}
	return tr.Run()
}

func main() {
	var dirs directory_flag_value
	var j int
	var v bool

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [--] <command>...\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Var(&dirs, "d", "Run only in this directory")
	flag.IntVar(&j, "j", 4, "Number of concurrent subprocesses")
	flag.BoolVar(&v, "v", false, "Print out the version and exit")
	flag.Parse()

	if v {
		fmt.Println(version)
		os.Exit(0)
	}

	// 1 or more arguments are required.
	cmd := flag.Args()
	if len(cmd) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// If -d is not given on the command-line, fall-back to sub.cnf.
	// If sub.cnf can't be read, then we just run in all visible directories.
	if len(dirs) == 0 {
		cnf_dirs, err := read_lines("sub.cnf")
		if err == nil {
			dirs = cnf_dirs
		} else {
			vis_dirs, err := list_visible_dirs(".")
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			dirs = vis_dirs
		}
		if len(dirs) == 0 {
			os.Exit(0)
		}
	}
	sort.Strings(dirs)

	if err := run(cmd[0], cmd[1:], dirs, j); err != nil {
		os.Exit(-1)
	}
}
