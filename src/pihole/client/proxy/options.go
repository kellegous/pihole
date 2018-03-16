package proxy

// Options ...
type Options struct {
	SSHAddr       string
	SSHUser       string
	SSHPrivateKey []byte
	APIAddr       string
	ClientHosts   []string
	ClientID      string
}
