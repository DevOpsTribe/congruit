package main

import (
	"./congruit-go/libs"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request, mystr string) {
	println(mystr)
}

func HelloServer(w http.ResponseWriter, req *http.Request, t string, StockRoomDir string, Debug bool) int {

	ExecutedWorks := 0

	w.Header().Set("Content-Type", "text/plain")

	if strings.EqualFold(t, (req.Header.Get("Token"))) {
		w.Write([]byte("Hi dude!\n"))

		works_ptr := []*congruit.Work{}
		workplaces_ptr_temp := []*congruit.WorkPlace{}
		workplaces_ptr := []*congruit.WorkPlace{}
		places_ptr := []*congruit.Place{}

		places_ptr, works_ptr, workplaces_ptr_temp = congruit.LoadStockroom(StockRoomDir, Debug)

		log.Println("A remote agent asks to run |" + req.Header.Get("Workplace") + "|")

		for i := range workplaces_ptr_temp {
			log.Println(req.Header.Get("Workplace") + "@" + strconv.Itoa(i+1))
			workplace := workplaces_ptr_temp[i]
			if strings.EqualFold(req.Header.Get("Workplace")+"@"+strconv.Itoa(i+1), workplace.Name) {
				workplaces_ptr = append(workplaces_ptr, workplaces_ptr_temp[i])
				log.Println("A remote command starts worksplace " + req.Header.Get("Workspace"))
			}
		}

		ExecutedWorks = congruit.ExecuteStockroom(Debug, places_ptr, works_ptr, workplaces_ptr)
		return ExecutedWorks

	} else {
		w.Write([]byte("Error. Not authorized\n"))
		return 0
	}
}

func main() {

	const version string = "1.0.0"

	StockRoomDir := flag.String("stockroom-dir", "./", "stockroom directory")
	CongruitVersion := flag.Bool("version", false, "show version")
	Debug := flag.Bool("debug", false, "print more logs")
	Supervisor := flag.Bool("supervisor", false, "enable supervisor mode")
	Friend := flag.Bool("friend", false, "enable friend mode")
	Token := flag.String("token", "nil", "token for talking with friends")

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

	if len(workplaces_ptr) == 0 {

		log.Printf("There are no workplaces to apply... Doing nothing...")

	} else {

		log.Printf("\n *** \n Going to apply workplaces \n *** \n")

	}

	log.Printf("Extecuted works: " + strconv.Itoa(ExecutedWorks))

	if *Friend {

		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			HelloServer(w, r, *Token, *StockRoomDir, *Debug)
		})

		err := http.ListenAndServeTLS(":443", "domain.crt", "domain.key", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}

	places_ptr, works_ptr, workplaces_ptr = congruit.LoadStockroom(*StockRoomDir, *Debug)

	ExecutedWorks = congruit.ExecuteStockroom(*Debug, places_ptr, works_ptr, workplaces_ptr)

	for *Supervisor {

		log.Printf("start congruit supervisor...")

		time.Sleep(2000 * time.Millisecond)

		ExecutedWorks = congruit.ExecuteStockroom(*Debug, places_ptr, works_ptr, workplaces_ptr)

		log.Printf("Extecuted works: " + strconv.Itoa(ExecutedWorks))
	}

}
