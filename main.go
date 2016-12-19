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
	kflag string
	fflag string
)

func init() {
	flag.StringVar(&vflag, "v", "", "do we need log?")
	flag.StringVar(&kflag, "k", "", "keep files?")
	flag.StringVar(&fflag, "f", "text.txt", "file to read")
	flag.Parse()
	if vflag == "" {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	voice := voices("voice.txt")
	file, err := os.Open(fflag)
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
				pause := exec.Command("sh", "-c", "cp pause.mp3 "+fmt.Sprintf("file_%06d.mp3 ", i))
				_, err := pause.Output()
				if err != nil {
					panic(err)
				}
			} else {
				line = strings.Replace(line, "[", "<emphasis level='strong'>", -1)
				line = strings.Replace(line, "]", "</emphasis>", -1)
				text := "<speak>" + line + "</speak>"
				fileext := fmt.Sprintf("file_%06d.mp3", i)
				if _, err := os.Stat(fileext); os.IsNotExist(err) {
					args := "aws polly synthesize-speech  --lexicon-names=\"lexicon\" --text-type ssml --text " + strconv.Quote(text) + " --output-format mp3 --voice-id " + currentvoice + " " + fileext
					log.Println(args)

					lsCmd := exec.Command("sh", "-c", args)
					_, err := lsCmd.Output()
					if err != nil {
						panic(err)
					}
				}
			}
			i = i + 1
		}
	}
	rm := exec.Command("sh", "-c", "rm result.mp3")
	_, err = rm.Output()
	if err != nil {
		log.Println(err)
	}
	cmd := "cat "
	for j := 0; j < i; j++ {
		cmd += fmt.Sprintf("file_%06d.mp3 ", j)
	}
	cmd += " > result.mp3"
	catCmd := exec.Command("sh", "-c", cmd)
	_, err = catCmd.Output()
	if err != nil {
		log.Println(err)
	} else {
		if kflag == "" {
			rmCmd := exec.Command("sh", "-c", "rm file*")
			_, err = rmCmd.Output()
			if err != nil {
				log.Println(err)
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

// aws polly put-lexicon \
// --name w3c \
// --content file://example.pls
