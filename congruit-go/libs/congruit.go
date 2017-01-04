package congruit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
		"bytes"
	"os/exec"
	"strings"
)

type WorkplaceConfigurationJson struct {
	Places []string
	Works  []string
}

type Work struct {
	Name    string
	Command string
}

type Place struct {
	Name    string
	Command string
}

type WorkPlace struct {
	Name   string
	Works  []string
	Places []string
}


func ExecuteStockroom(Debug bool,places_ptr []*Place, works_ptr []*Work, workplaces_ptr []*WorkPlace) (int){

	var goodplace bool

	var command string

	ExecutedWorks := 0


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
						log.Printf("Command is \n" + place.Command )
						command = place.Command
					}

				}

				if len(command) == 0 {
					log.Fatal("Error in loading places!")
				}

				if Debug {
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

				if Debug {
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
							if len(command) == 0 {
								log.Fatal("Error in loading places!")
							}
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

				if Debug {
					log.Printf("Workplace " + workplace.Name + " not needed here!")
				}

			}
		}
	return ExecutedWorks
}

func LoadStockroom(StockRoomDir string, Debug bool) ([]*Place, []*Work, []*WorkPlace) {

	log.Printf("Loading works...")

	places_ptr := []*Place{}
	works_ptr := []*Work{}
	workplaces_ptr := []*WorkPlace{}

	works, _ := ioutil.ReadDir(StockRoomDir + "/works")
	workplaces, _ := ioutil.ReadDir(StockRoomDir + "/workplaces_enabled")
	places, _ := ioutil.ReadDir(StockRoomDir + "/places")

	for _, w := range works {

		if Debug {
			log.Printf("Found work: " + w.Name())
		}

		thiswork := new(Work)
		thiswork.Name = w.Name()
		content, _ := ioutil.ReadFile(StockRoomDir + "/works/" + w.Name())
		thiswork.Command = string(content)
		if Debug {
			log.Printf("Load work: " + thiswork.Name)
		}
		works_ptr = append(works_ptr, thiswork)
	}

	for _, p := range places {

		if Debug {
			log.Printf("Found place: " + p.Name())
		}
		thisplace := new(Place)
		thisplace.Name = p.Name()
		content, _ := ioutil.ReadFile(StockRoomDir + "/places/" + p.Name())
		thisplace.Command = string(content)
		if Debug {
			log.Printf("Load Place: " + thisplace.Name)
		}
		places_ptr = append(places_ptr, thisplace)

	}

	for _, wp := range workplaces {

		if strings.EqualFold(wp.Name(), "README.md") == false {
			if Debug {
				log.Printf("Found workplace: " + wp.Name())
			}

			file, _ := os.Open(StockRoomDir + "/workplaces_enabled/" + wp.Name())
			decoder := json.NewDecoder(file)

			_, err := decoder.Token()
			if err != nil {
				log.Fatal(err)
			}

			cnt := 1

			for decoder.More() {

				configuration := WorkplaceConfigurationJson{}

				err := decoder.Decode(&configuration)
				if err != nil {
					log.Fatal(err)
				}

				thisworkplace := new(WorkPlace)
				thisworkplace.Name = wp.Name() + "@" + strconv.Itoa(cnt)
				thisworkplace.Works = configuration.Works
				thisworkplace.Places = configuration.Places
				log.Printf("Loading workplace: " + wp.Name() + "@" + strconv.Itoa(cnt))
				workplaces_ptr = append(workplaces_ptr, thisworkplace)
				cnt = cnt + 1
			}
		}
	}

	return places_ptr, works_ptr, workplaces_ptr

}
