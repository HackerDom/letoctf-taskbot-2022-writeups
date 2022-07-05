package api

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/blowfish"
)

func decrypt(key []byte, encrypted []byte) ([]byte, error) {
	newEnc := make([]byte, 0, len(encrypted))

	for i := 0; i < len(encrypted); i += 8 {
		penc, err := convertEndian(encrypted[i : i+8])
		if err != nil {
			return nil, fmt.Errorf("convert endian of key failed: %v", err)
		}
		newEnc = append(newEnc, penc...)
	}

	encrypted = newEnc

	cipher, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher failed: %v", err)
	}

	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i += cipher.BlockSize() {
		cipher.Decrypt(decrypted[i:i+cipher.BlockSize()], encrypted[i:i+cipher.BlockSize()])
	}

	newDec := make([]byte, 0, len(decrypted))
	for i := 0; i < len(decrypted); i += 8 {
		pdec, err := convertEndian(decrypted[i : i+8])
		if err != nil {
			return nil, fmt.Errorf("convert endian of key failed: %v", err)
		}
		newDec = append(newDec, pdec...)
	}

	return newDec, nil
}

func convertEndian(in []byte) ([]byte, error) {
	//Read byte array as uint32 (little-endian)
	var v1, v2 uint32
	buf := bytes.NewReader(in)
	if err := binary.Read(buf, binary.LittleEndian, &v1); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &v2); err != nil {
		return nil, err
	}

	//convert uint32 to byte array
	out := make([]byte, 8)
	binary.BigEndian.PutUint32(out, v1)
	binary.BigEndian.PutUint32(out[4:], v2)

	return out, nil
}
