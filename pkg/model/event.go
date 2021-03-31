package model

import (
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
)

// K8sConfigFileChangeEvent fs event msg
type K8sConfigFileChangeEvent struct {
	Msg  *execevent.KprobeEvent
	Spec *specs.FS
}

// K8sAPICallEvent net event msg
type K8sAPICallEvent struct {
	Msg  *netevent.HTTPNetData
	Spec *specs.API
}
