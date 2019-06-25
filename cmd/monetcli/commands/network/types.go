package network

type configurationRecord struct {
	//location of the network.toml file
	config     *configRecord    `mapstructure:"config"`
	validators *validatorRecord `mapstructure:"validators"`
	poa        *poaRecord       `mapstructure:"config"`
}

type configRecord struct {
	dataDir string `mapstructure:"datadir"`
}

type poaRecord struct {
	contractAddress string `mapstructure:"contractaddress"`
	contractName    string `mapstructure:"contractname"`
	compilerVersion string `mapstructure:"compilerversion"`
	byteCode        string `mapstructure:"bytecode"`
	abi             string `mapstructure:"abi"`
}

type validatorRecord struct {
	moniker            string `mapstructure:"monikers"`
	address            string `mapstructure:"addresses"`
	pubkeys            string `mapstructure:"pubkeys"`
	ip                 string `mapstructure:"ips"`
	isInitialValidator string `mapstructure:"isvalidator"`
}

type genesisAllocRecord struct {
	Balance string `json:"balance"`
	Moniker string `json:"moniker"`
}

type genesisAlloc map[string]*genesisAllocRecord

type genesisPOA struct {
	Address string `json:"address"`
	Abi     string `json:"abi"`
	Code    string `json:"code"`
}

type genesisFile struct {
	Alloc *genesisAlloc `json:"alloc"`
	Poa   *genesisPOA   `json:"poa"`
}

type peerRecord struct {
	NetAddr   string `json:"NetAddr"`
	PubKeyHex string `json:"PubKeyHex"`
	Moniker   string `json:"Moniker"`
}

type peerRecordList []*peerRecord

var (
	config       configurationRecord
	configConfig configRecord
	poa          poaRecord
)

const (
	defaultContractAddress = "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba"
	defaultContractName    = "genesis_array.sol"
)

func defaultConfig() {

	home, err := defaultHomeDir()
	if err == nil {
		networkViper.SetDefault("config.datadir", home)
	}
	networkViper.SetDefault("poa.contractaddress", defaultContractAddress)
	networkViper.SetDefault("poa.contractname", defaultContractName)
	networkViper.SetDefault("poa.compilerversion", "")
	networkViper.SetDefault("validators.monikers", "")
	networkViper.SetDefault("validators.addresses", "")
	networkViper.SetDefault("validators.pubkeys", "")
	networkViper.SetDefault("validators.ips", "")
	networkViper.SetDefault("validators.isvalidator", "")
}

func newConfigurationRecord() *configurationRecord {

	configConfig = configRecord{dataDir: ""}
	poa = poaRecord{contractAddress: defaultContractAddress, contractName: defaultContractName}

	config = configurationRecord{config: &configConfig, poa: &poa}

	return &config

}
