package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
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

	files, err = ioutil.ReadDir("./stockroom")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {

			if strings.EqualFold(file.Name(), "places") {

				log.Printf("Loading places...")

				places, _ := ioutil.ReadDir("stockroom/places")
				for _, p := range places {
					log.Printf("Found place: " + p.Name())
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
					log.Printf("Found work: " + w.Name())
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
					log.Printf("Found workplace: " + wp.Name())

					file, _ := os.Open("stockroom/workplaces_enabled/" + wp.Name())
					decoder := json.NewDecoder(file)


					//configuration := WorkplaceConfigurationJson{}

		// read open bracket
   		t, err := decoder.Token()
   		if err != nil {
   			log.Fatal(err)
   	  }

   	  	fmt.Printf("%T: %v\n", t, t)
   	  	cnt := 1
		// while the array contains values
   		for decoder.More() {
				configuration := WorkplaceConfigurationJson{}

   			// decode an array value (Message)
   			err := decoder.Decode(&configuration)
   			if err != nil {
   				log.Fatal(err)
   			}

   					thisworkplace := new(congruit.WorkPlace)
						thisworkplace.Name = wp.Name() + "@"  + strconv.Itoa(cnt)
						thisworkplace.Works = configuration.Works
						thisworkplace.Places = configuration.Places
						log.Printf("Loading workplace: " + wp.Name() + "@"  + strconv.Itoa(cnt))
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

  }else{

		log.Printf("*****************************************")
		log.Printf("*****************************************")
		log.Printf("Going to apply workplaces")
		log.Printf("*****************************************")
		log.Printf("*****************************************")

  }

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

			log.Printf("Executing Place:\n" + command)
			cmd := exec.Command("bash", "-c", command)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				goodplace = false
				break
			}
			log.Printf(out.String())
		}

		command = ""
		if goodplace == true {
			log.Printf("Your place is good! I can do works")


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
			if err != nil {

			}
			log.Printf(out.String())}




		} else {
			log.Printf("Your place is NOT good! I cannot this works for the workplace " + workplace.Name)
		}
	}
}
