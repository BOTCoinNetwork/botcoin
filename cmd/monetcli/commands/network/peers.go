package network

import (
	"fmt"

	"github.com/mosaicnetworks/monetd/src/common"
)

func PeersWizard(configDir string) error {

	var err error
	common.Message("Peers")

	tree, err := common.LoadTomlConfig(configDir)
	if err != nil {
		return err
	}

	peers := GetPeersLabelsList(tree)
	peers = append(peers, common.WizardExit)

peerloop:
	for {
		selectedPeer := RequestSelect("Choose a peer", peers, "")
		if selectedPeer == common.WizardExit {
			break peerloop
		}
	actionloop:
		for {
			action := RequestSelect("Choose an Action: ", []string{common.WizardView, common.WizardEdit, common.WizardDelete, common.WizardExit}, "")

			switch action {
			case common.WizardView:
				err = viewPeer(configDir, action)
			case common.WizardEdit:
				err = editPeer(configDir, action)
			case common.WizardDelete:
				err = deletePeer(configDir, action)
			case common.WizardExit:
				break actionloop

			}

			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("%v \n", peers)

	return nil
}

func viewPeer(configDir string, peername string) error {
	return nil
}

func editPeer(configDir string, peername string) error {
	return nil
}

func deletePeer(configDir string, peername string) error {
	return nil
}
