package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"time"

	"github.com/tarm/serial"
)

// Pre-compile regex and extract index
var frameSplitRegexp = regexp.MustCompile(`(?P<label>\S+)(?P<delimiter1>\s)(?P<data>\S+)(?P<delimiter2>\s)(?P<checksum>\S)`)
var labelIndex = frameSplitRegexp.SubexpIndex("label")
var delimiter1Index = frameSplitRegexp.SubexpIndex("delimiter1")
var dataIndex = frameSplitRegexp.SubexpIndex("data")
var checksumIndex = frameSplitRegexp.SubexpIndex("checksum")

var collectLabel = []string{
	"ADCO",
	"ISOUSC",
	"BASE",
	"BBRHCJB",
	"BBRHPJB",
	"BBRHCJW",
	"BBRHPJW",
	"BBRHCJR",
	"BBRHPJR",
	"PAPP",
	"IINST",
}

// Main function that loop on the serial port output
func getSerialTeleinfo(teleinfoMetric *TeleinfoMetrics) {
	go func() {
		config := &serial.Config{
			Name:        "/dev/ttyUSB0",
			Baud:        1200,
			ReadTimeout: time.Millisecond * 500,
			Size:        7,
			Parity:      serial.ParityEven,
			StopBits:    serial.Stop1,
		}

		stream, err := serial.OpenPort(config)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(stream)
		var label, data string
		for scanner.Scan() {
			//fmt.Println("frame dump:" + scanner.Text()) // Println will add back the final '\n'
			if label, data, err = frameExtractMode1(scanner.Text()); err != nil {
				fmt.Println(err.Error())
			} else {
				//fmt.Println(label + "\t" + data)
				if slices.Contains(collectLabel, label) {
					teleinfoMetric.Set(label, data)
				}
			}
			//fmt.Printf("%+v\n", teleinfoMetric)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
}

/*
	From: Enedis-NOI-CPT_02E.pdf
	LineFeed (removed by serial): 0x0A
	  Label
	  SP or HT: 0x20 or 0x09
	  Data
	  SP or HT: 0x20 or 0x09
	  Checksum
	CarryReturn (removed by serial): 0x0D
	Checksum in mode 1: Label + SP/HT + Data
	Checksum in mode 2: Label + SP/HT + Data + SP/HT
*/

// Implementation of Linky checksum, return true/false
func frameChecksum(payload string, checksum byte) bool {
	sum := 0
	for _, char := range payload {
		sum = sum + int(char)
	}
	sum = (sum & 0x3F) + 0x20
	return (sum == int(checksum))
}

// Return valid label ("étiquette") and data ("donnée") from the raw frame in mode 1 ("historique")
func frameExtractMode1(frame string) (string, string, error) {
	// Split attempt
	if frameSplit := frameSplitRegexp.FindStringSubmatch(frame); frameSplit == nil || len(frameSplit) != 6 {
		// Split failed
		return "", "", errors.New("teleinfo frame does not match the regexp (" + frame + "), ignoring")
	} else {
		// Split success
		if frameChecksum(frameSplit[labelIndex]+frameSplit[delimiter1Index]+frameSplit[dataIndex], frameSplit[checksumIndex][0]) {
			return frameSplit[labelIndex], frameSplit[dataIndex], nil
		} else {
			return "", "", errors.New("teleinfo frame checksum failed (" + frame + "), ignoring")
		}
	}
}

// Return valid label ("étiquette") and data ("donnée") from the raw frame in mode 2 ("standard")
/*func frameExtractMode2(frame string) (string, string, error) {
	//ToDo
}*/
