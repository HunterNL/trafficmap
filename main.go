package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const UpdateInterval = time.Minute * 5

type DripServ struct {
	sync.Mutex
	dripsMap   map[string]Drip
	DripsSlice []Drip `json:"drips"`
	LastUpdate time.Time
}

func newServ() DripServ {
	return DripServ{
		dripsMap:   make(map[string]Drip),
		DripsSlice: make([]Drip, 0),
	}
}

type Drip struct {
	Id          string `json:"id"`
	image       []byte
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Description string `json:"description"`
}

func (d *Drip) hasImage() bool {
	if d.image == nil {
		return false
	}

	if len(d.image) == 0 {
		return false
	}

	return true
}

func update(sourceUrl string, c <-chan time.Time, serv *DripServ) {
	for range c {
		start := time.Now()
		updateDrips(sourceUrl, serv)
		fmt.Printf("Updating from %v took %v\n", sourceUrl, time.Since(start))
	}
}

func main() {
	sourceUrl := flag.String("sourceURL", "http://opendata.ndw.nu/", "Full URL to retrieve the source data from")

	flag.Parse()

	serv := newServ()
	ticker := time.NewTicker(UpdateInterval)
	err := updateDrips(*sourceUrl, &serv)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Succesfully got data from %v\n", *sourceUrl)
	go update(*sourceUrl, ticker.C, &serv)

	// placeDripsFile()
	ServeData(&serv)
}

func placeDripsFile() {

	var err error
	dripsFile, err := os.ReadFile("./DRIPS.xml")
	if err != nil {
		fmt.Printf("Error reading drips file: %v\n", err)
		return
	}

	locFile, err := os.ReadFile("./LocatietabelDRIPS.xml")
	if err != nil {
		fmt.Printf("Error reading location file: %v\n", err)
		return
	}

	drips, err := parseDripsXML(dripsFile, locFile)

	if err != nil {
		fmt.Printf("Error unmarshalling xml: %v\n", err)
		return
	}

	os.Mkdir("./static", os.ModeType)
	os.Mkdir("./static/images", os.ModeType)

	for _, p := range drips {
		if p.image == nil { //Skip empty
			continue
		}
		// // imgData, err := base64.StdEncoding.DecodeString(p.image)
		// if err != nil {
		// 	fmt.Printf("Error decoding image %v\n", err)
		// }

		writeErr := os.WriteFile("./static/images/"+p.Id+".png", p.image, os.ModeType)
		if writeErr != nil {
			fmt.Println(writeErr)
		}
	}

	jsonFile, encodingErr := marshallJson(drips)
	if encodingErr != nil {
		fmt.Print(encodingErr)
		return
	}

	writeErr := os.WriteFile("./static/data.json", jsonFile, os.ModeType)
	if writeErr != nil {
		fmt.Println(writeErr)
	}
}
