package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	"flag"
	"fmt"
	"github.com/perriv/go-tasker"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var version = "0.3.0"

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
// Read the JSON configuration provided. If c is an empty string, just return
// the default configuration (the list of visible directories in the current
// working directory).
func read_config(c string) ([]string, error) {
	if c == "" {
		return list_visible_dirs(".")
	}

	data, err := ioutil.ReadFile(c)
	if err != nil {
		return nil, err
	}

	config := make([]string, 0)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func cmd_task(cmd *exec.Cmd, l *logger, data map[string][]byte) tasker.Task {
	return func() error {
		prefix := fmt.Sprintf("%s: %s", cmd.Dir, strings.Join(cmd.Args, " "))
		l.ok("%s: started\n", prefix)

		dat, err := cmd.CombinedOutput()

		if err == nil {
			l.good("%s: finished\n", prefix)
		} else {
			l.bad("%s: failed: %s\n", prefix, err)
		}

		data[cmd.Dir] = dat
		return err
	}
}

func main() {
	var j int
	var v bool
	var c string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [--] <command>...\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.StringVar(&c, "c", "", "Configuration file")
	flag.IntVar(&j, "j", -1, "Number of concurrent subprocesses")
	flag.BoolVar(&v, "v", false, "Print out the version and exit")
	flag.Parse()

	if v {
		fmt.Println(version)
		os.Exit(0)
	}

	// 1 or more arguments are required.
	cmd_format := flag.Args()
	if len(cmd_format) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	config, err := read_config(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	if len(config) == 0 {
		fmt.Println("Nothing to do")
		os.Exit(0)
	}
	sort.Strings(config)

	// Create the Tasker that runs all of the commands.
	tr, err := tasker.NewTasker(j)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// This logger coordinates printed among the tasks.
	l := new_logger()
	data := make(map[string][]byte)

	for _, dir := range config {

		// A command may include {} to interpolate the current directory.
		name := strings.Replace(cmd_format[0], "{}", dir, -1)
		args := make([]string, 0)
		for _, arg_f := range cmd_format[1:] {
			args = append(args, strings.Replace(arg_f, "{}", dir, -1))
		}

		// Add a task that runs the interpolated command in the current directory.
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		tr.Add(dir, nil, cmd_task(cmd, l, data))
	}

	err = tr.Run()

	for _, dir := range config {
		if len(data[dir]) == 0 {
			continue
		}
		lines := bytes.Split(data[dir], []byte{10})
		for _, line := range lines {
			os.Stdout.WriteString(dir)
			os.Stdout.WriteString(": ")
			os.Stdout.Write(line)
			os.Stdout.WriteString("\n")
		}
	}

	if err != nil {
		os.Exit(-1)
	}
}
