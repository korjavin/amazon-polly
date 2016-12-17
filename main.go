package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	vflag string
)

func init() {
	flag.StringVar(&vflag, "v", "", "do we need log?")
	flag.Parse()
	if vflag == "" {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	voice := voices("voice.txt")
	file, err := os.Open("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	currentvoice := ""
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 1 {
			end := line[len(line)-1:]
			if end == ":" {
				log.Println("Changing voice to " + line + voice[line[:len(line)-1]])
				currentvoice = voice[line[:len(line)-1]]
			} else {
				text := "<speak>" + line + "</speak>"
				fileext := fmt.Sprintf("file_%06d.mp3", i)
				args := "aws polly synthesize-speech --text-type ssml --text " + strconv.Quote(text) + " --output-format mp3 --voice-id " + currentvoice + " " + fileext
				log.Println(args)

				lsCmd := exec.Command("sh", "-c", args)
				_, err := lsCmd.Output()
				if err != nil {
					panic(err)
				}
				i = i + 1
			}
		}
	}
}

func voices(filename string) map[string]string {
	m1 := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := scanner.Text()
		parts := strings.Split(str, ":")
		if len(parts) > 1 {
			m1[parts[0]] = parts[1]
		}
	}
	return m1

}
