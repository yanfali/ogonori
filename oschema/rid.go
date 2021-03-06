package oschema

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ClusterIDInvalid  = -1
	ClusterPosInvalid = -1
)

type ORID struct {
	ClusterID  int16
	ClusterPos int64
}

//
// Returns an ORID with the default "invalid" settings.
// Invalid settings indicate that the Document has not yet been saved
// to the DB (which assigns it a valid RID) or it indicates that
// it is not a true Document with a Class
// (e.g., it is a result of a Property query)
//
func NewORID() ORID {
	return ORID{ClusterID: ClusterIDInvalid, ClusterPos: ClusterPosInvalid}
}

func (r ORID) String() string {
	return fmt.Sprintf("#%d:%d", r.ClusterID, r.ClusterPos)
}

//
// NewORIDFromString converts a string of form #N:M or N:M
// to an ORID struct. Make sure to get the string format correctly,
// as this function panics if any error occurs.
//
func NewORIDFromString(s string) ORID {
	noPrefix := s
	if strings.HasPrefix(s, "#") {
		noPrefix = s[1:]
	}
	toks := strings.Split(noPrefix, ":")
	if len(toks) != 2 {
		panic(fmt.Errorf("Invalid RID string to NewORIDFromString: %s", s))
	}
	id, err := strconv.ParseInt(toks[0], 10, 16)
	if err != nil {
		panic(fmt.Errorf("Invalid RID string to NewORIDFromString: %s", s))
	}
	pos, err := strconv.ParseInt(toks[1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("Invalid RID string to NewORIDFromString: %s", s))
	}
	return ORID{ClusterID: int16(id), ClusterPos: pos}
}
