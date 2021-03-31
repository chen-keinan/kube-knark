package kplugin

import "plugin"

//K8sAPICallHook hold the plugin symbol for K8s API Call Hook
type K8sAPICallHook struct {
	Plugins []plugin.Symbol
}

//K8sFileConfigChangeHook hold the plugin symbol for K8s File Config Change Hook
type K8sFileConfigChangeHook struct {
	Plugins []plugin.Symbol
}
