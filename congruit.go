package main

import (
	"./congruit-go/libs"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

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

	places_ptr, works_ptr, workplaces_ptr = congruit.LoadStockroom(*StockRoomDir, *Debug)

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
