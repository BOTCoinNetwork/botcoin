package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/src/engine"
	"github.com/mosaicnetworks/monetd/src/babble"
	"github.com/mosaicnetworks/monetd/src/common"
	"github.com/mosaicnetworks/monetd/src/configuration"
	"github.com/spf13/cobra"
)

/*******************************************************************************
RunCmd
*******************************************************************************/

//newRunCmd returns the command that starts the daemon
func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a MONET node",
		Long: `
Run a MONET node.
	
Start a daemon which acts as a full node on a MONET network. All data and 
configuration are stored under a directory [datadir] controlled by the 
--datadir flag ($HOME/.monet by default on UNIX systems). 

[datadir] must contain a set of files defining the network that this node is 
attempting to join or create. Please refer to monetd config for tools to manage 
this configuration. 

Further options pertaining to the operation of the node are read from the 
[datadir]/monetd.toml file, or overwritten by the following flags. The following 
command displays the expected output:

monetd config location `,

		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			common.DebugMessage(fmt.Sprintf("Base Config: %+v", configuration.Global.BaseConfig))
			common.DebugMessage(fmt.Sprintf("Babble Config: %+v", configuration.Global.Babble))
			common.DebugMessage(fmt.Sprintf("Eth Config: %+v", configuration.Global.Eth))
			return nil
		},

		RunE: runMonet,
	}

	bindFlags(cmd)

	return cmd
}

func bindFlags(cmd *cobra.Command) {
	// EVM-Lite and Babble share the same API address
	cmd.Flags().String("api-listen", configuration.Global.APIAddr, "IP:PORT of Monet HTTP API service")

	// Babble config
	cmd.Flags().String("babble.listen", configuration.Global.Babble.BindAddr, "IP:PORT of Babble node")
	cmd.Flags().Duration("babble.heartbeat", configuration.Global.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Duration("babble.timeout", configuration.Global.Babble.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("babble.cache-size", configuration.Global.Babble.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("babble.sync-limit", configuration.Global.Babble.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("babble.max-pool", configuration.Global.Babble.MaxPool, "Max number of pool connections")
	cmd.Flags().Bool("babble.bootstrap", configuration.Global.Babble.Bootstrap, "Bootstrap Babble from database")

	// Eth config
	cmd.Flags().Int("eth.cache", configuration.Global.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")
}

/*******************************************************************************
READ CONFIG AND RUN
*******************************************************************************/

// Run the EVM-Lite / Babble engine
func runMonet(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(
		configuration.Global.ToBabbleConfig(),
		configuration.Global.Logger("babble-proxy"))

	engine, err := engine.NewEngine(*configuration.Global.ToEVMLConfig(), babble)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
