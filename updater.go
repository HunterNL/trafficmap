package main

import (
	"bytes"
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

func updateDrips(baseUrl string, serv *DripServ) error {

	sourceURL, err := url.Parse(baseUrl)

	if err != nil {
		return err
	}

	//Fetching
	dripResp, err := http.Get("http://" + sourceURL.Host + "/" + path.Join(sourceURL.Path, dripStatusFile))
	if err != nil {
		return err
	}
	defer dripResp.Body.Close()

	if dripResp.StatusCode != 200 {
		return fmt.Errorf("server responded with %v", dripResp.Status)
	}

	dripGz, err := io.ReadAll(dripResp.Body)
	if err != nil {
		return err
	}

	// os.WriteFile("./test.gz", dripGz, os.ModeAppend)

	locResp, err := http.Get("http://" + sourceURL.Host + "/" + path.Join(sourceURL.Path, dripLocationFile))
	if err != nil {
		return err
	}
	defer locResp.Body.Close()

	if locResp.StatusCode != 200 {
		return fmt.Errorf("server responded with %v", locResp.Status)
	}

	locGz, err := io.ReadAll(locResp.Body)
	if err != nil {
		return err
	}

	r, err := gzip.NewReader(bytes.NewReader(dripGz))

	if err != nil {
		return err
	}

	dripsFile, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	locR, err := gzip.NewReader(bytes.NewReader(locGz))
	if err != nil {
		return err
	}
	locFile, err := io.ReadAll(locR)
	if err != nil {
		return err
	}

	allDrips, err := parseDripsXML(dripsFile, locFile)
	if err != nil {
		return err
	}

	drips := make([]Drip, 0, len(allDrips))
	for _, d := range allDrips {
		if d.hasImage() {
			drips = append(drips, d)
		}
	}

	serv.Lock()
	defer serv.Unlock()

	serv.LastUpdate = time.Now()

	serv.DripsSlice = drips
	for _, drip := range drips {
		serv.dripsMap[drip.Id] = drip
	}

	return nil
}
