//go:build windows && amd64

package license

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/s77rt/xor"
)

var (
	EncryptionKey string

	EncryptedServerName     string
	EncryptedServerLocalIP  string
	EncryptedServerPublicIP string
	EncryptedExpDate        string

	Signature string
)

type License struct {
	ServerName     string
	ServerLocalIP  net.IP
	ServerPublicIP net.IP
	ExpDate        time.Time
}

func Read() (*License, error) {
	encryptionKeyX := b64Encode(xor.XOR([]byte(EncryptionKey), []byte("RDPCloud")))

	readSignature := b64Encode(xor.XOR([]byte(EncryptedServerName+EncryptedServerLocalIP+EncryptedServerPublicIP+EncryptedExpDate+"SIGNATURE"), []byte(encryptionKeyX)))
	if readSignature != Signature {
		return nil, errors.New("Signature does not match")
	}

	var serverName string
	decodedServerName, err := b64Decode(EncryptedServerName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to decode server name (%v)", err))
	}
	decryptedServerName := xor.XOR(decodedServerName, []byte(encryptionKeyX))
	serverName = string(decryptedServerName)

	var serverLocalIP net.IP
	decodedServerLocalIP, err := b64Decode(EncryptedServerLocalIP)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to decode server local ip (%v)", err))
	}
	decryptedServerLocalIP := xor.XOR(decodedServerLocalIP, []byte(encryptionKeyX))
	serverLocalIP = net.ParseIP(string(decryptedServerLocalIP))

	var serverPublicIP net.IP
	decodedServerPublicIP, err := b64Decode(EncryptedServerPublicIP)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to decode server public ip (%v)", err))
	}
	decryptedServerPublicIP := xor.XOR(decodedServerPublicIP, []byte(encryptionKeyX))
	serverPublicIP = net.ParseIP(string(decryptedServerPublicIP))

	var expDate time.Time
	decodedExpDate, err := b64Decode(EncryptedExpDate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to decode exp date (%v)", err))
	}
	decryptedExpDate := xor.XOR(decodedExpDate, []byte(encryptionKeyX))
	if len(decryptedExpDate) > 0 {
		expDate, err = time.Parse("2006-01-02", string(decryptedExpDate))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to parse exp date (%v)", err))
		}
	}

	return &License{
		ServerName:     serverName,
		ServerLocalIP:  serverLocalIP,
		ServerPublicIP: serverPublicIP,
		ExpDate:        expDate,
	}, nil
}

func b64Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func b64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
