package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func getmsg(w http.ResponseWriter, r *http.Request) {
	voices := r.FormValue("voices")
	text := r.FormValue("text")

	voice := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(voices))
	for scanner.Scan() {
		str := scanner.Text()
		parts := strings.Split(str, ":")
		if len(parts) > 1 {
			voice[parts[0]] = parts[1]
		}
	}

	scanner = bufio.NewScanner(strings.NewReader(text))
	i := 0
	currentvoice := ""
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 1 {
			end := line[len(line)-1:]
			if end == ":" {
				log.Println("Changing voice to " + line + voice[line[:len(line)-1]])
				currentvoice = voice[line[:len(line)-1]]
				// pause := exec.Command("sh", "-c", "cp pause.mp3 "+fmt.Sprintf("file_%06d.mp3 ", i))
				// _, err := pause.Output()
				// if err != nil {
				// 	panic(err)
				// }
			} else {
				line = strings.TrimSpace(line)
				if line != "" {
					// text := "<speak>" + line + "</speak>"
					fileext := fmt.Sprintf("file_%06d.mp3", i)
					// if _, err := os.Stat(fileext); os.IsNotExist(err) {
					// args := "aws polly synthesize-speech --text-type ssml --text " + strconv.Quote(text) + " --output-format mp3 --voice-id " + currentvoice + " " + fileext
					// log.Println(args)

					// lsCmd := exec.Command("sh", "-c", args)
					// _, err := lsCmd.Output()
					res := makeSpeech(line, currentvoice)
					file, err := os.Create(fileext)
					if err != nil {
						log.Fatal(err)
					}

					_, err = io.Copy(file, res)
					if err != nil {
						log.Fatal(err)
					}
					file.Close()
					// }
					i = i + 1
				}
			}
		}
	}
	rm := exec.Command("sh", "-c", "rm result.mp3 ara_MP3WRAP.mp3")
	_, err := rm.Output()
	if err != nil {
		log.Println(err)
	}
	// mp3wrap result.mp3 file_000000.mp3 file_000001.mp3 file
	// ffmpeg -i "concat:file1.mp3|file2.mp3" -acodec copy output.mp3
	// cmd := "cat "
	cmd := "mp3wrap ara "
	for j := 0; j < i; j++ {
		cmd += fmt.Sprintf("file_%06d.mp3 ", j)
	}
	// cmd += " result.mp3"
	log.Println(cmd)
	catCmd := exec.Command("sh", "-c", cmd)
	_, err = catCmd.Output()
	if err != nil {
		log.Panicf("catcmd: %v \n", err)
	} else {
		pause := exec.Command("sh", "-c", "mv ara_MP3WRAP.mp3 result.mp3")
		_, err = pause.Output()
		if err != nil {
			log.Panicf("cp wrap  %v \n", err)
		}
		rmCmd := exec.Command("sh", "-c", "rm file*")
		_, err = rmCmd.Output()
		if err != nil {
			log.Panicf("rm file %v \n ", err)
		}
	}

	file, err := os.Open("result.mp3")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Disposition", "attachment; filename=result.mp3")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, file)
	defer file.Close()
}

func main() {
	http.HandleFunc("/polly", getmsg)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	http.ListenAndServe(":8989", nil)
}
