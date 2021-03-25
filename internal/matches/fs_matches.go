package matches

import (
	"bytes"
	"github.com/chen-keinan/kube-knark/internal/routes"
)

//FSMatches Object
type FSMatches struct {
	fsCommandMap map[string]interface{}
	Cache        map[string]*routes.FS
}

//NewFSMatches create new file system matches instance
func NewFSMatches(fsCommandMap map[string]interface{}, cache map[string]*routes.FS) *FSMatches {
	return &FSMatches{fsCommandMap: fsCommandMap, Cache: cache}
}

//Match match fs command to specified commands pattern
func (mr *FSMatches) Match(cmd []string, sb *bytes.Buffer) bool {
	if len(cmd) == 0 {
		return false
	}
	return mr.recursiveMatch(mr.fsCommandMap, cmd, sb)
}

func (mr *FSMatches) recursiveMatch(currMap map[string]interface{}, cmd []string, sb *bytes.Buffer) bool {
	if len(currMap) == 0 {
		return true
	}
	if len(cmd) == 0 {
		return false
	}
	subMap, ok := currMap[cmd[0]]
	if !ok { // stay with current map
		return mr.recursiveMatch(currMap, cmd[1:], sb)
	} //continue to sub map
	sb.WriteString(cmd[0] + "_")
	return mr.recursiveMatch(subMap.(map[string]interface{}), cmd[1:], sb)
}
