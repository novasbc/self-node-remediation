/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	ConfigCRName                         = "self-node-remediation-config"
	defaultWatchdogPath                  = "/dev/watchdog"
	DefaultSafeToAssumeNodeRebootTimeout = 180
	defaultIsSoftwareRebootEnabled       = true
	defaultMinPeersForRemediation        = 1
)

// SelfNodeRemediationConfigSpec defines the desired state of SelfNodeRemediationConfig
type SelfNodeRemediationConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// WatchdogFilePath is the watchdog file path that should be available on each node, e.g. /dev/watchdog
	// +kubebuilder:default=/dev/watchdog
	WatchdogFilePath string `json:"watchdogFilePath,omitempty"`

	// SafeTimeToAssumeNodeRebootedSeconds is the time after which the healthy self node remediation
	// agents will assume the unhealthy node has been rebooted, and it is safe to recover affected workloads.
	// This is extremely important as starting replacement Pods while they are still running on the failed
	// node will likely lead to data corruption and violation of run-once semantics.
	// In an effort to prevent this, the operator ignores values lower than a minimum calculated from the
	// ApiCheckInterval, ApiServerTimeout, MaxApiErrorThreshold, PeerDialTimeout, and PeerRequestTimeout fields.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=180
	SafeTimeToAssumeNodeRebootedSeconds int `json:"safeTimeToAssumeNodeRebootedSeconds,omitempty"`

	// Valid time units are "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="5s"
	// +kubebuilder:validation:Pattern="^(0|([0-9]+(\\.[0-9]+)?(ms|s|m|h)))$"
	// +kubebuilder:validation:Type:=string
	PeerApiServerTimeout *metav1.Duration `json:"peerApiServerTimeout,omitempty"`

	// the frequency for api-server connectivity check
	// Valid time units are "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="15s"
	// +kubebuilder:validation:Pattern="^(0|([0-9]+(\\.[0-9]+)?(ms|s|m|h)))$"
	// +kubebuilder:validation:Type:=string
	// the frequency for api-server connectivity check
	ApiCheckInterval *metav1.Duration `json:"apiCheckInterval,omitempty"`

	// Valid time units are "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="15m"
	// +kubebuilder:validation:Pattern="^(0|([0-9]+(\\.[0-9]+)?(ms|s|m|h)))$"
	// +kubebuilder:validation:Type:=string
	PeerUpdateInterval *metav1.Duration `json:"peerUpdateInterval,omitempty"`

	// Valid time units are "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="5s"
	// +kubebuilder:validation:Pattern="^(0|([0-9]+(\\.[0-9]+)?(ms|s|m|h)))$"
	// +kubebuilder:validation:Type:=string
	// timeout for each api-connectivity check
	ApiServerTimeout *metav1.Duration `json:"apiServerTimeout,omitempty"`

	// Valid time units are "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="5s"
	// +kubebuilder:validation:Pattern="^(0|([0-9]+(\\.[0-9]+)?(ms|s|m|h)))$"
	// +kubebuilder:validation:Type:=string
	// timeout for establishing connection to peer
	PeerDialTimeout *metav1.Duration `json:"peerDialTimeout,omitempty"`

	// Valid time units are "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="5s"
	// +kubebuilder:validation:Pattern="^(0|([0-9]+(\\.[0-9]+)?(ms|s|m|h)))$"
	// +kubebuilder:validation:Type:=string
	// timeout for each peer request
	PeerRequestTimeout *metav1.Duration `json:"peerRequestTimeout,omitempty"`

	// +optional
	// +kubebuilder:default:=3
	// +kubebuilder:validation:Minimum=1
	// after this threshold, the node will start contacting its peers
	MaxApiErrorThreshold int `json:"maxApiErrorThreshold,omitempty"`

	// IsSoftwareRebootEnabled indicates whether self node remediation agent will do software reboot,
	// if the watchdog device can not be used or will use watchdog only,
	// without a fallback to software reboot
	// +kubebuilder:default=true
	IsSoftwareRebootEnabled bool `json:"isSoftwareRebootEnabled,omitempty"`

	// EndpointHealthCheckUrl is an url that self node remediation agents which run on control-plane node will try to access when they can't contact their peers.
	// This is a part of self diagnostics which will decide whether the node should be remediated or not.
	// It will be ignored when empty (which is the default).
	EndpointHealthCheckUrl string `json:"endpointHealthCheckUrl,omitempty"`

	// HostPort is used for internal communication between SNR agents.
	// +optional
	// +kubebuilder:default:=30001
	// +kubebuilder:validation:Minimum=1
	HostPort int `json:"hostPort,omitempty"`

	// CustomDsTolerations allows to add custom tolerations snr agents that are running on the ds in order to support remediation for different types of nodes.
	CustomDsTolerations []v1.Toleration `json:"customDsTolerations,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=0
	// Minimum number of peer workers/control nodes to attempt to contact before deciding if node is unhealthy or not
	MinPeersForRemediation int `json:"minPeersForRemediation,omitempty"`
}

// SelfNodeRemediationConfigStatus defines the observed state of SelfNodeRemediationConfig
type SelfNodeRemediationConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=snrc;snrconfig

// SelfNodeRemediationConfig is the Schema for the selfnoderemediationconfigs API in which a user can configure the self node remediation agents
// +operator-sdk:csv:customresourcedefinitions:resources={{"SelfNodeRemediationConfig","v1alpha1","selfnoderemediationconfigs"}}
type SelfNodeRemediationConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SelfNodeRemediationConfigSpec   `json:"spec,omitempty"`
	Status SelfNodeRemediationConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SelfNodeRemediationConfigList contains a list of SelfNodeRemediationConfig
type SelfNodeRemediationConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SelfNodeRemediationConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SelfNodeRemediationConfig{}, &SelfNodeRemediationConfigList{})
}

func NewDefaultSelfNodeRemediationConfig() SelfNodeRemediationConfig {
	return SelfNodeRemediationConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: ConfigCRName,
		},
		Spec: SelfNodeRemediationConfigSpec{
			WatchdogFilePath:                    defaultWatchdogPath,
			SafeTimeToAssumeNodeRebootedSeconds: DefaultSafeToAssumeNodeRebootTimeout,
			IsSoftwareRebootEnabled:             defaultIsSoftwareRebootEnabled,
			MinPeersForRemediation:              defaultMinPeersForRemediation,
		},
	}
}
