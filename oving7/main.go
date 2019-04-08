package main

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Serving", req.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello world"))
}

func main() {

	message := []byte("important message!")

	// Hash'en av sertifikatet etter kyptering med VeriSign's private nøkkel
	// Hash'en av sertifikatet må være lik signaturfeltet etter dekryptering med VeriSign's offentlige nøkkel
	signature, _ := hex.DecodeString("9c3e8d77333fcee3885747250fd48c8a6a5a8e62c24f8ef5f578c752469880409f69fa94a70dae0f71acc7a3988cc81e66881cbc75d5096dedfeeb3d17fb88fd27abe5d32f3b705a11045a91b5b5986f34948009e9b35e8026f986ae871e986392ae37e0458223d62b05fbb50935f63fa920590454d7851d35bf7b3d4cf0752c4683666bcb0398843d141113f32442f8d38f7910a43102da331a6e56fd2a3b3dbe49abf15b4e93c5a81341ed9f87e6bd972536e185e2cde096105db51de519f980901585b2c312b8a097853434bf144a3f14182f2d1b971169280b15061b781a21b8954c626aa4d9417275c1b1812eb0b9770b8320db2f1093f6e775105d39d5")

	rsaKey := []byte(`
-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEA639u2haGdEoEQ5wf7lfTHEvDW2FuLBNmZgailV3N9L2JCI9NKtk1
QOlEW2t6jweRfzjNf7Qs9XZkk6v6hveW2AZAYuhbNxQFT1FOk+Ez2RFVLLNZfIc+
sXD0VURkORY7m+CFHfT+pf6hlLrvZONEWdJ1ZmxDtMOH6hTESCOooxdJ8m2+WsA5
GuzOvaagZD/P4Gf9uoVjk/+G4jsB3YyaGAu+hs/Xx/ti9xPwFtCiUloJlUxhsDz9
my67QMmPype4vv1w2Hhaj3UabCQi5qj4JgSctNayRy73Wk0iXtos1s2S38CUsUuS
L7oZWDeIi2pZS0NT7e8cZllAHgSuX8MW+wIDAQAB
-----END RSA PUBLIC KEY-----`)

	// TODO:
	// dekrypter signaturen
	// hash meldingen
	// 	sjekk at de to er like

	publicKey := BytesToPublicKey(rsaKey)

	fmt.Printf("%+v\n", publicKey)

	hashed := sha1.Sum(message)
	m := new(big.Int)
	m.SetBytes(signature)

	encrypted := decrypt(new(big.Int), publicKey, m)

	fmt.Println(hex.EncodeToString(hashed[:]))
	fmt.Println(hex.EncodeToString(encrypted.Bytes()[235:]))

	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	// block, _ := pem.Decode(pub)
	// var pk rsa.PublicKey
	// asn1.Unmarshal(block.Bytes, &pk)
	// return &pk

	block, _ := pem.Decode(pub)

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	return key
}

func decrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}
