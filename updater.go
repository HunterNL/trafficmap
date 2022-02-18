package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"
)

func updateDrips(baseUrl string, serv *DripServ) error {

	//Fetching
	dripResp, err := http.Get(baseUrl + "DRIPS.xml.gz")
	if err != nil {
		return err
	}
	defer dripResp.Body.Close()

	dripGz, err := io.ReadAll(dripResp.Body)
	if err != nil {
		return err
	}

	// os.WriteFile("./test.gz", dripGz, os.ModeAppend)

	locResp, err := http.Get(baseUrl + "LocatietabelDRIPS.xml.gz")
	if err != nil {
		return err
	}
	defer locResp.Body.Close()
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

	drips, err := parseDripsXML(dripsFile, locFile)
	if err != nil {
		return err
	}

	serv.Lock()
	defer serv.Unlock()

	serv.LastUpdate = time.Now()

	serv.DripsSlice = drips
	for _, drip := range drips {
		serv.dripsMap[drip.Id] = drip
	}

	fmt.Println("Drips updated")

	return nil
}
