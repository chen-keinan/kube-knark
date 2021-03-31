package model

import (
	"github.com/chen-keinan/kube-knark/pkg/model/execevent"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"github.com/chen-keinan/kube-knark/pkg/model/specs"
)

// FilesystemEvt fs event msg
type FilesystemEvt struct {
	Msg  *execevent.KprobeEvent
	Spec *specs.FS
}

// NetEvt net event msg
type NetEvt struct {
	Msg  *netevent.HTTPNetData
	Spec *specs.API
}
