package main

import (
	"flag"
	"fmt"
	"log"
	congruit "main/congruit-go/libs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func HelloServer(w http.ResponseWriter, req *http.Request, t string, StockRoomDir string, Debug bool, GitRepo string) int {

	if len(GitRepo) > 0 {
		StockRoomDir = congruit.CloneRepo(GitRepo)
	}

	ExecutedWorks := 0

	w.Header().Set("Content-Type", "text/plain")

	if strings.EqualFold(t, (req.Header.Get("Token"))) {

		w.Write([]byte("Hi dude!\n"))

		wp := strings.Split(req.Header.Get("Workplaces"), ",")

		for w := range wp {
			log.Println("A remote agent asks to run " + wp[w])
			_, err := os.Stat(StockRoomDir + "/workplaces_enabled/" + wp[w])
			if err != nil {
				err := os.Link(StockRoomDir+"/workplaces/"+wp[w], StockRoomDir+"/workplaces_enabled/"+wp[w])
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if len(req.Header.Get("Workplace")) > 0 {

			log.Println("A remote agent asks to run " + req.Header.Get("Workplace"))
			_, err := os.Stat(StockRoomDir + "/workplaces_enabled/" + req.Header.Get("Workplace"))

			if err != nil {
				err := os.Link(StockRoomDir+"/workplaces/"+req.Header.Get("Workplace"), StockRoomDir+"/workplaces_enabled/"+req.Header.Get("Workplace"))
				if err != nil {
					log.Fatal("Error in enabling workplace!")
				}
			}

		}

		places, works, workplaces := congruit.LoadStockroom(StockRoomDir, Debug)

		ExecutedWorks = congruit.ExecuteStockroom(Debug, places, works, workplaces)
		w.Write([]byte("\n Remote executed works: " + strconv.Itoa(ExecutedWorks) + "\n"))
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
	ssl_cert := flag.String("ssl_cert", "./", "path of a ssl certificate")
	ssl_key := flag.String("ssl_key", "./", "path of a ssl key")

	ExecutedWorks := 0

	flag.Parse()

	fmt.Printf("                         _ _   \n")
	fmt.Printf(" ___ ___ ___ ___ ___ _ _|_| |_ \n")
	fmt.Printf("|  _| . |   | . |  _| | | |  _|\n")
	fmt.Printf("|___|___|_|_|_  |_| |___|_|_|  \n")
	fmt.Printf("            |___|              \n")
	fmt.Println("Version:", version)

	if *CongruitVersion {
		return
	}

	if len(*GitRepo) > 0 {
		*StockRoomDir = congruit.CloneRepo(*GitRepo)
	}

	wp := strings.Split(*WorkPlaces, ",")

	os.Mkdir(*StockRoomDir+"/workplaces_enabled/", 0755)

	for w := range wp {
		_, err := os.Stat(*StockRoomDir + "/workplaces_enabled/" + wp[w])
		if err != nil {
			err := os.Link(*StockRoomDir+"/workplaces/"+wp[w], *StockRoomDir+"/workplaces_enabled/"+wp[w])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if *Friend {

		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			HelloServer(w, r, *Token, *StockRoomDir, *Debug, *GitRepo)
		})

		err := http.ListenAndServeTLS(":8443", *ssl_cert, *ssl_key, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}

	places, works, workplaces := congruit.LoadStockroom(*StockRoomDir, *Debug)

	ExecutedWorks = congruit.ExecuteStockroom(*Debug, places, works, workplaces)

	log.Printf("Executed works: " + strconv.Itoa(ExecutedWorks))

	for *Supervisor {

		log.Printf("start congruit supervisor...")

		time.Sleep(5000 * time.Millisecond)

		ExecutedWorks = congruit.ExecuteStockroom(*Debug, places, works, workplaces)

		log.Printf("Executed works: " + strconv.Itoa(ExecutedWorks))
	}

}
