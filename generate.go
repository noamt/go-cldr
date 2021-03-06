// +build ignore

// Generate generate.go. Can by running
// go generate
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"golang.org/x/text/unicode/cldr"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const coreURL = "https://unicode.org/Public/cldr/37/core.zip"

func main() {
	os.Exit(generateSupplemental())
}

func generateSupplemental() int {
	workingDirectory, wdError := os.Getwd()
	if wdError != nil {
		log.Println("Failed to get working directory", wdError)
		return 1
	}
	decoded, cldrError := getCLDR()
	if cldrError != nil {
		log.Println("Failed to get CLDR", cldrError)
		return 1
	}
	if generateError := generateSupplementalFromTemplate(workingDirectory, decoded); generateError != nil {
		log.Println("Failed to generateSupplemental from template", generateError)
		return 1
	}
	return 0
}

func generateSupplementalFromTemplate(workingDirectory string, decoded *cldr.CLDR) error {
	t := template.Must(template.New("supplemental").Parse(supplementalTemplate))

	var buf bytes.Buffer
	executeError := t.Execute(&buf, map[string]interface{}{
		"FirstDays": getFirstDays(decoded),
	})
	if executeError != nil {
		return fmt.Errorf("failed to execute supplemental template: %w", executeError)
	}
	formattedTemplate, formatError := format.Source(buf.Bytes())
	if formatError != nil {
		return fmt.Errorf("failed to format supplemental output: %w", formatError)
	}

	supplementalFilePath := filepath.Join(workingDirectory, "supplemental", "supplemental.go")
	os.Remove(supplementalFilePath)
	supplementalFile, supplementalFileError := os.OpenFile(supplementalFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if supplementalFileError != nil {
		return fmt.Errorf("failed to open file %s: %w", supplementalFilePath, supplementalFileError)
	}
	defer supplementalFile.Close()
	_, copyError := io.Copy(supplementalFile, bytes.NewReader(formattedTemplate))
	if copyError != nil {
		return fmt.Errorf("failed to write to file %s: %w", supplementalFilePath, supplementalFileError)
	}

	return nil
}

func getFirstDays(decoded *cldr.CLDR) map[string]string {
	firstDays := map[string]string{}
	for _, firstDay := range decoded.Supplemental().WeekData.FirstDay {
		if firstDay.Alt != "" {
			continue
		}
		day := firstDay.Day
		territories := strings.Fields(firstDay.Territories)
		for _, territory := range territories {
			firstDays[territory] = day
		}
	}
	return firstDays
}

func getCLDR() (*cldr.CLDR, error) {
	println("Getting core file from", coreURL)
	getCLDRCore, getCLDRCoreError := http.Get(coreURL)
	if getCLDRCoreError != nil {
		return nil, fmt.Errorf("failed to get core zip file: %w", getCLDRCoreError)
	}
	if getCLDRCore.StatusCode >= http.StatusBadRequest {
		body, _ := ioutil.ReadAll(getCLDRCore.Body)
		return nil, fmt.Errorf("failed to get core zip file. Status: %d. Error: %s", getCLDRCore.StatusCode, string(body))
	}
	defer getCLDRCore.Body.Close()

	println("Decoding...")
	var d cldr.Decoder
	d.SetDirFilter("supplemental")
	decoded, decodeError := d.DecodeZip(getCLDRCore.Body)
	if decodeError != nil {
		return nil, fmt.Errorf("failed to get core zip file: %w", getCLDRCoreError)
	}
	return decoded, nil
}

var supplementalTemplate = `
// Code generated by go generate; DO NOT EDIT.
package supplemental

import (
	"golang.org/x/text/language"
	"time"
)

const defaultTerritory = "001"

var territoryFirstDays = map[string]string{
{{ range $territory, $firstDay := .FirstDays }}"{{ $territory }}": "{{ $firstDay }}",
{{ end }}
}

type firstDays struct{}
var FirstDay = firstDays{}

func (f *firstDays) ByRegion(region language.Region) time.Weekday {
	firstDay, ok := territoryFirstDays[region.String()]
	if !ok {
		firstDay = territoryFirstDays[defaultTerritory]
	}

	switch firstDay {
	case "fri":
		return time.Friday
	case "sat":
		return time.Saturday
	case "mon":
		return time.Monday
	default:
		return time.Sunday
	}
}
`
