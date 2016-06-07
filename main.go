package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	LF = 10

	DATA    = "DATA"
	GROUP   = "GROUP"
	HEADING = "HEADING"
	TYPE    = "TYPE"
	UNIT    = "UNIT"
)

// AGS order
type Node struct {
	Group   string     `json:"GROUP"`
	Heading []string   `json:"HEADING"`
	Unit    []string   `json:"UNIT"`
	Type    []string   `json:"TYPE"`
	Data    [][]string `json:"DATA"`
}

var Tabs = make(map[string][]Node)

func main() {

	// func rude(cache string)
	rude(os.Args[1])

	// func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	bytes, _ := json.MarshalIndent(Tabs, "", "  ")
	fmt.Printf("%s", bytes)
}

/*
	// alternate to pretty print (above)
	for name, nodes := range Tabs {
		fmt.Println(name)
		for _, node := range nodes {
			fmt.Println("GROUP " + node.Group)
			fmt.Print("HEADING ")
			fmt.Println(node.Heading)
			fmt.Print("TYPE ")
			fmt.Println(node.Type)
			fmt.Print("UNIT ")
			fmt.Println(node.Unit)
			fmt.Print("DATA ")
			fmt.Println(node.Data)
			fmt.Println()
		}
	}
*/

/*
Function rude is a rudimentary parser for AGS standard data. It takes as an argument the name of
the directory containing the AGS definition files.
*/
func rude(cache string) {

	// func ReadDir(dirname string) ([]os.FileInfo, error)
	files, err := ioutil.ReadDir(cache)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		name := file.Name()
		//fmt.Println(name)

		// func ReadFile(filename string) ([]byte, error)
		bloc, err := ioutil.ReadFile(cache + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		// func NewBuffer(buf []byte) *Buffer
		buf := bytes.NewBuffer(bloc)
		node := new(Node)

		for {
			// func (b *Buffer) ReadString(delim byte) (line string, err error)
			line, err := buf.ReadString(LF)
			if err != nil {
				break
			}
			// func Trim(s string, cutset string) string
			line = strings.Trim(line, "\r\n")
			// AGS format terminated with redundant \r\n
			if len(line) == 0 {
				if Tabs[name] == nil {
					Tabs[name] = make([]Node, 0)
				}
				Tabs[name] = append(Tabs[name], *node)
				node = new(Node)
				continue
			}

			// used only during DATA cycle
			data := []string{}

			head := true
			mode := ""

			// func Split(s, sep string) []string
			for _, token := range strings.Split(line, ",") {
				// strip quotes
				token = strings.Trim(token, "\"")
				// set new mode
				if head {
					mode = token
					head = false
					continue
				}

				switch mode {
				case DATA:
					data = append(data, token)
				case GROUP:
					node.Group = token
				case HEADING:
					node.Heading = append(node.Heading, token)
				case TYPE:
					node.Type = append(node.Type, token)
				case UNIT:
					node.Unit = append(node.Unit, token)
				}
			}

			if mode == DATA {
				node.Data = append(node.Data, data)
			}
		}
	}
}
