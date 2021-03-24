package matches

//FSMatches Object
type FSMatches struct {
	fsCommandMap map[string]interface{}
}

//NewFSMatches create new file system matches instance
func NewFSMatches(fsCommandMap map[string]interface{}) *FSMatches {
	return &FSMatches{fsCommandMap: fsCommandMap}
}

//Match match fs command to specified commands pattern
func (mr *FSMatches) Match(cmd []string) bool {
	if len(cmd) == 0 {
		return false
	}
	return mr.recursiveMatch(mr.fsCommandMap, cmd)
}

func (mr *FSMatches) recursiveMatch(currMap map[string]interface{}, cmd []string) bool {
	if len(currMap) == 0 {
		return true
	}
	if len(cmd) == 0 {
		return false
	}
	subMap, ok := currMap[cmd[0]]
	if !ok { // stay with current map
		return mr.recursiveMatch(currMap, cmd[1:])
	} //continue to sub map
	return mr.recursiveMatch(subMap.(map[string]interface{}), cmd[1:])
}
