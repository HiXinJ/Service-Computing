package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

var progname string

type selpgArgs struct {
	startPage  int
	endPage    int
	inFilename string
	pageLen    int
	pageType   bool //True for form-feed-delimited page type, fasle for lines-delimited page type
	printDest  string
}

var pstartPage = pflag.IntP("startPage", "s", -1, "start page")
var pendPage = pflag.IntP("endPage", "e", -1, "end page")
var ppageLen = pflag.IntP("pageLen", "l", 72, "length of page")
var ppageType = pflag.BoolP("pageType", "f", false, "page type")
var a = pflag.Int("ae", 1, "efa")
var pprintDest = pflag.StringP("printDest", "d", "", "destination to print")

func main() {
	progname = os.Args[0]
	pflag.Parse()

	var sa selpgArgs
	sa.startPage = *pstartPage
	sa.endPage = *pendPage
	sa.pageLen = *ppageLen
	sa.pageType = *ppageType
	sa.printDest = *pprintDest
	ValidateArgs(sa)

	sa.inFilename = pflag.Arg(0)

	processInput(sa)
}

// ValidateArgs 判断传入参数的正确性
func ValidateArgs(sa selpgArgs) {
	if pflag.NArg() > 1 {
		pflag.Usage()
		os.Exit(1)
	}

	if sa.startPage < 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", progname, sa.startPage)
		os.Exit(1)
	}
	if sa.startPage > sa.endPage {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", progname, sa.endPage)
		os.Exit(1)
	}
	if sa.pageLen <= 0 {
		fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", progname, sa.pageLen)
		os.Exit(1)
	}

}

func processInput(sa selpgArgs) {
	var cmd *exec.Cmd
	var fin io.ReadCloser
	var fout io.WriteCloser
	var lineCtr, // line counter
		pageCtr int // page counter
	var err error
	var line []byte
	if sa.inFilename == "" {
		fin = os.Stdin
	} else {

		if fin, err = os.Open(sa.inFilename); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if sa.printDest == "" {
		fout = os.Stdout
	} else {

		cmd = exec.Command("lp", "-d", sa.printDest)
		cmd.Stderr = os.Stderr
		fout, _ = cmd.StdinPipe()
	}

	finReader := bufio.NewReader(fin)
	//lines-delimited page type
	if sa.pageType == false {
		lineCtr = 0
		pageCtr = 1
		for {
			if line, err = finReader.ReadBytes('\n'); err == io.EOF {
				break
			}
			lineCtr++
			if lineCtr > sa.pageLen {
				pageCtr++
				lineCtr = 1
			}
			if pageCtr >= sa.startPage && pageCtr <= sa.endPage {
				fmt.Fprintf(fout, string(line))
			}
		}
	} else {
		pageCtr = 1
		var c byte
		for {

			if c, err = finReader.ReadByte(); err == io.EOF {
				break
			} else if c == '\f' {
				pageCtr++
			}

			if pageCtr >= sa.startPage && pageCtr <= sa.endPage {
				fmt.Fprintf(fout, string(c))
			}
		}
	}
	if sa.printDest != "" {
		cmd.Run()
	}
	if pageCtr < sa.startPage {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d), no output written\n", progname, sa.startPage, pageCtr)
	} else if pageCtr < sa.endPage {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, sa.endPage, pageCtr)
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}
