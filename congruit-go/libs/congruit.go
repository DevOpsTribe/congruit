package congruit

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

func ExecuteStockroom(Debug bool, places []Place, works []Work, workplaces []WorkPlace) int {

	var goodplace bool

	var command string

	ExecutedWorks := 0

	for i := range workplaces {

		goodplace = true
		workplace := workplaces[i]
		log.Printf("Processing workplace " + workplace.Name)

		log.Printf("Checking places of workplace " + workplace.Name)

		for k := range workplace.Places {

			z := workplace.Places[k]
			log.Printf("Testing place " + z)
			place_name := ""

			for p := range places {

				place := places[p]

				if strings.EqualFold(z, place.Name) {
					command = place.Command
					place_name = place.Name
				}

			}

			if len(command) == 0 {
				log.Fatal("Error in loading places!")
			}

			cmd := exec.Command("bash", "-c", command)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()

			if err != nil {
				goodplace = false
				log.Printf("Place " + place_name + " does not return 0. This is not a good place to run " + workplace.Name)
				fmt.Println(err)
				break
			}

			log.Printf("Place " + place_name + " returns 0")

			if Debug {
				//log.Printf("\n\n output of \n" + command + "\n is " + out.String() + "\n\n")
			}
		}

		command = ""

		if goodplace == true {

			for k := range workplace.Works {
				j := workplace.Works[k]

				for w := range works {

					work := works[w]

					if strings.EqualFold(j, work.Name) {
						command = work.Command
						if len(command) == 0 {
							log.Fatal("Error in loading places!")
						}
					}
				}

				log.Printf("Executing work: \n" + command + "\n\n")

				cmd := exec.Command("bash", "-c", command)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()

				ExecutedWorks = ExecutedWorks + 1

				if err != nil {

				}
				if len(out.String()) > 0 {
					log.Printf("command output: " + out.String() + "\n")
				}
			}

		} else {

			if Debug {
				log.Printf("Workplace " + workplace.Name + " not needed here")
			}

		}
	}
	return ExecutedWorks
}

func LoadStockroom(StockRoomDir string, Debug bool) ([]Place, []Work, []WorkPlace) {

	log.Printf("Loading works...")

	places := make([]Place, 0)
        works := make([]Work, 0)
        workplaces := make([]WorkPlace, 0)

	raw_works, _ := ioutil.ReadDir(StockRoomDir + "/works")
	raw_workplaces, _ := ioutil.ReadDir(StockRoomDir + "/workplaces_enabled")
	raw_places, _ := ioutil.ReadDir(StockRoomDir + "/places")

	for _, w := range raw_works {

		if Debug {
			log.Printf("Found work: " + w.Name())
		}

		content, _ := ioutil.ReadFile(StockRoomDir + "/works/" + w.Name())
		thiswork := Work{Name: w.Name(), Command: string(content)}
		if Debug {
			log.Printf("Load work: " + thiswork.Name)
		}
		works = append(works, thiswork)
	}

	for _, p := range raw_places {

		if Debug {
			log.Printf("Found place: " + p.Name())
		}
		content, _ := ioutil.ReadFile(StockRoomDir + "/places/" + p.Name())
		thisplace := Place{Name: p.Name(), Command: string(content)}
		if Debug {
			log.Printf("Load Place: " + thisplace.Name)
		}
		places = append(places, thisplace)

	}

	for _, wp := range raw_workplaces {

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

				thisworkplace := WorkPlace{Name: wp.Name() + "@" + strconv.Itoa(cnt), Works: configuration.Works,
                                   Places: configuration.Places}
				log.Printf("Loading workplace: " + wp.Name() + "@" + strconv.Itoa(cnt))
				workplaces = append(workplaces, thisworkplace)
				cnt = cnt + 1
			}
		}
	}

	return places, works, workplaces

}

func CloneRepo (GitRepo string) string {

  StockRoomDir := ""

	if _, err := os.Stat("/tmp/stockroom"); err == nil {
		os.RemoveAll("/tmp/stockroom")
	}

	cmd := exec.Command("git", "clone", GitRepo, "/tmp/stockroom")

	StockRoomDir = "/tmp/stockroom"

	err := cmd.Run()

	if err != nil {

		log.Fatal("Error when pull stockroom repo")

	}

	StockRoomDir = "/tmp/stockroom"
  return StockRoomDir
}
