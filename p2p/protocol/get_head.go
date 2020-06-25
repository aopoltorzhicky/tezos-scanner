package protocol

import (
	"fmt"
	"time"
)

func (peer *Peer) getHead() (head CurrentHeadMsg, err error) {
	// Receiving first maeesage from node GetCurrentBranch
	msg, msgType, err := peer.ReceivePeerMessage()
	if err != nil {
		return head, err
	}

	if msgType != GetCurrentBranchTag {
		return head, fmt.Errorf("message has a different type: %d", msgType)
	}

	currentBranch := msg.(GetCurrentBranchMsg).Branch
	// Response on node request CurrentBranchMsg
	if err = peer.SendMessage(CurrentHeadMsg{ChainID: currentBranch}); err != nil {
		return
	}
	if _, _, err = peer.ReceivePeerMessage(); err != nil {
		return
	}

	// Request current head
	if err = peer.SendMessage(GetCurrentHeadMsg{ChainID: currentBranch}); err != nil {
		return
	}
	// Receiving current head
	msg, msgType, err = peer.ReceivePeerMessage()
	if err != nil {
		return
	}
	if msgType != CurrentHeadTag {
		return head, fmt.Errorf("message has a different type: %d", msgType)
	}
	head = msg.(CurrentHeadMsg)
	return
}

// UpdateSyncState -
func (peer *Peer) UpdateSyncState(syncedTime int64) error {
	currentHead, err := peer.getHead()
	if err != nil {
		return err
	}
	diff := time.Now().Unix() - currentHead.CurrentBlockHeader.Timestamp
	peer.Synced = diff < syncedTime
	return nil
}
