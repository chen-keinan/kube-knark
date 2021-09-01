package main

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model"
)

//OnK8sFileConfigChangeHook this plugin method accept k8s file config change event
//(event include file change command chmod / chown) and it send the event command args to an external events service
//nolint
func OnK8sFileConfigChangeHook(k8sConfigFileChangeEvent model.K8sConfigFileChangeEvent) error {
	fmt.Println("this is OnK8sFileConfigChangeHook plugin")
	return nil
}
