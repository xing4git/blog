package main

import (
	"fmt"
	"os"
	// "io"
	"bufio"
	"bytes"
	"github.com/xing4git/cmdutils"
	"io/ioutil"
	"regexp"
	"sort"
)

var (
	readme = bytes.NewBuffer(make([]byte, 0))
	change = bytes.NewBuffer(make([]byte, 0))
	commit = bytes.NewBuffer(make([]byte, 0))
)

var (
	reg = regexp.MustCompile("^_\\d{4}-\\d{2}-\\d{2}-")
	filenames = make([]string, 0)
)

func main() {
	dir, err := os.Open(".")
	checkErr(err)
	defer dir.Close()

	fis, err := dir.Readdir(0)
	checkErr(err)

	for i := 0; i < len(fis); i++ {
		fi := fis[i]
		filename := fi.Name()
		if reg.Match([]byte(filename)) {
			filenames = append(filenames, filename)
		}
	}

	sort.Strings(filenames)
	fmt.Println(filenames)

	readme.WriteString("Index\n")
	readme.WriteString("-----\n\n")
	change.WriteString("update README.md;")

	for key, value := range filenames {
		realname := value[12:]
		readme.WriteString("####" + decorateFilename(realname) + "\n")

		file, err := os.Open(value)
		checkErr(err)
		defer file.Close()

		readme.WriteString(value[1:11] + "\n")
		buf := bufio.NewReader(file)
		line, err := buf.ReadString('\n')
		checkErr(err)
		readme.WriteString(string(line[0:len(line)-1]))
		line, err = buf.ReadString('\n')
		checkErr(err)
		readme.WriteString(string(line[0:len(line)-1]))
		readme.WriteString("...[Read More](golang/" + realname + ")\n\n")

		file.Seek(0, os.SEEK_SET)
		nbytes, err := ioutil.ReadAll(file)
		checkErr(err)
		var tempbuf []byte
		tempbuf = append(tempbuf, []byte(decorateFilename(realname) + "\n" + "----\n\n")...)
		nbytes = append(tempbuf, nbytes...)

		nbytes = append(nbytes, []byte("\n\n" + "links\n" + "-----\n")...)
		if key != 0 {
			previous := filenames[key-1][12:]
			nbytes = append(nbytes, []byte("+ previous: [" + decorateFilename(previous) + "](" + previous + ")\n")...)
		}
		if key != len(filenames)-1 {
			next := filenames[key+1][12:]
			nbytes = append(nbytes, []byte("+ next: [" + decorateFilename(next) + "](" + next + ")\n")...)
		}

		pbytes, err := ioutil.ReadFile(realname)
		if err == nil && compareBytes(nbytes, pbytes) == 0 {
			fmt.Println("no change file:", realname)
			continue
		}

		change.WriteString("update " + realname + ";")
		err = ioutil.WriteFile(realname, nbytes, 0664)
		checkErr(err)
	}

	err = ioutil.WriteFile("README.md", readme.Bytes(), 0664)
	checkErr(err)

	commit.WriteString("git add .\n")
	commit.WriteString("git commit -a -m '" + change.String() + "'\n")
	commit.WriteString("git push -u origin master\n")
	fmt.Println(commit.String())
	ret, err := cmdutils.BashExecute(commit.String())
	checkErr(err)
	fmt.Println(ret)
}

func decorateFilename(str string) string {
	bytes := []byte(str)
	var suffixStart int = -1
	for i := len(str)-1; i >= 0; i-- {
		if bytes[i] == '.' {
			suffixStart = i
			break
		}
	}

	if suffixStart != -1 {
		bytes = bytes[:suffixStart]
	}

	for pos, v := range bytes {
		if v == '-' {
			bytes[pos] = ' '
		}
	}

	return string(bytes)
}

func compareBytes(pbytes []byte, nbytes []byte) int {
	for i := 0; i < len(pbytes) && i < len(nbytes); i++ {
		if pbytes[i] > nbytes[i] {
			return 1
		} else if pbytes[i] < nbytes[i] {
			return -1
		}
	}
	if len(pbytes) > len(nbytes) {
		return 1
	} else if len(pbytes) < len(nbytes) {
		return -1
	}
	return 0
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
