package network

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/evm-lite/src/currency"

	"github.com/mosaicnetworks/monetd/src/common"
)

var nodeNamePrefix = "node"

func getNodesWithNames(srcFile string, numNodes int, numValidators int, initialIP string) ([]node, string, error) {
	var rtn []node
	lastDigit := 0
	IPStem := ""
	var err error

	if initialIP != "" {
		splitIP := strings.Split(initialIP, ".")

		if len(splitIP) != 4 {
			return rtn, "", errors.New("malformed initial IP: " + initialIP)
		}

		lastDigit, err = strconv.Atoi(splitIP[3])
		if err != nil {
			fmt.Println("lastDigit Set to Zero")
			lastDigit = 0
		} else {
			IPStem = strings.Join(splitIP[:3], ".") + "."
		}
	}

	netaddr := ""
	// if no input file specified, then we generate node0, node1 etc
	if srcFile == "" {
		for i := 0; i < numNodes; i++ {
			if IPStem != "" {
				netaddr = IPStem + strconv.Itoa(lastDigit+i)
			}

			rtn = append(rtn, node{Moniker: nodeNamePrefix + strconv.Itoa(i),
				NetAddr: netaddr, Validator: (numValidators < 1 || i < numValidators),
				Tokens: defaultTokens, Address: "", PubKeyHex: "", NonNode: false})
		}
		return rtn, netaddr, nil
	}

	// Read file line by line.

	file, err := os.Open(srcFile)

	if err != nil {
		common.ErrorMessage("failed opening file: ", err)
		return rtn, netaddr, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	i := 1
	for scanner.Scan() {

		var moniker string

		if IPStem != "" {
			netaddr = IPStem + strconv.Itoa(lastDigit+i-1)
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		} // Ignore blank lines

		if line[:1] == "#" {
			continue
		}

		validator := (numValidators < 1 || i <= numValidators)
		tokens := defaultTokens
		nonnode := false

		if strings.Contains(line, ",") {
			arrLine := strings.Split(line, ",")
			moniker = arrLine[0]
			common.DebugMessage("Setting moniker to " + moniker)
			if len(arrLine) > 1 {
				if strings.TrimSpace(arrLine[1]) != "" {
					netaddr = arrLine[1]
				}
			}

			if len(arrLine) > 2 {
				if strings.TrimSpace(arrLine[2]) != "" {
					tokens = currency.ExpandCurrencyString(arrLine[2])
				}
			}

			if len(arrLine) > 3 && len(strings.TrimSpace(arrLine[3])) > 0 {
				validator, _ = strconv.ParseBool(arrLine[3])
			}

			if len(arrLine) > 4 && len(strings.TrimSpace(arrLine[4])) > 0 {
				nonnode, _ = strconv.ParseBool(arrLine[4])
			}

		} else {
			moniker = line
		}

		if !common.CheckMoniker(moniker) {
			return rtn, netaddr, errors.New("node name " + moniker + " contains invalid characters")
		}

		rtn = append(rtn, node{Moniker: moniker,
			NetAddr: netaddr, Validator: validator,
			Tokens: tokens, Address: netaddr, PubKeyHex: "", NonNode: nonnode})

		if i >= numNodes && numNodes > 0 {
			break
		}
		i++
	}

	file.Close()
	return rtn, netaddr, nil
}
