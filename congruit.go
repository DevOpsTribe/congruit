package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"./congruit-go/libs"
)

type WorkplaceConfigurationJson struct {
	Places []string
	Works  []string
}

func main() {

	log.Printf("*****************************************")
	log.Printf("*****************************************")
	log.Printf("Thank you for use Congruit!							 ")
	log.Printf("*****************************************")
	log.Printf("*****************************************")

	works_ptr := []*congruit.Work{}
	workplaces_ptr := []*congruit.WorkPlace{}
	places_ptr := []*congruit.Place{}

	files, err := ioutil.ReadDir("/tmp/congruit_temp")

	var goodplace bool

	if err != nil {
		os.Mkdir("/tmp/congruit_temp", 700)
	}

	files, err = ioutil.ReadDir("./legion")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			//log.Printf("Found directory: " + file.Name())

			if strings.EqualFold(file.Name(), "places") {

				log.Printf("Loading places...")

				places, _ := ioutil.ReadDir("legion/places")
				for _, p := range places {
					log.Printf("Found place: " + p.Name())
					thisplace := new(congruit.Place)
					thisplace.Name = p.Name()
					content1, _ := ioutil.ReadFile("legion/places/" + p.Name())
					thisplace.Command = string(content1)
					places_ptr = append(places_ptr, thisplace)

				}

			} else if strings.EqualFold(file.Name(), "works") {

				log.Printf("Loading works...")

				works, _ := ioutil.ReadDir("legion/works")
				for _, w := range works {
					log.Printf("Found work: " + w.Name())
					thiswork := new(congruit.Work)
					thiswork.Name = w.Name()
					content2, _ := ioutil.ReadFile("legion/works/" + w.Name())
					thiswork.Command = string(content2)
					works_ptr = append(works_ptr, thiswork)
				}

			} else if strings.EqualFold(file.Name(), "workplaces") {

				log.Printf("Loading workplaces...")

				workplaces, _ := ioutil.ReadDir("legion/workplaces")
				for _, wp := range workplaces {
					log.Printf("Found workplace: " + wp.Name())

					file, _ := os.Open("legion/workplaces/" + wp.Name())
					decoder := json.NewDecoder(file)
					configuration := WorkplaceConfigurationJson{}
					err := decoder.Decode(&configuration)
					if err != nil {
						fmt.Println("error:", err)
					}

					thisworkplace := new(congruit.WorkPlace)
					thisworkplace.Name = wp.Name()
					thisworkplace.Works = configuration.Works
					thisworkplace.Places = configuration.Places
					workplaces_ptr = append(workplaces_ptr, thisworkplace)
				}
			}
		}
	}

	var command string
	log.Printf("*****************************************")
	log.Printf("*****************************************")
	log.Printf("Going to apply workplaces")
	log.Printf("*****************************************")
	log.Printf("*****************************************")

	for i := range workplaces_ptr {
		goodplace = true
		workplace := workplaces_ptr[i]
		log.Printf("Workplace: " + workplace.Name)
		log.Printf("Checking all places")
		for k := range workplace.Places {
			z := workplace.Places[k]
			log.Printf("Place: " + z)

			for p := range places_ptr {
				place := places_ptr[p]
				if strings.EqualFold(z, place.Name) {
					command = place.Command
				}
			}

			log.Printf(command)
			cmd := exec.Command("bash", "-c", command)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				goodplace = false
			}
			log.Printf(out.String())
		}

		command = ""
		if goodplace == true {
			log.Printf("Your place is good! I can do works")
		} else {
			log.Printf("Your place is NOT good! I cannot this works for the workplace " + workplace.Name)
			break
		}

		for k := range workplace.Works {
			j := workplace.Works[k]
			for w := range works_ptr {
				work := works_ptr[w]
				if strings.EqualFold(j, work.Name) {
					command = work.Command
				}
			}
			log.Printf(" foo " + command)
			cmd := exec.Command("bash", "-c", command)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {

			}
			log.Printf(out.String())
		}
	}
}
