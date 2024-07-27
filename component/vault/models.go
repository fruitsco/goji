package vault

type Secret struct {
	Name    string
	Version int
	Payload []byte
}
