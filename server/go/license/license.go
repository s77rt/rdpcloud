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

	EncryptedServerName string
	EncryptedServerIP   string
	EncryptedExpDate    string

	Signature string
)

type License struct {
	ServerName string
	ServerIP   net.IP
	ExpDate    time.Time
}

func Read() (*License, error) {
	encryptionKeyX := b64Encode(xor.XOR([]byte(EncryptionKey), []byte("RDPCloud")))

	readSignature := b64Encode(xor.XOR([]byte(EncryptedServerName+EncryptedServerIP+EncryptedExpDate+"SIGNATURE"), []byte(encryptionKeyX)))
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

	var serverIP net.IP
	decodedServerIP, err := b64Decode(EncryptedServerIP)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to decode server ip (%v)", err))
	}
	decryptedServerIP := xor.XOR(decodedServerIP, []byte(encryptionKeyX))
	serverIP = net.ParseIP(string(decryptedServerIP))

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
		ServerName: serverName,
		ServerIP:   serverIP,
		ExpDate:    expDate,
	}, nil
}

func b64Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func b64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
