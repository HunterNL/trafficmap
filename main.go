package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	Id           string `json:"id"`
	image        []byte
	Lat          string   `json:"lat"`
	Lon          string   `json:"lon"`
	Name         string   `json:"name"`
	ImageWidth   int      `json:"imageWidth"`
	ImageHeight  int      `json:"imageHeight"`
	Working      bool     `json:"working"`
	RoadId       string   `json:"roadId"`
	RoadSide     string   `json:"roadSide"`
	RoadOffset   int      `json:"roadOffset"`
	Organization string   `json:"organization"`
	TextLines    []string `json:"text"`
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
	downloadOnly := flag.Bool("download", false, "Only download images and quit")
	outDir := flag.String("outdir", ".", "Output directory for files")
	host := flag.String("host", "0.0.0.0", "Network addres to use")
	port := flag.Int("port", 3000, "Port to serve http on")

	flag.Parse()

	if *downloadOnly {
		error := outputImages(*sourceUrl, *outDir)
		if error != nil {
			log.Fatalln(error)
		}
		return
	}

	serv := newServ()
	ticker := time.NewTicker(UpdateInterval)
	err := updateDrips(*sourceUrl, &serv)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Succesfully got data from %v\n", *sourceUrl)
	go update(*sourceUrl, ticker.C, &serv)

	// placeDripsFile()
	ServeData(*host, *port, &serv)
}

// // Ensures a given directory relative to the workig directory exists
// func ensureTargetDirectoryExists(dirName string) error {
// 	stat, err := os.Stat(dirName)

// 	if errors.Is(err, fs.ErrNotExist) {
// 		err = os.MkdirAll(dirName, os.ModeDir)
// 		if err != nil {
// 			return fmt.Errorf("error creating directory: %w", err)
// 		}
// 	}

// 	if err != nil {
// 		return fmt.Errorf("error getting target directory stats: %w", err)
// 	}

// 	if stat.IsDir() {
// 		return nil
// 	} else {
// 		return fmt.Errorf("target path is not a directory")
// 	}

// }

func outputImages(baseUrl, outDir string) error {
	dripsFile, err := getFile(baseUrl, dripStatusFile, true)
	if err != nil {
		return err
	}

	images, err := imagesFromFile(dripsFile)
	if err != nil {
		return err
	}

	targetDir := filepath.ToSlash(filepath.Clean(outDir))

	err = os.MkdirAll(targetDir, os.ModeDir)
	if err != nil {
		return fmt.Errorf("error while ensuring output directory exists: %w", err)
	}

	for id, img := range images {
		fileName := filepath.ToSlash(filepath.Join("./"+targetDir, id+".png"))

		err := os.WriteFile(fileName, img, os.ModeType)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error writing file", err)
		}
	}

	path, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}

	fmt.Printf("Written %v images to %v\n", len(images), path)

	return nil

}
