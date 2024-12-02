package network

type GetStatusMessage struct{}
type StatusMessage struct {
	// the id of the server
	ID            string
	Version       uint32
	CurrentHeight uint32
}

type GetBlocksMessage struct {
	From uint32
	// If To is 0 the maximum blocks will be returned.
	To uint32
}
