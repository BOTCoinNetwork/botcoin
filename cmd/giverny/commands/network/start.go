package network

import (
	"errors"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/mosaicnetworks/monetd/src/config"

	"github.com/mosaicnetworks/monetd/cmd/giverny/configuration"
	"github.com/mosaicnetworks/monetd/src/common"
	monetconfig "github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/mosaicnetworks/monetd/src/docker"
	"github.com/mosaicnetworks/monetd/src/files"
	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//CLI flags
var forceNetwork = false
var useExisting = false
var startNodes = false

type copyRecord struct {
	from string
	to   string
}

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [network]",
		Short: "start a docker network",
		Long: `
giverny network start

Starts a network. Does not start individual nodes. The --force-network parameter
stops and restarts the network. 
		`,
		Args: cobra.ExactArgs(1),
		RunE: networkStart,
	}

	addStartFlags(cmd)

	return cmd
}

func addStartFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&forceNetwork, "force-network", forceNetwork, "force network down if already exists")
	cmd.Flags().BoolVar(&useExisting, "use-existing", useExisting, "use existing network if already exists")
	cmd.Flags().BoolVar(&startNodes, "start-nodes", startNodes, "start nodes")

	//	cmd.Flags().StringVar(&passwordFile, "passfile", "", "file containing the passphrase")
	viper.BindPFlags(cmd.Flags())
}

func networkStart(cmd *cobra.Command, args []string) error {
	network := args[0]

	if err := startDockerNetwork(network); err != nil {
		return err
	}

	return nil
}

func startDockerNetwork(networkName string) error {

	// Set some paths
	thisNetworkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	networkTomlFile := filepath.Join(thisNetworkDir, networkTomlFileName)

	// Check expect config exists
	if !files.CheckIfExists(thisNetworkDir) {
		return errors.New("cannot find the configuration folder, " + thisNetworkDir + " for " + networkName)
	}

	if !files.CheckIfExists(networkTomlFile) {
		return errors.New("cannot find the configuration file: " + networkTomlFile)
	}

	// Load Toml file to a tree
	tree, err := files.LoadToml(networkTomlFile)
	if err != nil {
		common.ErrorMessage("Cannot load network.toml file: ", networkTomlFile)
		return err
	}

	var dockerNetworkName, dockerbaseip, dockersubnet, dockeriprange, dockergateway string

	if tree.HasPath([]string{"docker", "name"}) {
		dockerNetworkName = tree.GetPath([]string{"docker", "name"}).(string)
	}
	if tree.HasPath([]string{"docker", "baseip"}) {
		dockerbaseip = tree.GetPath([]string{"docker", "baseip"}).(string)
	}
	if tree.HasPath([]string{"docker", "subnet"}) {
		dockersubnet = tree.GetPath([]string{"docker", "subnet"}).(string)
	}
	if tree.HasPath([]string{"docker", "iprange"}) {
		dockeriprange = tree.GetPath([]string{"docker", "iprange"}).(string)
	}
	if tree.HasPath([]string{"docker", "gateway"}) {
		dockergateway = tree.GetPath([]string{"docker", "gateway"}).(string)
	}

	common.DebugMessage("Configuring Network ", dockerNetworkName)
	common.DebugMessage("Base IP:    ", dockerbaseip)
	common.DebugMessage("IP Range:   ", dockeriprange)
	common.DebugMessage("Subnet:     ", dockersubnet)
	common.DebugMessage("Gateway:    ", dockergateway)

	if dockerNetworkName == "" {
		return errors.New("network " + networkName + " is not configured as a docker network")
	}

	// Create a Docker Client

	common.DebugMessage("Connecting to Docker Client")

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	// Create a Docker Network

	networkID, err := docker.SafeCreateNetwork(cli, dockerNetworkName,
		dockersubnet, dockeriprange, dockergateway, forceNetwork, useExisting)

	if err != nil {
		return err
	}

	common.DebugMessage("Created Network " + networkID)

	// Next we build the docker configurations to get all of the configs ready to
	// push

	nodes, err := exportDockerConfigs(tree, networkName)
	if err != nil {
		return err
	}

	if startNodes {
		for _, n := range nodes {
			common.DebugMessage("Starting node " + n)
			if err := pushDockerNode(networkName, n, networkID, imgName, imgIsRemote); err != nil {
				return err
			}
		}

	}

	return nil
}

//TODO return a []string of node names. Those names are passed back and used to loop through and
//push nodes.
func exportDockerConfigs(tree *toml.Tree, networkName string) ([]string, error) {
	var rtn []string

	// Configure some paths
	networkDir := filepath.Join(configuration.GivernyConfigDir, givernyNetworksDir, networkName)
	dockerDir := filepath.Join(networkDir, givernyDockerDir)
	err := files.CreateDirsIfNotExists([]string{dockerDir})
	if err != nil {
		return nil, err
	}

	//TODO rewrite using the unmarshalling code in the add command.
	//Need to abstract the node config creation into a single function
	//So we can call it from the add command.

	// Process the toml tree
	nodesquery, err := query.CompileAndExecute("$.nodes", tree)
	if err != nil {
		common.ErrorMessage("Error loading nodes")
		return nil, err
	}

	for _, value := range nodesquery.Values() {
		if reflect.TypeOf(value).String() == "[]*toml.Tree" {
			nodes := value.([]*toml.Tree)

			for _, tr := range nodes { // loop around nodes
				// Data wrangling
				var addr, moniker, netaddr, pubkey, tokens string
				var validator bool

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

				if tr.HasPath([]string{"validator"}) {
					validator = tr.GetPath([]string{"validator"}).(bool)
				}

				// Build output files

				if moniker != "" { // Should not be blank here, but safety first
					nodeDir := filepath.Join(dockerDir, moniker)
					// Docker container will always use .monet
					monetDir := filepath.Join(nodeDir, monetconfig.MonetdTomlDirDot)

					common.DebugMessage("Creating config in " + nodeDir)
					err := files.CreateDirsIfNotExists([]string{
						nodeDir,
						monetDir,
						filepath.Join(monetDir, monetconfig.BabbleDir),
						filepath.Join(monetDir, monetconfig.EthDir),
						filepath.Join(monetDir, monetconfig.KeyStoreDir),
					})
					if err != nil {
						return nil, err
					}

					copying := []copyRecord{
						{from: filepath.Join(networkDir, monetconfig.GenesisJSON),
							to: filepath.Join(monetDir, monetconfig.EthDir, monetconfig.GenesisJSON)},
						{from: filepath.Join(networkDir, monetconfig.PeersJSON),
							to: filepath.Join(monetDir, monetconfig.BabbleDir, monetconfig.PeersJSON)},
						{from: filepath.Join(networkDir, monetconfig.PeersGenesisJSON),
							to: filepath.Join(monetDir, monetconfig.BabbleDir, monetconfig.PeersGenesisJSON)},
						{from: filepath.Join(networkDir, monetconfig.MonetTomlFile),
							to: filepath.Join(monetDir, monetconfig.MonetTomlFile)},
						{from: filepath.Join(networkDir, monetconfig.KeyStoreDir, moniker+".json"),
							to: filepath.Join(monetDir, monetconfig.KeyStoreDir, moniker+".json")},
						{from: filepath.Join(networkDir, monetconfig.KeyStoreDir, moniker+".txt"),
							to: filepath.Join(monetDir, monetconfig.KeyStoreDir, moniker+".txt")},
					}

					for _, f := range copying {
						files.CopyFileContents(f.from, f.to)
					}

					// Write a node description file containing all of the parameters needed to start a container
					// Saves having to load and parse network.toml for every node
					nodeConfigFile := filepath.Join(dockerDir, moniker+".toml")
					nodeConfig := dockerNodeConfig{Moniker: moniker, NetAddr: strings.Split(netaddr, ":")[0]}

					tomlBytes, err := toml.Marshal(nodeConfig)
					if err != nil {
						return nil, err
					}

					err = files.WriteToFile(nodeConfigFile, string(tomlBytes))
					if err != nil {
						return nil, err
					}

					// Debug messages to kill the not used warning
					//TODO review and remove the items not used.
					common.DebugMessage("Address   : ", addr)
					common.DebugMessage("PubKey    : ", pubkey)
					common.DebugMessage("tokens    : ", tokens)
					common.DebugMessage("Validator : ", strconv.FormatBool(validator))

					// Need to edit monetd.toml and set datadir and listen appropriately

					err = config.SetLocalParamsInToml("/.monet", filepath.Join(monetDir, monetconfig.MonetTomlFile), netaddr)
					if err != nil {
						return nil, err
					}

					// Need to generate private key
					err = config.GenerateBabblePrivateKey(monetDir, moniker)
					if err != nil {
						return nil, err
					}

				}

				rtn = append(rtn, moniker)

			}
		}
	}

	return rtn, nil
}
