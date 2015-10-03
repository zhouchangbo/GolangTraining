package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"fmt"
)

var data = make(map[string]string)

func handle(conn net.Conn) {
	defer conn.Close()

	// NewScanner returns a new Scanner to read from r.
	// The split function defaults to ScanLines.
	scanner := bufio.NewScanner(conn)
	// Scan advances the Scanner to the next token, which will then be
	// available through the Bytes or Text method.
	for scanner.Scan() {
		// Text returns the most recent token generated by a call to Scan
		// as a newly allocated string holding its bytes.
		ln := scanner.Text()
		// Fields splits the string s around each instance of one or more consecutive white space
		// characters, as defined by unicode.IsSpace, returning an array of substrings of s or an
		// empty list if s contains only white space.
		fs := strings.Fields(ln)
		// skip blank lines
		if len(fs) < 2 {
			fmt.Fprintln(conn, "FROM SERVER - USAGE <GET | SET | DEL> <KEY> [VAL]")
			continue
		}

		switch fs[0] {
		case "GET":
			key := fs[1]
			value := data[key]
			fmt.Fprintf(conn, "%s\n", value)
		case "SET":
			if len(fs) != 3 {
				io.WriteString(conn, "EXPECTED VALUE\n")
				continue
			}
			key := fs[1]
			value := fs[2]
			data[key] = value
		case "DEL":
			key := fs[1]
			delete(data, key)
		default:
			fmt.Println(ln)
			ln = fmt.Sprint("FROM SERVER - USAGE <GET | SET | DEL> <KEY> [VAL]\n" +
			"INVALID COMMAND: "+fs[0]+"\n\n")
			io.WriteString(conn, ln)
		}
	}
}

func main() {
	li, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		handle(conn)
	}
}
