package activityhelpers

import "strings"

const (
	cmdSymbol = "!"
)

/*IsOnCommand - check if message contains cmd
returns message, without command text
*/
func IsOnCommand(msgText string, cmdList []string) (string, bool) {
	// lets find delimiter position
	cmdSymbIdx := strings.IndexAny(msgText, "/"+cmdSymbol)
	if cmdSymbIdx == -1 {
		return "", false
	}

	found := false
	msg := msgText[cmdSymbIdx+len(cmdSymbol) : len(msgText)] //remove cmd string from query
	var strToFind string
	for _, cmd := range cmdList {
		if strings.HasPrefix(msg, cmd) {
			strToFind = strings.TrimSpace(strings.TrimPrefix(msg, cmd))
			found = true
		}
	}

	return strToFind, found
}
