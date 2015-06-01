package main

import (
	"github.com/ziutek/dvb/ts"
	"github.com/ziutek/dvb/ts/psi"
	"io"
	"log"
	"os"
	"strconv"
)

const usage = `Usage: ss2ts PID
ss2ts encapsulates section stream from stdin into transport stream on stdout.`

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Fatal(usage)
	}

	pid, err := strconv.ParseUint(os.Args[1], 0, 64)
	checkErr(err)
	if uint64(pid) > 8191 {
		log.Fatal(pid, "isn't valid PID")
	}

	s := make(psi.Section, psi.SectionMaxLen)
	r := psi.NewSectionStreamReader(os.Stdin, pid != 20)
	w := psi.NewSectionEncoder(ts.PktStreamWriter{os.Stdout}, int16(pid))
	for {
		err := r.ReadSection(s)
		if err != nil {
			if err == io.EOF {
				break
			}
			checkErr(err)
		}
		/*log.Println(
			"TableId:", s.TableId(),
			"TableIdExt:", s.TableIdExt(),
			"GenericSyntax:", s.GenericSyntax(),
			"PrivateSyntax:", s.PrivateSyntax(),
			"Len:", s.Len(),
			"Version:", s.Version(),
			"Current:", s.Current(),
			"Number:", s.Number(),
			"LastNumber:", s.LastNumber(),
		)*/
		checkErr(w.WriteSection(s))
	}
	checkErr(w.Flush())
}
