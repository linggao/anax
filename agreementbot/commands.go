package agreementbot

import (
	"fmt"
	"github.com/open-horizon/anax/events"
	"github.com/open-horizon/anax/exchange"
	"github.com/open-horizon/anax/policy"
)

// ==============================================================================================================
type AgreementTimeoutCommand struct {
	AgreementId string
	Protocol    string
	Reason      uint
}

func (a AgreementTimeoutCommand) ShortString() string {
	return fmt.Sprintf("%v", a)
}

func NewAgreementTimeoutCommand(agreementId string, protocol string, reason uint) *AgreementTimeoutCommand {
	return &AgreementTimeoutCommand{
		AgreementId: agreementId,
		Protocol:    protocol,
		Reason:      reason,
	}
}

// ==============================================================================================================
type PolicyChangedCommand struct {
	Msg events.PolicyChangedMessage
}

func (p PolicyChangedCommand) ShortString() string {
	return fmt.Sprintf("%v", p)
}

func NewPolicyChangedCommand(msg events.PolicyChangedMessage) *PolicyChangedCommand {
	return &PolicyChangedCommand{
		Msg: msg,
	}
}

// ==============================================================================================================
type PolicyDeletedCommand struct {
	Msg events.PolicyDeletedMessage
}

func (p PolicyDeletedCommand) ShortString() string {
	return fmt.Sprintf("%v", p)
}

func NewPolicyDeletedCommand(msg events.PolicyDeletedMessage) *PolicyDeletedCommand {
	return &PolicyDeletedCommand{
		Msg: msg,
	}
}

// ==============================================================================================================
type NewProtocolMessageCommand struct {
	Message   []byte
	MessageId int
	From      string
	PubKey    []byte
}

func (p NewProtocolMessageCommand) ShortString() string {
	return fmt.Sprintf("%v", p)
}

func NewNewProtocolMessageCommand(msg []byte, msgId int, deviceId string, pubkey []byte) *NewProtocolMessageCommand {
	return &NewProtocolMessageCommand{
		Message:   msg,
		MessageId: msgId,
		From:      deviceId,
		PubKey:    pubkey,
	}
}

// ==============================================================================================================
type BlockchainEventCommand struct {
	Msg events.EthBlockchainEventMessage
}

func (e BlockchainEventCommand) ShortString() string {
	return e.Msg.ShortString()
}

func NewBlockchainEventCommand(msg events.EthBlockchainEventMessage) *BlockchainEventCommand {
	return &BlockchainEventCommand{
		Msg: msg,
	}
}

// ==============================================================================================================
type WorkloadUpgradeCommand struct {
	Msg events.ABApiWorkloadUpgradeMessage
}

func (e WorkloadUpgradeCommand) ShortString() string {
	return e.Msg.ShortString()
}

func NewWorkloadUpgradeCommand(msg events.ABApiWorkloadUpgradeMessage) *WorkloadUpgradeCommand {
	return &WorkloadUpgradeCommand{
		Msg: msg,
	}
}

// ==============================================================================================================
type MakeAgreementCommand struct {
	ProducerPolicy policy.Policy               // the producer policy received from the exchange
	ConsumerPolicy policy.Policy               // the consumer policy we're matched up with
	Org            string                      // the org of the consumer policy
	Device         exchange.SearchResultDevice // the device entry in the exchange
}

func (e MakeAgreementCommand) ShortString() string {
	return fmt.Sprintf("Produder Policy: %v, ConsumerPolicy: %v, Org: %v, Device: %v", e.ProducerPolicy.Header.Name, e.ConsumerPolicy.Header.Name, e.Org, e.Device)
}

func NewMakeAgreementCommand(pPol policy.Policy, cPol policy.Policy, org string, dev exchange.SearchResultDevice) *MakeAgreementCommand {
	return &MakeAgreementCommand{
		ProducerPolicy: pPol,
		ConsumerPolicy: cPol,
		Org:            org,
		Device:         dev,
	}
}

// ==============================================================================================================
type ClientInitializedCommand struct {
	Msg events.BlockchainClientInitializedMessage
}

func (e ClientInitializedCommand) ShortString() string {
	return e.Msg.ShortString()
}

func NewClientInitializedCommand(msg *events.BlockchainClientInitializedMessage) *ClientInitializedCommand {
	return &ClientInitializedCommand{
		Msg: *msg,
	}
}

// ==============================================================================================================
type ClientStoppingCommand struct {
	Msg events.BlockchainClientStoppingMessage
}

func (e ClientStoppingCommand) ShortString() string {
	return e.Msg.ShortString()
}

func NewClientStoppingCommand(msg *events.BlockchainClientStoppingMessage) *ClientStoppingCommand {
	return &ClientStoppingCommand{
		Msg: *msg,
	}
}

// ==============================================================================================================
type AccountFundedCommand struct {
	Msg events.AccountFundedMessage
}

func (e AccountFundedCommand) ShortString() string {
	return e.Msg.ShortString()
}

func NewAccountFundedCommand(msg *events.AccountFundedMessage) *AccountFundedCommand {
	return &AccountFundedCommand{
		Msg: *msg,
	}
}
