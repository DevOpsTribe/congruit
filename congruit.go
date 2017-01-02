package main

import (
	"./congruit-go/libs"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type WorkplaceConfigurationJson struct {
	Places []string
	Works  []string
}

func main() {

	const version string = "1.0.0"

	StockRoomDir := flag.String("stockroom-dir", "./", "stockroom directory")
	CongruitVersion := flag.Bool("version", false, "show version")
	Debug := flag.Bool("debug", false, "print more logs")
	ExecutedWorks := 0

	flag.Parse()

	fmt.Printf("                         _ _   \n")
	fmt.Printf(" ___ ___ ___ ___ ___ _ _|_| |_ \n")
	fmt.Printf("|  _| . |   | . |  _| | | |  _|\n")
	fmt.Printf("|___|___|_|_|_  |_| |___|_|_|  \n")
	fmt.Printf("            |___|              \n")
	fmt.Println("Version:", version, "\n")

	if *CongruitVersion {
		return
	}

	works_ptr := []*congruit.Work{}
	workplaces_ptr := []*congruit.WorkPlace{}
	places_ptr := []*congruit.Place{}

	var goodplace bool

	files, err := ioutil.ReadDir(*StockRoomDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {

			if strings.EqualFold(file.Name(), "places") {

				log.Printf("Loading places...")

				places, _ := ioutil.ReadDir("stockroom/places")

				for _, p := range places {

					if *Debug {
						log.Printf("Found place: " + p.Name())
					}
					thisplace := new(congruit.Place)
					thisplace.Name = p.Name()
					content1, _ := ioutil.ReadFile("stockroom/places/" + p.Name())
					thisplace.Command = string(content1)
					places_ptr = append(places_ptr, thisplace)

				}

			} else if strings.EqualFold(file.Name(), "works") {

				log.Printf("Loading works...")

				works, _ := ioutil.ReadDir("stockroom/works")

				for _, w := range works {

					if *Debug {
						log.Printf("Found work: " + w.Name())
					}
					thiswork := new(congruit.Work)
					thiswork.Name = w.Name()
					content2, _ := ioutil.ReadFile("stockroom/works/" + w.Name())
					thiswork.Command = string(content2)
					works_ptr = append(works_ptr, thiswork)

				}

			} else if strings.EqualFold(file.Name(), "workplaces") {

				log.Printf("Loading workplaces...")

				workplaces, _ := ioutil.ReadDir("stockroom/workplaces_enabled")

				for _, wp := range workplaces {

					if *Debug {
						log.Printf("Found workplace: " + wp.Name())
					}

					file, _ := os.Open("stockroom/workplaces_enabled/" + wp.Name())
					decoder := json.NewDecoder(file)

					_, err := decoder.Token()
					if err != nil {
						log.Fatal(err)
					}

					//fmt.Printf("%T: %v\n", t, t)
					cnt := 1

					for decoder.More() {

						configuration := WorkplaceConfigurationJson{}

						err := decoder.Decode(&configuration)
						if err != nil {
							log.Fatal(err)
						}

						thisworkplace := new(congruit.WorkPlace)
						thisworkplace.Name = wp.Name() + "@" + strconv.Itoa(cnt)
						thisworkplace.Works = configuration.Works
						thisworkplace.Places = configuration.Places
						log.Printf("Loading workplace: " + wp.Name() + "@" + strconv.Itoa(cnt))
						workplaces_ptr = append(workplaces_ptr, thisworkplace)
						cnt = cnt + 1
					}
				}
			}
		}
	}

	var command string

	if len(workplaces_ptr) == 0 {

		log.Printf("There are no workplaces to apply... Doing nothing...")

	} else {

		log.Printf("\n *** \n Going to apply workplaces \n *** \n")

	}

	for i := range workplaces_ptr {

		goodplace = true
		workplace := workplaces_ptr[i]
		log.Printf("Workplace: " + workplace.Name)

		log.Printf("Checking places...")

		for k := range workplace.Places {

			z := workplace.Places[k]
			log.Printf("Testing Place: " + z)

			for p := range places_ptr {

				place := places_ptr[p]

				if strings.EqualFold(z, place.Name) {
					command = place.Command
				}

			}

			if *Debug {
				log.Printf("Executing Place:\n" + command)
			}

			cmd := exec.Command("bash", "-c", command)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()

			if err != nil {
				goodplace = false
				break
			}

			if *Debug {
				log.Printf("Place execution output: " + out.String())
			}
		}

		command = ""

		if goodplace == true {

			for k := range workplace.Works {
				j := workplace.Works[k]

				for w := range works_ptr {

					work := works_ptr[w]

					if strings.EqualFold(j, work.Name) {
						command = work.Command
					}
				}

				log.Printf("Executing Work:\n" + command)

				cmd := exec.Command("bash", "-c", command)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()

				ExecutedWorks = ExecutedWorks + 1

				if err != nil {

				}
				log.Printf(out.String())
			}

		} else {

			if *Debug {
				log.Printf("Workplace " + workplace.Name + " not needed here!")
			}

		}
	}

	log.Printf("Extecuted works: " + strconv.Itoa(ExecutedWorks))

}
