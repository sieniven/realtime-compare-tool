package compare

import "fmt"

var (
	ErrCtxCancelled = fmt.Errorf("context cancelled - stopping")
)

const (
	DefaultHeightSyncRange = 5
	DefaultMismatchCount   = 5
	DefaultChannelSize     = 1000
	DefaultCacheSize       = 1000
)
