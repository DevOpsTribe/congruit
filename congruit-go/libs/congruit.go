package congruit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
