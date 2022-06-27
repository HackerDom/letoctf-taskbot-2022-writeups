package sizes

type ByteSize uint64

const (
	B ByteSize = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)
