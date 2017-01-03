package main

import (
	"./congruit-go/libs"
	"flag"
	"fmt"
	"log"
	"strconv"
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

	places_ptr, works_ptr, workplaces_ptr = congruit.LoadStockroom(*StockRoomDir, *Debug)

	if len(workplaces_ptr) == 0 {

		log.Printf("There are no workplaces to apply... Doing nothing...")

	} else {

		log.Printf("\n *** \n Going to apply workplaces \n *** \n")

	}

	ExecutedWorks = congruit.ExecuteStockroom(*Debug, places_ptr, works_ptr, workplaces_ptr)

	log.Printf("Extecuted works: " + strconv.Itoa(ExecutedWorks))

}
