package network

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mosaicnetworks/monetd/src/config"
	"github.com/mosaicnetworks/monetd/src/contract"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/mosaicnetworks/monetd/src/types"
	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
	"github.com/spf13/cobra"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [network_name]",
		Short: "create the configuration for a multi-node network",
		Long: `
giverny network build
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkBuild,
	}

	return cmd
}

func networkBuild(cmd *cobra.Command, args []string) error {
	return buildNetwork(strings.TrimSpace(args[0]))
}

//buildNetwork builds the network. It is called directly from new command as well.
func buildNetwork(networkName string) error {
	if !common.CheckMoniker(networkName) {
		return errors.New("network name, " + networkName + ", is invalid")
	}

	// Check all the files and directories we expect actually exist
	thisNetworkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)

	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}

	tree, err := files.LoadToml(networkTomlFile)
	if err != nil {
		common.ErrorMessage("Cannot load network.toml file: ", networkTomlFile)
		return err
	}

	common.DebugMessage("Building network " + networkName)

	err = dumpPeersJSON(tree, thisNetworkDir)
	if err != nil {
		common.ErrorMessage("Error writing peers json file")
		return err
	}

	return nil
}

func dumpPeersJSON(tree *toml.Tree, thisNetworkDir string) error {

	var peers types.PeerRecordList

	if tree.HasPath([]string{"Name"}) {
		netName := tree.GetPath([]string{"Name"}).(string)
		common.DebugMessage("Network Name ", netName)
	}

	common.DebugMessage("dumpPeersJSON")

	nodesquery, err := query.CompileAndExecute("$.nodes", tree)
	if err != nil {
		common.ErrorMessage("Error loading nodes")
		return err
	}

	var alloc = make(config.GenesisAlloc)

	for _, value := range nodesquery.Values() {

		//		common.DebugMessage(reflect.TypeOf(value).String())
		//		common.DebugMessage("Found a value: "+strconv.Itoa(i), value)

		if reflect.TypeOf(value).String() == "[]*toml.Tree" {
			nodes := value.([]*toml.Tree)

			for _, tr := range nodes {
				var addr, moniker, netaddr, pubkey, tokens string

				if tr.HasPath([]string{"moniker"}) {
					moniker = tr.GetPath([]string{"moniker"}).(string)
				}
				if tr.HasPath([]string{"netaddr"}) {
					netaddr = tr.GetPath([]string{"netaddr"}).(string)
					if !strings.Contains(netaddr, ":") {
						netaddr += ":" + monetconfig.DefaultGossipPort
					}
				}
				if tr.HasPath([]string{"pubkey"}) {
					pubkey = tr.GetPath([]string{"pubkey"}).(string)
				}
				if tr.HasPath([]string{"tokens"}) {
					tokens = tr.GetPath([]string{"tokens"}).(string)
				}
				if tr.HasPath([]string{"address"}) {
					addr = tr.GetPath([]string{"address"}).(string)
				}

				rec := config.GenesisAllocRecord{Moniker: moniker, Balance: tokens}
				alloc[addr] = &rec

				if tr.HasPath([]string{"validator"}) && (!tr.GetPath([]string{"validator"}).(bool)) {
					continue
				}

				peers = append(peers, &types.PeerRecord{Moniker: moniker,
					NetAddr:   netaddr,
					PubKeyHex: pubkey})

				// If we reach this point this node is a validator
				if err := createKeyFileIfNotExists(thisNetworkDir, moniker, addr, pubkey); err != nil {
					return err
				}

			}

		}
	}

	peersJSONOut, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}

	jsonFileName := filepath.Join(thisNetworkDir, monetconfig.PeersJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut))
	if err != nil {
		return err
	}

	// Write copy of peers.json to peers.genesis.json
	jsonFileName = filepath.Join(thisNetworkDir, monetconfig.PeersGenesisJSON)
	err = files.WriteToFile(jsonFileName, string(peersJSONOut))
	if err != nil {
		return err
	}

	err = BuildGenesisJSON(thisNetworkDir, peers, monetconfig.DefaultContractAddress, alloc)
	if err != nil {
		return err
	}

	return err
}

func createKeyFileIfNotExists(configDir string, moniker string, addr string, pubkey string) error {
	keyfile := filepath.Join(configDir, monetconfig.KeyStoreDir, moniker+".json")
	if files.CheckIfExists(keyfile) {
		return nil
	} // If exists, nothing to do

	type minjson struct {
		Address string `json:"address"`
		Pub     string `json:"pub"`
	}

	j := minjson{Address: addr, Pub: pubkey}
	out, err := json.Marshal(j)
	if err != nil {
		return err
	}

	err = files.WriteToFile(keyfile, string(out))
	if err != nil {
		return err
	}

	return nil
}

//BuildGenesisJSON compiles and build a genesis.json file
func BuildGenesisJSON(configDir string, peers types.PeerRecordList, contractAddress string, alloc config.GenesisAlloc) error {
	var genesis config.GenesisFile

	common.DebugMessage("buildGenesisJSON")

	finalSource, err := contract.GetFinalSoliditySource(peers)
	if err != nil {
		return err
	}

	common.DebugMessage("Source Loaded")

	genesispoa, err := config.BuildGenesisPOAJSON(finalSource, configDir, contractAddress, false)
	if err != nil {
		return err
	}
	genesis.Poa = &genesispoa

	common.DebugMessage("POA Section Build")

	/* alloc, err := buildGenesisAlloc(filepath.Join(configDir, monetconfig.KeyStoreDir))
	if err != nil {
		return err
	} */
	genesis.Alloc = &alloc

	common.DebugMessage("Alloc Built")

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	common.DebugMessage("Write Genesis.json")
	jsonFileName := filepath.Join(configDir, monetconfig.GenesisJSON)
	files.WriteToFile(jsonFileName, string(genesisjson))

	return nil
}
