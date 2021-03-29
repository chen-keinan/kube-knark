package common

const (
	//KubeKnarkHome env variable
	KubeKnarkHome = "KUBE_KNARK_HOME"
	//KprobeSourceFile kprobe program source file
	KprobeSourceFile = "kprobe.c"
	//BpfHeaderFile bpf header file
	BpfHeaderFile = "bpf.h"
	//BpfHelperHeaderFile bpf helper header file
	BpfHelperHeaderFile = "bpf_helpers.h"
	//KprobeCompiledFile kprobe program compiled file
	KprobeCompiledFile = "kprobe.elf"
	//GET Method
	GET = "GET"
	//DELETE Method
	DELETE = "DELETE"
	//PUT Method
	PUT = "PUT"
	//POST Method
	POST = "POST"
	//PATCH Method
	PATCH = "PATCH"
	//Workload spec
	Workload = "workload.yml"
	//Services spec
	Services = "services.yml"
	//ConfigAndStorage spec
	ConfigAndStorage = "config_and_storage.yml"
	//ConfigFilesPermission spec
	ConfigFilesPermission = "config_files_permission.yml"
	//Authentication spec
	Authentication = "authentication.yml"
	//Authorization spec
	Authorization = "authorization.yml"
	//Policy spec
	Policy = "policy.yml"
	//Extend spec
	Extend = "extend.yml"
	//Cluster spec
	Cluster = "cluster.yml"
)
