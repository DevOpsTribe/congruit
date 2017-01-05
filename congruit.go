package main

import (
	"./congruit-go/libs"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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

		log.Println("A remote agent asks to run " + req.Header.Get("Workplace") )

		for i := range workplaces_ptr_temp {

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

	const version string = "1.1.0"

	StockRoomDir := flag.String("stockroom-dir", "./", "stockroom directory")
	GitRepo := flag.String("gitrepo", "", "Git repo of stockroom")
	WorkPlaces := flag.String("workplaces", "", "Workplace to apply")
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

	if len(*GitRepo) > 0 {

		if _, err := os.Stat("/tmp/stockroom"); err == nil {
			os.RemoveAll("/tmp/stockroom")
		}

		cmd := exec.Command("git", "clone", *GitRepo, "/tmp/stockroom")

		*StockRoomDir = "/tmp/stockroom"

		err := cmd.Run()

		if err != nil {

			log.Fatal("Error when pull stockroom repo")

		}
	}

	wp := strings.Split(*WorkPlaces, ",")

	for w := range wp {
		_, err := os.Stat(*StockRoomDir + "/workplaces_enabled/" + wp[w])
		if err != nil {
			err := os.Link(*StockRoomDir+"/workPlaces/"+wp[w], *StockRoomDir+"/workplaces_enabled/"+wp[w])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	works_ptr := []*congruit.Work{}
	workplaces_ptr := []*congruit.WorkPlace{}
	places_ptr := []*congruit.Place{}

	if *Friend {

		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			HelloServer(w, r, *Token, *StockRoomDir, *Debug)
		})

		err := http.ListenAndServeTLS(":8443", "domain.crt", "domain.key", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}

	places_ptr, works_ptr, workplaces_ptr = congruit.LoadStockroom(*StockRoomDir, *Debug)

	ExecutedWorks = congruit.ExecuteStockroom(*Debug, places_ptr, works_ptr, workplaces_ptr)

	log.Printf("Extecuted works: " + strconv.Itoa(ExecutedWorks))

	for *Supervisor {

		log.Printf("start congruit supervisor...")

		time.Sleep(5000 * time.Millisecond)

		ExecutedWorks = congruit.ExecuteStockroom(*Debug, places_ptr, works_ptr, workplaces_ptr)

		log.Printf("Extecuted works: " + strconv.Itoa(ExecutedWorks))
	}

}
