package events

import (
	"fmt"
	"github.com/open-horizon/anax/containermessage"
	"github.com/open-horizon/anax/persistence"
	"net/url"
	"time"
)

type Event struct {
	Id EventId
}

func (e Event) String() string {
	return fmt.Sprintf("%v", e.Id)
}

type EventId string

// event constants are declared here for all workers to ensure uniqueness of constant values
const (
	// blockchain-related
	NOOP                  EventId = "NOOP"
	AGREEMENT_ACCEPTED    EventId = "AGREEMENT_ACCEPTED"
	AGREEMENT_ENDED       EventId = "AGREEMENT_ENDED"
	AGREEMENT_CREATED     EventId = "AGREEMENT_CREATED"
	AGREEMENT_REGISTERED  EventId = "AGREEMENT_REGISTERED"
	ACCOUNT_FUNDED        EventId = "ACCOUNT_FUNDED"
	BC_CLIENT_INITIALIZED EventId = "BC_CLIENT_INITIALIZED"
	BC_CLIENT_STOPPING    EventId = "BC_CLIENT_STOPPING"
	BC_EVENT              EventId = "BC_EVENT"
	BC_NEEDED             EventId = "BC_NEEDED"
	ALL_STOP              EventId = "ALL_STOP"

	// exchange related
	RECEIVED_EXCHANGE_DEV_MSG EventId = "RECEIVED_EXCHANGE_DEV_MSG"

	// image fetching related
	IMAGE_FETCHED          EventId = "IMAGE_FETCHED"
	IMAGE_DATA_ERROR       EventId = "IMAGE_DATA_ERROR"
	IMAGE_FETCH_ERROR      EventId = "IMAGE_FETCH_ERROR"
	IMAGE_FETCH_AUTH_ERROR EventId = "IMAGE_FETCH_AUTH_ERROR"
	IMAGE_SIG_VERIF_ERROR  EventId = "IMAGE_SIG_VERIF_ERROR"

	// container-related
	EXECUTION_FAILED    EventId = "EXECUTION_FAILED"
	EXECUTION_BEGUN     EventId = "EXECUTION_BEGUN"
	WORKLOAD_DESTROYED  EventId = "WORKLOAD_DESTROYED"
	CONTAINER_STOPPING  EventId = "CONTAINER_STOPPING"
	CONTAINER_DESTROYED EventId = "CONTAINER_DESTROYED"
	CONTAINER_MAINTAIN  EventId = "CONTAINER_MAINTAIN"
	LOAD_CONTAINER      EventId = "LOAD_CONTAINER"
	START_MICROSERVICE  EventId = "START_MICROSERVICE"
	CANCEL_MICROSERVICE EventId = "CANCEL_MICROSERVICE"
	NEW_BC_CLIENT       EventId = "NEW_BC_CONTAINER"
	IMAGE_LOAD_FAILED   EventId = "IMAGE_LOAD_FAILED"

	// policy-related
	NEW_POLICY     EventId = "NEW_POLICY"
	CHANGED_POLICY EventId = "CHANGED_POLICY"
	DELETED_POLICY EventId = "DELETED_POLICY"

	// exchange-related
	NEW_DEVICE_REG             EventId = "NEW_DEVICE_REG"
	NEW_DEVICE_CONFIG_COMPLETE EventId = "NEW_DEVICE_CONFIG_COMPLETE"
	NEW_AGBOT_REG              EventId = "NEW_AGBOT_REG"

	// agreement-related
	AGREEMENT_REACHED        EventId = "AGREEMENT_REACHED"
	DEVICE_AGREEMENTS_SYNCED EventId = "DEVICE_AGREEMENTS_SYNCED"
	DEVICE_CONTAINERS_SYNCED EventId = "DEVICE_CONTAINERS_SYNCED"
	WORKLOAD_UPGRADE         EventId = "WORKLOAD_UPGRADE"

	// Node related
	START_UNCONFIGURE       EventId = "UNCONFIGURE_NODE"
	UNCONFIGURE_COMPLETE    EventId = "UNCONFIGURE_COMPLETE"
	WORKER_STOP             EventId = "WORKER_STOP"
	START_AGBOT_QUIESCE     EventId = "AGBOT_QUIESCE"
	AGBOT_QUIESCE_COMPLETE  EventId = "AGBOT_QUIESCE_COMPLETE"
	NODE_HEARTBEAT_FAILED   EventId = "HEARTBEAT_FAILED"
	NODE_HEARTBEAT_RESTORED EventId = "HEARTBEAT_RESTORED"
)

type EndContractCause string

const (
	AG_TERMINATED EndContractCause = "AG_TERMINATED"
	AG_ERROR      EndContractCause = "AG_ERROR"
	AG_FULFILLED  EndContractCause = "AG_FULFILLED"
)

type Message interface {
	Event() Event
	ShortString() string
}

type LaunchContext interface {
	ContainerConfig() ContainerConfig
	ShortString() string
}

type MicroserviceSpec struct {
	SpecRef string
	Org     string
	Version string
	MsdefId string
}

type AgreementLaunchContext struct {
	AgreementProtocol    string
	AgreementId          string
	Configure            ContainerConfig
	ConfigureRaw         []byte
	EnvironmentAdditions *map[string]string // provided by platform, not but user
	Microservices        []MicroserviceSpec // for ms split.
}

func (c AgreementLaunchContext) String() string {
	return fmt.Sprintf("AgreementProtocol: %v, AgreementId: %v, Configure: %v, EnvironmentAdditions: %v, Microservices: %v", c.AgreementProtocol, c.AgreementId, c.Configure, c.EnvironmentAdditions, c.Microservices)
}

func (c AgreementLaunchContext) ShortString() string {
	return fmt.Sprintf("AgreementProtocol: %v, AgreementId: %v", c.AgreementProtocol, c.AgreementId)
}

func (c AgreementLaunchContext) ContainerConfig() ContainerConfig {
	return c.Configure
}

type ImageDockerAuth struct {
	Registry string `json:"registry"`
	UserName string `json:"username"`
	Password string `json:"token"`
}

func (s ImageDockerAuth) String() string {
	return fmt.Sprintf(
		"Registry: %v, "+
			"UserName: %v, "+
			"Password: %v",
		s.Registry, s.UserName, "******")
}

type ContainerConfig struct {
	TorrentURL          url.URL           `json:"torrent_url"`
	TorrentSignature    string            `json:"torrent_signature"`
	Deployment          string            `json:"deployment"`           // A stringified (and escaped) JSON structure.
	DeploymentSignature string            `json:"deployment_signature"` // Digital signature of the Deployment string.
	DeploymentUserInfo  string            `json:"deployment_user_info"`
	Overrides           string            `json:"overrides"`
	ImageDockerAuths    []ImageDockerAuth `json:"image_auths"`
}

func (c ContainerConfig) String() string {
	return fmt.Sprintf("TorrentURL: %v, TorrentSignature: %v, Deployment: %v, DeploymentSignature: %v, DeploymentUserInfo: %v, Overrides: %v, ImageDockerAuths: %v", c.TorrentURL.String(), c.TorrentSignature, c.Deployment, c.DeploymentSignature, c.DeploymentUserInfo, c.Overrides, c.ImageDockerAuths)
}

func NewContainerConfig(torrentURL url.URL, torrentSignature string, deployment string, deploymentSignature string, deploymentUserInfo string, overrides string, imageDockerAuths []ImageDockerAuth) *ContainerConfig {
	return &ContainerConfig{
		TorrentURL:          torrentURL,
		TorrentSignature:    torrentSignature,
		Deployment:          deployment,
		DeploymentSignature: deploymentSignature,
		DeploymentUserInfo:  deploymentUserInfo,
		Overrides:           overrides,
		ImageDockerAuths:    imageDockerAuths,
	}
}

type BlockchainConfig struct {
	Type string
	Name string
	Org  string
}

type ContainerLaunchContext struct {
	Configure            ContainerConfig
	EnvironmentAdditions *map[string]string
	Blockchain           BlockchainConfig
	Name                 string // used as the docker network name and part of container name. For microservice it is the ms instance key
	AgreementId          string
	Microservices        []MicroserviceSpec                     // Service dependencies go here. Microservices (in the workload/microservice model) never have dependencies.
	ServicePathElement   persistence.ServiceInstancePathElement // The service that we're trying to start.
}

func (c ContainerLaunchContext) String() string {
	return fmt.Sprintf("ContainerConfig: %v, EnvironmentAdditions: %v, Blockchain: %v, Name: %v, AgreementId: %v, ServiceDependencies: %v, ThisService: %v", c.Configure, c.EnvironmentAdditions, c.Blockchain, c.Name, c.AgreementId, c.Microservices, c.ServicePathElement)
}

func (c ContainerLaunchContext) ShortString() string {
	return c.String()
}

func (c ContainerLaunchContext) ContainerConfig() ContainerConfig {
	return c.Configure
}

func (c ContainerLaunchContext) GetAgreementId() string {
	return c.AgreementId
}

func (c ContainerLaunchContext) GetMicroservices() []MicroserviceSpec {
	return c.Microservices
}

func (c ContainerLaunchContext) GetServicePathElement() *persistence.ServiceInstancePathElement {
	return &c.ServicePathElement
}

func NewContainerLaunchContext(config *ContainerConfig, envAdds *map[string]string, bc BlockchainConfig, name string, agId string, mss []MicroserviceSpec, spe *persistence.ServiceInstancePathElement) *ContainerLaunchContext {

	spe_temp := spe
	if spe_temp == nil {
		spe_temp = persistence.NewServiceInstancePathElement("", "", "")
	}

	return &ContainerLaunchContext{
		Configure:            *config,
		EnvironmentAdditions: envAdds,
		Blockchain:           bc,
		Name:                 name,
		AgreementId:          agId,
		Microservices:        mss,
		ServicePathElement:   *spe_temp,
	}
}

// Anax device side fires this event when it needs to download and load a container.
type LoadContainerMessage struct {
	event         Event
	launchContext *ContainerLaunchContext
}

func (e LoadContainerMessage) String() string {
	return fmt.Sprintf("event: %v, launch context: %v", e.event, e.launchContext)
}

func (e LoadContainerMessage) ShortString() string {
	return e.String()
}

func (e *LoadContainerMessage) Event() Event {
	return e.event
}

func (e *LoadContainerMessage) LaunchContext() *ContainerLaunchContext {
	return e.launchContext
}

func NewLoadContainerMessage(id EventId, lc *ContainerLaunchContext) *LoadContainerMessage {

	return &LoadContainerMessage{
		event: Event{
			Id: id,
		},
		launchContext: lc,
	}
}

// This event indicates that a new microservice has been created in the form of a policy file
type PolicyCreatedMessage struct {
	event    Event
	fileName string
}

func (e PolicyCreatedMessage) String() string {
	return fmt.Sprintf("event: %v, file: %v", e.event, e.fileName)
}

func (e PolicyCreatedMessage) ShortString() string {
	return e.String()
}

func (e PolicyCreatedMessage) Event() Event {
	return e.event
}

func (e *PolicyCreatedMessage) PolicyFile() string {
	return e.fileName
}

func NewPolicyCreatedMessage(id EventId, policyFileName string) *PolicyCreatedMessage {

	return &PolicyCreatedMessage{
		event: Event{
			Id: id,
		},
		fileName: policyFileName,
	}
}

// This event indicates that a policy file has changed. It might also be a new policy file in the agbot.
type PolicyChangedMessage struct {
	event    Event
	fileName string
	name     string
	policy   string
	org      string
}

func (e PolicyChangedMessage) String() string {
	return fmt.Sprintf("event: %v, file: %v, name: %v, org: %v, policy: %v", e.event, e.fileName, e.name, e.org, e.policy)
}

func (e PolicyChangedMessage) ShortString() string {
	return e.String()
}

func (e *PolicyChangedMessage) Event() Event {
	return e.event
}

func (e *PolicyChangedMessage) PolicyFile() string {
	return e.fileName
}

func (e *PolicyChangedMessage) PolicyName() string {
	return e.name
}

func (e *PolicyChangedMessage) Org() string {
	return e.org
}

func (e *PolicyChangedMessage) PolicyString() string {
	return e.policy
}

func NewPolicyChangedMessage(id EventId, policyFileName string, policyName string, org string, policy string) *PolicyChangedMessage {

	return &PolicyChangedMessage{
		event: Event{
			Id: id,
		},
		fileName: policyFileName,
		name:     policyName,
		policy:   policy,
		org:      org,
	}
}

// This event indicates that a policy file was deleted.
type PolicyDeletedMessage struct {
	event    Event
	fileName string
	name     string
	policy   string
	org      string
}

func (e PolicyDeletedMessage) String() string {
	return fmt.Sprintf("event: %v, file: %v, name: %v, org: %v, policy: %v", e.event, e.fileName, e.name, e.org, e.policy)
}

func (e PolicyDeletedMessage) ShortString() string {
	return e.String()
}

func (e *PolicyDeletedMessage) Event() Event {
	return e.event
}

func (e *PolicyDeletedMessage) PolicyFile() string {
	return e.fileName
}

func (e *PolicyDeletedMessage) PolicyName() string {
	return e.name
}

func (e *PolicyDeletedMessage) Org() string {
	return e.org
}

func (e *PolicyDeletedMessage) PolicyString() string {
	return e.policy
}

func NewPolicyDeletedMessage(id EventId, policyFileName string, policyName string, org string, policy string) *PolicyDeletedMessage {

	return &PolicyDeletedMessage{
		event: Event{
			Id: id,
		},
		fileName: policyFileName,
		name:     policyName,
		policy:   policy,
		org:      org,
	}
}

// This event indicates that the edge device has been registered in the exchange
type EdgeRegisteredExchangeMessage struct {
	event     Event
	device_id string
	token     string
	org       string
	pattern   string
}

func (e EdgeRegisteredExchangeMessage) String() string {
	return fmt.Sprintf("event: %v, device_id: %v, token: %v, org: %v, pattern: %v", e.event, e.device_id, e.token, e.org, e.pattern)
}

func (e EdgeRegisteredExchangeMessage) ShortString() string {
	return e.String()
}

func (e *EdgeRegisteredExchangeMessage) Event() Event {
	return e.event
}

func (e *EdgeRegisteredExchangeMessage) DeviceId() string {
	return e.device_id
}

func (e *EdgeRegisteredExchangeMessage) Token() string {
	return e.token
}

func (e *EdgeRegisteredExchangeMessage) Org() string {
	return e.org
}

func (e *EdgeRegisteredExchangeMessage) Pattern() string {
	return e.pattern
}

func NewEdgeRegisteredExchangeMessage(evId EventId, device_id string, token string, org string, pattern string) *EdgeRegisteredExchangeMessage {

	return &EdgeRegisteredExchangeMessage{
		event: Event{
			Id: evId,
		},
		device_id: device_id,
		token:     token,
		org:       org,
		pattern:   pattern,
	}
}

// This event indicates that the edge device configuration is complete
type EdgeConfigCompleteMessage struct {
	event Event
}

func (e EdgeConfigCompleteMessage) String() string {
	return fmt.Sprintf("event: %v", e.event)
}

func (e EdgeConfigCompleteMessage) ShortString() string {
	return e.String()
}

func (e *EdgeConfigCompleteMessage) Event() Event {
	return e.event
}

func NewEdgeConfigCompleteMessage(evId EventId) *EdgeConfigCompleteMessage {

	return &EdgeConfigCompleteMessage{
		event: Event{
			Id: evId,
		},
	}
}

// Anax device side fires this event when an agreement is reached so that it can begin
// downloading containers. The Agreement is not final until it is seen in the blockchain.
type AgreementReachedMessage struct {
	event         Event
	launchContext *AgreementLaunchContext
}

func (e AgreementReachedMessage) String() string {
	return fmt.Sprintf("event: %v, launch context: %v", e.event, e.launchContext)
}

func (e AgreementReachedMessage) ShortString() string {
	return fmt.Sprintf("event: %v, launch context: %v", e.event, e.launchContext.ShortString())
}

func (e *AgreementReachedMessage) Event() Event {
	return e.event
}

func (e *AgreementReachedMessage) LaunchContext() *AgreementLaunchContext {
	return e.launchContext
}

func NewAgreementMessage(id EventId, lc *AgreementLaunchContext) *AgreementReachedMessage {

	return &AgreementReachedMessage{
		event: Event{
			Id: id,
		},
		launchContext: lc,
	}
}

type TorrentMessage struct {
	event                 Event
	DeploymentDescription *containermessage.DeploymentDescription
	LaunchContext         interface{}
}

// fulfill interface of events.Message
func (b *TorrentMessage) Event() Event {
	return b.event
}

func (b *TorrentMessage) String() string {
	return fmt.Sprintf("event: %v, deploymentDescription: %v, launchContext: %v", b.event, b.DeploymentDescription, b.LaunchContext)
}

func (b *TorrentMessage) ShortString() string {
	return fmt.Sprintf("event: %v, deploymentDescription: %v, launchContext: %v", b.event, b.DeploymentDescription, b.LaunchContext)
}

func NewTorrentMessage(id EventId, deploymentDescription *containermessage.DeploymentDescription, launchContext interface{}) *TorrentMessage {

	return &TorrentMessage{
		event: Event{
			Id: id,
		},
		DeploymentDescription: deploymentDescription,
		LaunchContext:         launchContext,
	}
}

// Governance messages
type GovernanceMaintenanceMessage struct {
	event             Event
	AgreementProtocol string
	AgreementId       string
	Deployment        persistence.DeploymentConfig
}

func (m *GovernanceMaintenanceMessage) Event() Event {
	return m.event
}

func (m GovernanceMaintenanceMessage) String() string {
	depStr := ""
	if m.Deployment != nil {
		depStr = m.Deployment.ToString()
	}
	return fmt.Sprintf("Event: %v, AgreementProtocol: %v, AgreementId: %v, Deployment: %v", m.event, m.AgreementProtocol, m.AgreementId, depStr)
}

func (m GovernanceMaintenanceMessage) ShortString() string {
	return m.String()
}

func NewGovernanceMaintenanceMessage(id EventId, protocol string, agreementId string, deployment persistence.DeploymentConfig) *GovernanceMaintenanceMessage {
	return &GovernanceMaintenanceMessage{
		event: Event{
			Id: id,
		},
		AgreementProtocol: protocol,
		AgreementId:       agreementId,
		Deployment:        deployment,
	}
}

type GovernanceWorkloadCancelationMessage struct {
	GovernanceMaintenanceMessage
	Message
	Cause EndContractCause
}

func (m *GovernanceWorkloadCancelationMessage) Event() Event {
	return m.event
}

func (m GovernanceWorkloadCancelationMessage) String() string {
	depStr := ""
	if m.Deployment != nil {
		depStr = m.Deployment.ToString()
	}
	return fmt.Sprintf("Event: %v, AgreementProtocol: %v, AgreementId: %v, Deployment: %v, Cause: %v", m.event, m.AgreementProtocol, m.AgreementId, depStr, m.Cause)
}

func (m GovernanceWorkloadCancelationMessage) ShortString() string {
	return m.String()
}

func NewGovernanceWorkloadCancelationMessage(id EventId, cause EndContractCause, protocol string, agreementId string, deployment persistence.DeploymentConfig) *GovernanceWorkloadCancelationMessage {

	govMaint := NewGovernanceMaintenanceMessage(id, protocol, agreementId, deployment)

	return &GovernanceWorkloadCancelationMessage{
		GovernanceMaintenanceMessage: *govMaint,
		Cause: cause,
	}
}

//Workload messages
type WorkloadMessage struct {
	event             Event
	AgreementProtocol string
	AgreementId       string
	Deployment        persistence.DeploymentConfig
}

func (m WorkloadMessage) String() string {
	depStr := ""
	if m.Deployment != nil {
		depStr = m.Deployment.ToString()
	}
	return fmt.Sprintf("event: %v, AgreementProtocol: %v, AgreementId: %v, Deployment: %v", m.event.Id, m.AgreementProtocol, m.AgreementId, depStr)
}

func (m WorkloadMessage) ShortString() string {
	return m.String()
}

func (b WorkloadMessage) Event() Event {
	return b.event
}

func NewWorkloadMessage(id EventId, protocol string, agreementId string, deployment persistence.DeploymentConfig) *WorkloadMessage {

	return &WorkloadMessage{
		event: Event{
			Id: id,
		},
		AgreementProtocol: protocol,
		AgreementId:       agreementId,
		Deployment:        deployment,
	}
}

//Container messages
type ContainerMessage struct {
	event         Event
	LaunchContext ContainerLaunchContext
	ServiceName   string
	ServicePort   string
}

func (m ContainerMessage) String() string {
	return fmt.Sprintf("event: %v, ServiceName: %v, ServicePort: %v, LaunchContext: %v", m.event.Id, m.ServiceName, m.ServicePort, m.LaunchContext)
}

func (m ContainerMessage) ShortString() string {
	return m.String()
}

func (b ContainerMessage) Event() Event {
	return b.event
}

func NewContainerMessage(id EventId, lc ContainerLaunchContext, serviceName string, servicePort string) *ContainerMessage {

	return &ContainerMessage{
		event: Event{
			Id: id,
		},
		LaunchContext: lc,
		ServiceName:   serviceName,
		ServicePort:   servicePort,
	}
}

//Container stop message
type ContainerStopMessage struct {
	event         Event
	ContainerName string
	Org           string
}

func (m ContainerStopMessage) String() string {
	return fmt.Sprintf("event: %v, ContainerName: %v, Org: %v", m.event.Id, m.ContainerName, m.Org)
}

func (m ContainerStopMessage) ShortString() string {
	return m.String()
}

func (b ContainerStopMessage) Event() Event {
	return b.event
}

func NewContainerStopMessage(id EventId, containerName string, org string) *ContainerStopMessage {

	return &ContainerStopMessage{
		event: Event{
			Id: id,
		},
		ContainerName: containerName,
		Org:           org,
	}
}

//Container Shutdown message
type ContainerShutdownMessage struct {
	event         Event
	ContainerName string
	Org           string
}

func (m ContainerShutdownMessage) String() string {
	return fmt.Sprintf("event: %v, ContainerName: %v, Org: %v", m.event.Id, m.ContainerName, m.Org)
}

func (m ContainerShutdownMessage) ShortString() string {
	return m.String()
}

func (b ContainerShutdownMessage) Event() Event {
	return b.event
}

func NewContainerShutdownMessage(id EventId, containerName string, org string) *ContainerShutdownMessage {

	return &ContainerShutdownMessage{
		event: Event{
			Id: id,
		},
		ContainerName: containerName,
		Org:           org,
	}
}

// Api messages
type ApiAgreementCancelationMessage struct {
	event             Event
	AgreementProtocol string
	AgreementId       string
	Deployment        persistence.DeploymentConfig
	Cause             EndContractCause
}

func (m *ApiAgreementCancelationMessage) Event() Event {
	return m.event
}

func (m ApiAgreementCancelationMessage) String() string {
	depStr := ""
	if m.Deployment != nil {
		depStr = m.Deployment.ToString()
	}
	return fmt.Sprintf("Event: %v, AgreementProtocol: %v, AgreementId: %v, Deployment: %v, Cause: %v", m.event, m.AgreementProtocol, m.AgreementId, depStr, m.Cause)
}

func (m ApiAgreementCancelationMessage) ShortString() string {
	return m.String()
}

func NewApiAgreementCancelationMessage(id EventId, cause EndContractCause, protocol string, agreementId string, deployment persistence.DeploymentConfig) *ApiAgreementCancelationMessage {
	return &ApiAgreementCancelationMessage{
		event: Event{
			Id: id,
		},
		AgreementProtocol: protocol,
		AgreementId:       agreementId,
		Deployment:        deployment,
		Cause:             cause,
	}
}

// Agbot Api messages
type ABApiAgreementCancelationMessage struct {
	event             Event
	AgreementProtocol string
	AgreementId       string
}

func (m *ABApiAgreementCancelationMessage) Event() Event {
	return m.event
}

func (m ABApiAgreementCancelationMessage) String() string {
	return fmt.Sprintf("Event: %v, AgreementProtocol: %v, AgreementId: %v", m.event, m.AgreementProtocol, m.AgreementId)
}

func (m ABApiAgreementCancelationMessage) ShortString() string {
	return m.String()
}

func NewABApiAgreementCancelationMessage(id EventId, protocol string, agreementId string) *ABApiAgreementCancelationMessage {
	return &ABApiAgreementCancelationMessage{
		event: Event{
			Id: id,
		},
		AgreementProtocol: protocol,
		AgreementId:       agreementId,
	}
}

type ABApiWorkloadUpgradeMessage struct {
	event             Event
	AgreementProtocol string
	AgreementId       string
	DeviceId          string
	PolicyName        string
}

func (m *ABApiWorkloadUpgradeMessage) Event() Event {
	return m.event
}

func (m ABApiWorkloadUpgradeMessage) String() string {
	return fmt.Sprintf("Event: %v, AgreementProtocol: %v, AgreementId: %v, DeviceId: %v, PolicyName: %v", m.event, m.AgreementProtocol, m.AgreementId, m.DeviceId, m.PolicyName)
}

func (m ABApiWorkloadUpgradeMessage) ShortString() string {
	return m.String()
}

func NewABApiWorkloadUpgradeMessage(id EventId, protocol string, agreementId string, deviceId string, policyName string) *ABApiWorkloadUpgradeMessage {
	return &ABApiWorkloadUpgradeMessage{
		event: Event{
			Id: id,
		},
		AgreementProtocol: protocol,
		AgreementId:       agreementId,
		DeviceId:          deviceId,
		PolicyName:        policyName,
	}
}

// Initialization and restart messages
type InitAgreementCancelationMessage struct {
	event             Event
	AgreementProtocol string
	AgreementId       string
	Deployment        persistence.DeploymentConfig
	Reason            uint
}

func (m *InitAgreementCancelationMessage) Event() Event {
	return m.event
}

func (m InitAgreementCancelationMessage) String() string {
	depStr := ""
	if m.Deployment != nil {
		depStr = m.Deployment.ToString()
	}
	return fmt.Sprintf("Event: %v, AgreementProtocol: %v, AgreementId: %v, Deployment: %v, Reason: %v", m.event, m.AgreementProtocol, m.AgreementId, depStr, m.Reason)
}

func (m InitAgreementCancelationMessage) ShortString() string {
	return m.String()
}

func NewInitAgreementCancelationMessage(id EventId, reason uint, protocol string, agreementId string, deployment persistence.DeploymentConfig) *InitAgreementCancelationMessage {
	return &InitAgreementCancelationMessage{
		event: Event{
			Id: id,
		},
		AgreementProtocol: protocol,
		AgreementId:       agreementId,
		Deployment:        deployment,
		Reason:            reason,
	}
}

// Account funded message
type AccountFundedMessage struct {
	event       Event
	Account     string
	Time        uint64
	bcType      string
	bcInstance  string
	bcOrg       string
	serviceName string
	servicePort string
	colonusDir  string
}

func (m *AccountFundedMessage) Event() Event {
	return m.event
}

func (m AccountFundedMessage) String() string {
	return fmt.Sprintf("Event: %v, Account: %v, Time: %v, Type: %v, Instance: %v, Org: %v, ServiceName: %v, ServicePort: %v, ColonusDir: %v", m.event, m.Account, m.Time, m.bcType, m.bcInstance, m.bcOrg, m.serviceName, m.servicePort, m.colonusDir)
}

func (m AccountFundedMessage) ShortString() string {
	return m.String()
}

func (m AccountFundedMessage) BlockchainType() string {
	return m.bcType
}

func (m AccountFundedMessage) BlockchainInstance() string {
	return m.bcInstance
}

func (m AccountFundedMessage) BlockchainOrg() string {
	return m.bcOrg
}

func (m AccountFundedMessage) ServiceName() string {
	return m.serviceName
}

func (m AccountFundedMessage) ServicePort() string {
	return m.servicePort
}

func (m AccountFundedMessage) ColonusDir() string {
	return m.colonusDir
}

func NewAccountFundedMessage(id EventId, acct string, bcType string, bcName string, bcOrg string, serviceName string, servicePort string, colonusDir string) *AccountFundedMessage {
	return &AccountFundedMessage{
		event: Event{
			Id: id,
		},
		Account:     acct,
		Time:        uint64(time.Now().Unix()),
		bcType:      bcType,
		bcInstance:  bcName,
		bcOrg:       bcOrg,
		serviceName: serviceName,
		servicePort: servicePort,
		colonusDir:  colonusDir,
	}
}

// Blockchain client initialized message
type BlockchainClientInitializedMessage struct {
	event       Event
	Time        uint64
	bcType      string
	bcInstance  string
	bcOrg       string
	serviceName string
	servicePort string
	colonusDir  string
}

func (m *BlockchainClientInitializedMessage) Event() Event {
	return m.event
}

func (m BlockchainClientInitializedMessage) String() string {
	return fmt.Sprintf("Event: %v, Time: %v, Type: %v, Instance: %v, Org: %v, ServiceName: %v, ServicePort: %v, ColonusDir: %v", m.event, m.Time, m.bcType, m.bcInstance, m.bcOrg, m.serviceName, m.servicePort, m.colonusDir)
}

func (m BlockchainClientInitializedMessage) ShortString() string {
	return m.String()
}

func (m BlockchainClientInitializedMessage) BlockchainType() string {
	return m.bcType
}

func (m BlockchainClientInitializedMessage) BlockchainInstance() string {
	return m.bcInstance
}

func (m BlockchainClientInitializedMessage) BlockchainOrg() string {
	return m.bcOrg
}

func (m BlockchainClientInitializedMessage) ServiceName() string {
	return m.serviceName
}

func (m BlockchainClientInitializedMessage) ServicePort() string {
	return m.servicePort
}

func (m BlockchainClientInitializedMessage) ColonusDir() string {
	return m.colonusDir
}

func NewBlockchainClientInitializedMessage(id EventId, bcType string, bcName string, bcOrg string, serviceName string, servicePort string, colonusDir string) *BlockchainClientInitializedMessage {
	return &BlockchainClientInitializedMessage{
		event: Event{
			Id: id,
		},
		Time:        uint64(time.Now().Unix()),
		bcType:      bcType,
		bcInstance:  bcName,
		bcOrg:       bcOrg,
		serviceName: serviceName,
		servicePort: servicePort,
		colonusDir:  colonusDir,
	}
}

// Blockchain client Stopping message
type BlockchainClientStoppingMessage struct {
	event      Event
	Time       uint64
	bcType     string
	bcInstance string
	bcOrg      string
}

func (m *BlockchainClientStoppingMessage) Event() Event {
	return m.event
}

func (m BlockchainClientStoppingMessage) String() string {
	return fmt.Sprintf("Event: %v, Time: %v, Type: %v, Instance: %v, Org: %v", m.event, m.Time, m.bcType, m.bcInstance, m.bcOrg)
}

func (m BlockchainClientStoppingMessage) ShortString() string {
	return m.String()
}

func (m BlockchainClientStoppingMessage) BlockchainType() string {
	return m.bcType
}

func (m BlockchainClientStoppingMessage) BlockchainInstance() string {
	return m.bcInstance
}

func (m BlockchainClientStoppingMessage) BlockchainOrg() string {
	return m.bcOrg
}

func NewBlockchainClientStoppingMessage(id EventId, bcType string, bcName string, org string) *BlockchainClientStoppingMessage {
	return &BlockchainClientStoppingMessage{
		event: Event{
			Id: id,
		},
		Time:       uint64(time.Now().Unix()),
		bcType:     bcType,
		bcInstance: bcName,
		bcOrg:      org,
	}
}

// Report of blockchains that are needed
type ReportNeededBlockchainsMessage struct {
	event     Event
	Time      uint64
	bcType    string
	neededBCs map[string]map[string]bool
}

func (m *ReportNeededBlockchainsMessage) Event() Event {
	return m.event
}

func (m ReportNeededBlockchainsMessage) String() string {
	return fmt.Sprintf("Event: %v, Time: %v, Type: %v, Needed Blockchains: %v", m.event, m.Time, m.bcType, m.neededBCs)
}

func (m ReportNeededBlockchainsMessage) ShortString() string {
	return m.String()
}

func (m ReportNeededBlockchainsMessage) BlockchainType() string {
	return m.bcType
}

func (m ReportNeededBlockchainsMessage) NeededBlockchains() map[string]map[string]bool {
	return m.neededBCs
}

func NewReportNeededBlockchainsMessage(id EventId, bcType string, neededBCs map[string]map[string]bool) *ReportNeededBlockchainsMessage {
	return &ReportNeededBlockchainsMessage{
		event: Event{
			Id: id,
		},
		Time:      uint64(time.Now().Unix()),
		bcType:    bcType,
		neededBCs: neededBCs,
	}
}

// Blockchain event occurred
type EthBlockchainEventMessage struct {
	event    Event
	rawEvent string
	protocol string
	name     string
	org      string
	Time     uint64
}

func (m *EthBlockchainEventMessage) Event() Event {
	return m.event
}

func (m *EthBlockchainEventMessage) RawEvent() string {
	return m.rawEvent
}

func (m *EthBlockchainEventMessage) Name() string {
	return m.name
}

func (m *EthBlockchainEventMessage) Org() string {
	return m.org
}

func (m EthBlockchainEventMessage) String() string {
	return fmt.Sprintf("Event: %v, Name: %v, Org: %v, Protocol: %v, Raw Event: %v, Time: %v", m.event, m.name, m.org, m.protocol, m.rawEvent, m.Time)
}

func (m EthBlockchainEventMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, Name: %v, Org: %v, Protocol: %v, Time: %v", m.event, m.name, m.org, m.protocol, m.Time)
}

func NewEthBlockchainEventMessage(id EventId, ev string, name string, org string, protocol string) *EthBlockchainEventMessage {
	return &EthBlockchainEventMessage{
		event: Event{
			Id: id,
		},
		rawEvent: ev,
		protocol: protocol,
		name:     name,
		org:      org,
		Time:     uint64(time.Now().Unix()),
	}
}

// Exchange message received event occurred
type ExchangeDeviceMessage struct {
	event           Event
	exchangeMessage []byte
	protocolMessage string
	Time            uint64
}

func (m *ExchangeDeviceMessage) Event() Event {
	return m.event
}

func (m *ExchangeDeviceMessage) ExchangeMessage() []byte {
	return m.exchangeMessage
}

func (m *ExchangeDeviceMessage) ProtocolMessage() string {
	return m.protocolMessage
}

func (m *ExchangeDeviceMessage) ShortProtocolMessage() string {
	end := 200
	if len(m.protocolMessage) < end {
		end = len(m.protocolMessage)
	}
	return m.protocolMessage[:end]
}

func (m ExchangeDeviceMessage) String() string {
	return fmt.Sprintf("Event: %v, ProtocolMessage: %v, Time: %v, ExchangeMessage: %s", m.event, m.protocolMessage, m.Time, m.exchangeMessage)
}

func (m ExchangeDeviceMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, ProtocolMessage: %v, Time: %v", m.event, m.ShortProtocolMessage(), m.Time)
}

func NewExchangeDeviceMessage(id EventId, exMsg []byte, pMsg string) *ExchangeDeviceMessage {
	return &ExchangeDeviceMessage{
		event: Event{
			Id: id,
		},
		exchangeMessage: exMsg,
		protocolMessage: pMsg,
		Time:            uint64(time.Now().Unix()),
	}
}

// Make sure eth container is up and running
type NewBCContainerMessage struct {
	event         Event
	exchangeURL   string
	exchangeId    string
	exchangeToken string
	instance      string
	typeName      string
	org           string
	Time          uint64
}

func (m *NewBCContainerMessage) Event() Event {
	return m.event
}

func (m *NewBCContainerMessage) ExchangeURL() string {
	return m.exchangeURL
}

func (m *NewBCContainerMessage) ExchangeId() string {
	return m.exchangeId
}

func (m *NewBCContainerMessage) ExchangeToken() string {
	return m.exchangeToken
}

func (m *NewBCContainerMessage) Instance() string {
	return m.instance
}

func (m *NewBCContainerMessage) TypeName() string {
	return m.typeName
}

func (m *NewBCContainerMessage) Org() string {
	return m.org
}

func (m NewBCContainerMessage) String() string {
	return fmt.Sprintf("Event: %v, Type: %v, Instance: %v, Org: %v, Time: %v, ExchangeURL: %v, ExchangeId: %v, ExchangeToken: %v", m.event, m.typeName, m.instance, m.org, m.Time, m.exchangeURL, m.exchangeId, m.exchangeToken)
}

func (m NewBCContainerMessage) ShortString() string {
	return m.String()
}

func NewNewBCContainerMessage(id EventId, typeName string, instance string, org string, exchangeURL string, exchangeId string, exchangeToken string) *NewBCContainerMessage {
	return &NewBCContainerMessage{
		event: Event{
			Id: id,
		},
		exchangeURL:   exchangeURL,
		exchangeId:    exchangeId,
		exchangeToken: exchangeToken,
		instance:      instance,
		typeName:      typeName,
		org:           org,
		Time:          uint64(time.Now().Unix()),
	}
}

// Tell everyone that the device side of anax has synced up it's agreements wiht the exchange and blockchain
type DeviceAgreementsSyncedMessage struct {
	event     Event
	Completed bool
	Time      uint64
}

func (m *DeviceAgreementsSyncedMessage) Event() Event {
	return m.event
}

func (m *DeviceAgreementsSyncedMessage) IsCompleted() bool {
	return m.Completed
}

func (m DeviceAgreementsSyncedMessage) String() string {
	return fmt.Sprintf("Event: %v, Completed: %v, Time: %v", m.event, m.Completed, m.Time)
}

func (m DeviceAgreementsSyncedMessage) ShortString() string {
	return m.String()
}

func NewDeviceAgreementsSyncedMessage(id EventId, completed bool) *DeviceAgreementsSyncedMessage {
	return &DeviceAgreementsSyncedMessage{
		event: Event{
			Id: id,
		},
		Completed: completed,
		Time:      uint64(time.Now().Unix()),
	}
}

// Tell everyone that the device side of anax has synced up it's containers with the local DB
type DeviceContainersSyncedMessage struct {
	event     Event
	Completed bool
	Time      uint64
}

func (m *DeviceContainersSyncedMessage) Event() Event {
	return m.event
}

func (m *DeviceContainersSyncedMessage) IsCompleted() bool {
	return m.Completed
}

func (m DeviceContainersSyncedMessage) String() string {
	return fmt.Sprintf("Event: %v, Completed: %v, Time: %v", m.event, m.Completed, m.Time)
}

func (m DeviceContainersSyncedMessage) ShortString() string {
	return m.String()
}

func NewDeviceContainersSyncedMessage(id EventId, completed bool) *DeviceContainersSyncedMessage {
	return &DeviceContainersSyncedMessage{
		event: Event{
			Id: id,
		},
		Completed: completed,
		Time:      uint64(time.Now().Unix()),
	}
}

type MicroserviceMaintenanceMessage struct {
	event     Event
	MsInstKey string // the key to the microservice instance, it is used for network id and part of container name
}

func (m *MicroserviceMaintenanceMessage) Event() Event {
	return m.event
}

func (m MicroserviceMaintenanceMessage) String() string {
	return m.ShortString()
}

func (m MicroserviceMaintenanceMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, MsInstKey: %v", m.event, m.MsInstKey)
}

func NewMicroserviceMaintenanceMessage(id EventId, key string) *MicroserviceMaintenanceMessage {
	return &MicroserviceMaintenanceMessage{
		event: Event{
			Id: id,
		},
		MsInstKey: key,
	}
}

type MicroserviceCancellationMessage struct {
	event     Event
	MsInstKey string // the key to the microservice instance
}

func (m *MicroserviceCancellationMessage) Event() Event {
	return m.event
}

func (m MicroserviceCancellationMessage) String() string {
	return m.ShortString()
}

func (m MicroserviceCancellationMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, MsInstKey: %v", m.event, m.MsInstKey)
}

func NewMicroserviceCancellationMessage(id EventId, key string) *MicroserviceCancellationMessage {
	return &MicroserviceCancellationMessage{
		event: Event{
			Id: id,
		},
		MsInstKey: key,
	}
}

type MicroserviceContainersDestroyedMessage struct {
	event     Event
	MsInstKey string // the key to the microservice instance
}

func (m *MicroserviceContainersDestroyedMessage) Event() Event {
	return m.event
}

func (m MicroserviceContainersDestroyedMessage) String() string {
	return m.ShortString()
}

func (m MicroserviceContainersDestroyedMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, MsInstKey: %v", m.event, m.MsInstKey)
}

func NewMicroserviceContainersDestroyedMessage(id EventId, key string) *MicroserviceContainersDestroyedMessage {
	return &MicroserviceContainersDestroyedMessage{
		event: Event{
			Id: id,
		},
		MsInstKey: key,
	}
}

// Node lifecycle events
type NodeShutdownMessage struct {
	event      Event
	block      bool
	removeNode bool
}

func (n *NodeShutdownMessage) Event() Event {
	return n.event
}

func (n NodeShutdownMessage) String() string {
	return n.ShortString()
}

func (n NodeShutdownMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, Blocking: %v, RemoveNode: %v", n.event, n.block, n.removeNode)
}

func (n NodeShutdownMessage) Blocking() bool {
	return n.block
}

func (n NodeShutdownMessage) RemoveNode() bool {
	return n.removeNode
}

func NewNodeShutdownMessage(id EventId, blocking bool, removeNode bool) *NodeShutdownMessage {
	return &NodeShutdownMessage{
		event: Event{
			Id: id,
		},
		block:      blocking,
		removeNode: removeNode,
	}
}

type NodeShutdownCompleteMessage struct {
	event Event
	err   string
}

func (n *NodeShutdownCompleteMessage) Event() Event {
	return n.event
}

func (n NodeShutdownCompleteMessage) String() string {
	return n.ShortString()
}

func (n NodeShutdownCompleteMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, Error: %v", n.event, n.err)
}

func (n NodeShutdownCompleteMessage) Err() string {
	return n.err
}

func NewNodeShutdownCompleteMessage(id EventId, errorMsg string) *NodeShutdownCompleteMessage {
	return &NodeShutdownCompleteMessage{
		event: Event{
			Id: id,
		},
		err: errorMsg,
	}
}

// This is a special message that the message dispatcher knows about.
type WorkerStopMessage struct {
	event Event
	name  string
}

func (w *WorkerStopMessage) Event() Event {
	return w.event
}

func (w *WorkerStopMessage) String() string {
	return w.ShortString()
}

func (w *WorkerStopMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, Worker Name: %v", w.event, w.name)
}

func (w *WorkerStopMessage) Name() string {
	return w.name
}

func NewWorkerStopMessage(id EventId, name string) *WorkerStopMessage {
	return &WorkerStopMessage{
		event: Event{
			Id: id,
		},
		name: name,
	}
}

type AllBlockchainShutdownMessage struct {
	event Event
}

func (w *AllBlockchainShutdownMessage) Event() Event {
	return w.event
}

func (w *AllBlockchainShutdownMessage) String() string {
	return w.ShortString()
}

func (w *AllBlockchainShutdownMessage) ShortString() string {
	return fmt.Sprintf("Event: %v", w.event)
}

func NewAllBlockchainShutdownMessage(id EventId) *AllBlockchainShutdownMessage {
	return &AllBlockchainShutdownMessage{
		event: Event{
			Id: id,
		},
	}
}

type NodeHeartbeatStateChangeMessage struct {
	event   Event
	NodeOrg string
	NodeId  string
}

func (w *NodeHeartbeatStateChangeMessage) Event() Event {
	return w.event
}

func (w *NodeHeartbeatStateChangeMessage) String() string {
	return w.ShortString()
}

func (w *NodeHeartbeatStateChangeMessage) ShortString() string {
	return fmt.Sprintf("Event: %v, NodeOrg: %v, NodeId: %v", w.event, w.NodeOrg, w.NodeId)
}

func NewNodeHeartbeatStateChangeMessage(id EventId, node_org string, node_id string) *NodeHeartbeatStateChangeMessage {
	return &NodeHeartbeatStateChangeMessage{
		event: Event{
			Id: id,
		},
		NodeOrg: node_org,
		NodeId:  node_id,
	}
}
