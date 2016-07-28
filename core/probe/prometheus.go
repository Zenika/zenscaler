package probe

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Prometheus probe
//
// Metrics are retrived using HTTP/Text protocol
// (see https://prometheus.io/docs/instrumenting/exposition_formats/)
// TODO protobuf support ?
type Prometheus struct {
	URL string `json:"url"` // Sample: http://localhost:9100/metrics
	Key string `json:"key"` // Sample: node_cpu{cpu="cpu6",mode="idle"}
}

// Name of the probe
func (p Prometheus) Name() string {
	return "Prometheus probe for " + p.URL + " [" + p.Key + "]"
}

// Value make the request and parse content
func (p Prometheus) Value() (float64, error) {
	resp, err := http.Get(p.URL)
	if err != nil {
		return -1.0, err
	}
	if resp.StatusCode != http.StatusOK {
		return -1.0, err
	}
	fvalue, err := p.findValue(resp.Body)
	if err != nil {
		return -1.0, err
	}
	return fvalue, nil
}

// find the matching token and parse probe value
func (p Prometheus) findValue(body io.Reader) (float64, error) {
	scanLines := bufio.NewScanner(body)
	for scanLines.Scan() { // for each line
		// now we're splitting by spaces
		scanWords := bufio.NewScanner(strings.NewReader(scanLines.Text()))
		scanWords.Split(bufio.ScanWords)

		if scanWords.Scan() { // line is not empty
			switch scanWords.Text() {
			case p.Key:
				if scanWords.Scan() {
					fvalue, err := strconv.ParseFloat(scanWords.Text(), 64)
					if err != nil {
						fmt.Printf("%s", err)
						return -1.0, err
					}
					return fvalue, nil
				}
			default:
				// if it doesn't start with the key, discard the line
				break
			}

		}
	}
	return 0, fmt.Errorf("Token %s not found", p.Key)
}
