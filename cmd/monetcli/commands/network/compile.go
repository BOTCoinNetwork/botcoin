package network

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"github.com/mosaicnetworks/monetd/src/common"

	types "github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"
)

func compileConfig(cmd *cobra.Command, args []string) error {
	return CompileConfigWithParam(configDir)
}

//CompileConfigWithParam "finishes" the monetcli configuration, compiling the POA smart contract
//in preparation for a call to monetcli config publish
func CompileConfigWithParam(configDir string) error {
	// Load the Current Config

	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	// Retrieve and set the version number
	version, err := common.GetSolidityCompilerVersion()
	if err != nil {
		return err
	}

	tree.SetPath([]string{"poa", "compilerverison"}, version)

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

	if len(currentNodes) < 1 {
		return errors.New("Peerset is empty")
	}

	var alloc = make(common.GenesisAlloc)
	var peers common.PeerRecordList

	for _, value := range currentNodes {

		rawaddr := tree.GetPath([]string{"validators", value, "address"}).(string)
		rawmoniker := tree.GetPath([]string{"validators", value, "moniker"}).(string)
		rawpubkey := tree.GetPath([]string{"validators", value, "pubkey"}).(string)
		//	rawisvalidator := tree.GetPath([]string{"validators", value, "validator"}).(bool)
		rawip := tree.GetPath([]string{"validators", value, "ip"}).(string)

		// Convert Hex to Address and back out to get a EIP55 compliant address
		addr := types.HexToAddress(rawaddr).Hex()
		// Non-validators are added to the peer set, but not to the genesis peer set.
		peer := common.PeerRecord{NetAddr: rawip, PubKeyHex: rawpubkey, Moniker: rawmoniker}
		peers = append(peers, &peer)

		rec := common.GenesisAllocRecord{Moniker: rawmoniker, Balance: common.DefaultAccountBalance}
		alloc[addr] = &rec

	}

	/*
		//When contracts are "set" for a network, the solidity source is copied into the monetcli config directory
		//with a name of template.sol (defined by constant common.TemplateContract). Thus we can check just for that file.
		//If not found, then we download a fresh contract.
		filename := filepath.Join(configDir, common.TemplateContract)
		message("Checking for file: ", filename)

		soliditySource, err := common.GetSoliditySource(filename)

		if err != nil || strings.TrimSpace(soliditySource) == "" {
			return errors.New("no valid solidity contract source found")
		}

		finalSoliditySource, err := common.ApplyInitialWhitelistToSoliditySource(soliditySource, peers)*/
	finalSoliditySource, err := common.GetFinalSoliditySource(peers)

	if err != nil {
		message("Error building genesis contract:", err)
		return err
	}

	err = common.WriteToFile(filepath.Join(configDir, common.GenesisContract), finalSoliditySource)
	if err != nil {
		message("Error writing genesis contract:", err)
		return err
	}

	contractInfo, err := common.CompileSolidityContract(finalSoliditySource)
	if err != nil {
		message("Error compiling genesis contract:", err)
		return err
	}

	var poagenesis common.GenesisPOA

	// message("Contract Compiled: ", contractInfo)

	for k, v := range contractInfo {
		message("Processing Contract: ", k)
		jsonabi, err := json.MarshalIndent(v.Info.AbiDefinition, "", "\t")
		if err != nil {
			message("ABI error:", err)
			return err
		}

		tree.SetPath([]string{"poa", "contractclass"}, strings.TrimPrefix(k, "<stdin>:"))
		tree.SetPath([]string{"poa", "abi"}, string(jsonabi))

		common.WriteToFile(filepath.Join(configDir, common.GenesisABI), string(jsonabi))
		tree.SetPath([]string{"poa", "bytecode"}, strings.TrimPrefix(v.RuntimeCode, "0x"))

		poagenesis.Abi = string(jsonabi)
		poagenesis.Address = types.HexToAddress(tree.Get("poa.contractaddress").(string)).Hex() //EIP55 compliant
		poagenesis.Code = strings.TrimPrefix(v.RuntimeCode, "0x")

		message("Set Contract Items")
		break // We only have one contract ever so no need to loop. We use the for loop as k is indeterminate
	}

	err = common.SaveTomlConfig(configDir, tree)
	if err != nil {
		common.MessageWithType(common.MsgDebug, "Cannot save TOML file")
		return err
	}

	var genesis common.GenesisFile

	genesis.Alloc = &alloc
	genesis.Poa = &poagenesis

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	common.MessageWithType(common.MsgDebug, "Write Genesis.json")
	jsonFileName := filepath.Join(configDir, common.GenesisJSON)
	common.WriteToFile(jsonFileName, string(genesisjson))

	common.MessageWithType(common.MsgDebug, "Write Peers.json")

	peersjson, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}
	jsonFileName = filepath.Join(configDir, common.PeersJSON)
	common.WriteToFile(jsonFileName, string(peersjson))

	/*	peersjson, err = json.MarshalIndent(genesisPeers, "", "\t")
		if err != nil {
			return err
		}*/
	jsonFileName = filepath.Join(configDir, common.PeersGenesisJSON)
	common.WriteToFile(jsonFileName, string(peersjson))

	common.MessageWithType(common.MsgDebug, "Compilation Task Complete")

	return nil
}

/*
//CompileConfigWithParamb "finishes" the monetcli configuration, compiling the POA smart contract
//in preparation for a call to monetcli config publish
func CompileConfigWithParamb(configDir string) error {
	var soliditySource string
	// Load the Current Config

	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	// Retrieve and set the version number
	version, err := common.GetSolidityCompilerVersion()
	if err != nil {
		return err
	}

	tree.SetPath([]string{"poa", "compilerverison"}, version)

	//When contracts are "set" for a network, the solidity source is copied into the monetcli config directory
	//with a name of template.sol (defined by constant common.TemplateContract). Thus we can check just for that file.
	//If not found, then we download a fresh contract.
	filename := filepath.Join(configDir, common.TemplateContract)
	message("Checking for file: ", filename)

	if _, err := os.Stat(filename); err == nil {
		message("Opening: ", filename)
		file, err := os.Open(filename)
		if err != nil {
			message("Error opening: ", filename)
			return err
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			message("Error reading: ", filename)
			return err
		}

		soliditySource = string(b)
	} else { // NB, we do not write the downloaded template to file. Preferable to get fresh is regenerating.
		message("Loading: ", common.DefaultSolidityContract)
		resp, err := http.Get(common.DefaultSolidityContract)
		if err != nil {
			common.MessageWithType(common.MsgError, "Could not load the standard poa smart contract from GitHub. Aborting.")
			common.MessageWithType(common.MsgError, "You can specify the contract explicitly using the standard one from [...monetd]/smart-contract/genesis.sol")
			common.MessageWithType(common.MsgInformation, "monetcli network contract [...monetd]/smart-contract/genesis.sol")

			message("Error loading: ", common.DefaultSolidityContract)

			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			message("Error reading body of Solidity Contract")
			return err
		}

		soliditySource = string(body)
	}

	// message(soliditySource)

	currentNodes, err := GetPeersLabelsListFromToml(configDir)
	if err != nil {
		return err
	}

	if len(currentNodes) < 1 {
		return errors.New("Peerset is empty")
	}

	var consts, addTo, checks []string

	var alloc = make(common.GenesisAlloc)
	var peers common.PeerRecordList
	var genesisPeers common.PeerRecordList

	for i, value := range currentNodes {

		rawaddr := tree.GetPath([]string{"validators", value, "address"}).(string)
		rawmoniker := tree.GetPath([]string{"validators", value, "moniker"}).(string)
		rawpubkey := tree.GetPath([]string{"validators", value, "pubkey"}).(string)
		rawisvalidator := tree.GetPath([]string{"validators", value, "validator"}).(bool)
		rawip := tree.GetPath([]string{"validators", value, "ip"}).(string)

		// Convert Hex to Address and back out to get a EIP55 compliant address
		addr := types.HexToAddress(rawaddr).Hex()

		// Non-validators are added to the peer set, but not to the genesis peer set.
		peer := common.PeerRecord{NetAddr: rawip, PubKeyHex: rawpubkey, Moniker: rawmoniker}
		peers = append(peers, &peer)

		if rawisvalidator {
			consts = append(consts, "    address constant initWhitelist"+strconv.Itoa(i)+" = "+addr+";")
			consts = append(consts, "    bytes32 constant initWhitelistMoniker"+strconv.Itoa(i)+" = \""+rawmoniker+"\";")

			addTo = append(addTo, "     addToWhitelist(initWhitelist"+strconv.Itoa(i)+", initWhitelistMoniker"+strconv.Itoa(i)+");")
			checks = append(checks, " ( initWhitelist"+strconv.Itoa(i)+" == _address ) ")
			genesisPeers = append(genesisPeers, &peer)
		}

		rec := common.GenesisAllocRecord{Moniker: rawmoniker, Balance: common.DefaultAccountBalance}
		alloc[addr] = &rec

	}

	generatedSol := "GENERATED GENESIS BEGIN \n " +
		" \n" +
		strings.Join(consts, "\n") +
		" \n" +
		" \n" +
		" \n" +
		"    function processGenesisWhitelist() private \n" +
		"    { \n" +
		strings.Join(addTo, "\n") +
		" \n" +
		"    } \n" +
		" \n" +
		" \n" +
		"    function isGenesisWhitelisted(address _address) pure private returns (bool) \n" +
		"    { \n" +
		"        return ( " + strings.Join(checks, "||") + "); \n" +
		"    } \n" +

		" \n" +
		" //GENERATED GENESIS END \n "

	// replace

	reg := regexp.MustCompile(`(?s)GENERATED GENESIS BEGIN.*GENERATED GENESIS END`)
	finalSoliditySource := reg.ReplaceAllString(soliditySource, generatedSol)

	err = common.WriteToFile(filepath.Join(configDir, common.GenesisContract), finalSoliditySource)
	if err != nil {
		message("Error writing genesis contract:", err)
		return err
	}

	contractInfo, err := compile.CompileSolidityString("solc", finalSoliditySource)
	if err != nil {
		message("Error compiling genesis contract:", err)
		return err
	}

	var poagenesis common.GenesisPOA

	// message("Contract Compiled: ", contractInfo)

	for k, v := range contractInfo {
		message("Processing Contract: ", k)
		jsonabi, err := json.MarshalIndent(v.Info.AbiDefinition, "", "\t")
		if err != nil {
			message("ABI error:", err)
			return err
		}

		tree.SetPath([]string{"poa", "contractclass"}, strings.TrimPrefix(k, "<stdin>:"))
		tree.SetPath([]string{"poa", "abi"}, string(jsonabi))

		common.WriteToFile(filepath.Join(configDir, common.GenesisABI), string(jsonabi))
		tree.SetPath([]string{"poa", "bytecode"}, strings.TrimPrefix(v.RuntimeCode, "0x"))

		poagenesis.Abi = string(jsonabi)
		poagenesis.Address = types.HexToAddress(tree.Get("poa.contractaddress").(string)).Hex() //EIP55 compliant
		poagenesis.Code = strings.TrimPrefix(v.RuntimeCode, "0x")

		message("Set Contract Items")
		break // We only have one contract ever so no need to loop. We use the for loop as k is indeterminate
	}

	err = common.SaveTomlConfig(configDir, tree)
	if err != nil {
		common.MessageWithType(common.MsgDebug, "Cannot save TOML file")
		return err
	}

	var genesis common.GenesisFile

	genesis.Alloc = &alloc
	genesis.Poa = &poagenesis

	genesisjson, err := json.MarshalIndent(genesis, "", "\t")
	if err != nil {
		return err
	}

	common.MessageWithType(common.MsgDebug, "Write Genesis.json")
	jsonFileName := filepath.Join(configDir, common.GenesisJSON)
	common.WriteToFile(jsonFileName, string(genesisjson))

	common.MessageWithType(common.MsgDebug, "Write Peers.json")

	peersjson, err := json.MarshalIndent(peers, "", "\t")
	if err != nil {
		return err
	}
	jsonFileName = filepath.Join(configDir, common.PeersJSON)
	common.WriteToFile(jsonFileName, string(peersjson))

	peersjson, err = json.MarshalIndent(genesisPeers, "", "\t")
	if err != nil {
		return err
	}
	jsonFileName = filepath.Join(configDir, common.PeersGenesisJSON)
	common.WriteToFile(jsonFileName, string(peersjson))

	common.MessageWithType(common.MsgDebug, "Compilation Task Complete")

	return nil
}
*/
