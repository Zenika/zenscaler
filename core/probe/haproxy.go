package probe

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// HAproxy configurable probe.
//
// Metrics are retrived by accessing HAProxy command socket
type HAproxy struct {
	Socket string `json:"socket"`
	Type   string `json:"type"`
	Item   string `json:"item"`
}

// Name of the probe
func (ha HAproxy) Name() string {
	return "HAproxy probe for " + ha.Type + "." + ha.Item
}

// Value probe the target and report back values
func (ha HAproxy) Value() (float64, error) {
	statsMap, err := ha.getStats(ha.Type)
	if err != nil {
		return 0, fmt.Errorf("cannot probe hap: %s, check access rights", err)
	}
	value, err := strconv.ParseFloat(statsMap[ha.Item][1], 64)
	if err != nil {
		return 0, fmt.Errorf("Cannot parse float: %s", err)
	}
	return value, nil
}

// Some code from github.com/tnolet/haproxy-rest
// See https://cbonte.github.io/haproxy-dconv/configuration-1.5.html#show%20stat
// See https://www.datadoghq.com/blog/monitoring-haproxy-performance-metrics/
// See https://cbonte.github.io/haproxy-dconv/configuration-1.5.html#9.1
func (ha HAproxy) getStats(statsType string) (map[string][]string, error) {
	var cmdString string

	switch statsType {
	case "all":
		cmdString = "show stat -1\n"
	case "backend":
		cmdString = "show stat -1 2 -1\n"
	case "frontend":
		cmdString = "show stat -1 1 -1\n"
	case "server":
		cmdString = "show stat -1 4 -1\n"
	default:
		return nil, errors.New("Unknown stat type")
	}

	result, err := ha.HaproxyCmd(cmdString)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(strings.NewReader(strings.Trim(result, "# ")))
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.New("Failed to decode CSV")
	}

	// Turn records into a map
	mappedRecords := make(map[string][]string)
	lines := len(records)
	for i, headline := range records[0] {
		mappedRecords[headline] = make([]string, lines-1)
		for j := 1; j < lines; j++ {
			mappedRecords[headline][j-1] = records[j][i]
		}
	}
	return mappedRecords, nil
}

// HaproxyCmd execution on the unix socket
func (ha *HAproxy) HaproxyCmd(cmd string) (string, error) {
	// connect to haproxy
	conn, errConn := net.Dial("unix", ha.Socket)
	if errConn != nil {
		return "", fmt.Errorf("unable to connect to HAproxy socket: %s", errConn)
	}
	// #nosec close unix socket, no information written
	defer func() { _ = conn.Close() }()

	fmt.Fprint(conn, cmd)
	response := ""
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response += (scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return response, nil
}
