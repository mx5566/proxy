package proxy

import "testing"

func TestParseConfigFile(t *testing.T) {

	path := "./logs/tcpProxy.log"
	_ = parseConfigFile(path)
}
