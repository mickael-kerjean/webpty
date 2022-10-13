package ssl

import (
	"os"
)

var keyPEMPath string = "./cert/key.pem"
var certPEMPath string = "./cert/cert.pem"

func init() {
	os.MkdirAll("./cert", os.ModePerm)
}

func Clear() {
	clearPrivateKey()
	clearCert()
}
