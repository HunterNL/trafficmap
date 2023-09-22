package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const dripStatusFile = "DRIPS.xml.gz"
const dripLocationFile = "LocatietabelDRIPS.xml.gz"

// Get a file from the given base url, optionally decompressing it
// BaseURL is parsed and only the host(+port) and path is used
// protocol is always set to http and the given filename is appended
func getFile(baseUrl, filePath string, gZip bool) ([]byte, error) {
	sourceURL, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	response, err := http.Get("http://" + sourceURL.Host + "/" + path.Join(sourceURL.Path, filePath))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("server responded with %v", response.Status)
	}

	var reader io.ReadCloser

	if gZip {
		reader, err = gzip.NewReader(response.Body)
	} else {
		reader = response.Body
	}

	if err != nil {
		return nil, err
	}

	return io.ReadAll(reader)
}

func updateDrips(baseUrl string, serv *DripServ) error {
	dripsFile, dripsErr := getFile(baseUrl, dripStatusFile, true)
	if dripsErr != nil {
		return fmt.Errorf("could not retrieve drip status file: %w", dripsErr)
	}

	locationsFile, locErr := getFile(baseUrl, dripLocationFile, true)
	if locErr != nil {
		return fmt.Errorf("could not retrieve drip location file: %w", locErr)
	}

	allDrips, err := ParseDripsXML(dripsFile, locationsFile)

	if err != nil {
		return err
	}

	// We only care about drips with an image or text, filter out the rest
	drips := make([]Drip, 0, len(allDrips))
	for _, d := range allDrips {
		if d.hasImage() || d.hasText() {
			drips = append(drips, d)
		}
	}

	// For description parsing development
	// #TODO Seperate executable/submode?
	// sb := bytes.Buffer{}

	// for _, d := range allDrips {
	// 	desc := description.Parse(d.Description)

	// 	sb.WriteString(desc.Name)
	// 	sb.WriteRune('\n')
	// }

	// os.WriteFile("names.txt", sb.Bytes(), os.ModeAppend)

	serv.Lock()
	defer serv.Unlock()

	serv.LastUpdate = time.Now()
	serv.DripsSlice = drips

	for k := range serv.dripsMap {
		delete(serv.dripsMap, k)
	}

	for _, drip := range drips {
		serv.dripsMap[drip.Id] = drip
	}

	return nil
}
