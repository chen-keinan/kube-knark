package main

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model"
)

//OnK8sAPICallHook this plugin method accept k8s api call event
//(event include http request/response, and match spec data) and it send the event URI to an external event service
//nolint
func OnK8sAPICallHook(k8sAPICallEvent model.K8sAPICallEvent) error {
	fmt.Println("this is OnK8sAPICallHook plugin")
	return nil
}
