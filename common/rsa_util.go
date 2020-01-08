//author: https://github.com/zhaojunlike
package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

var (
	RsaPublicKeyError  = errors.New("public key error")
	RsaPrivateKeyError = errors.New("private key error")
)

// 加密
func RsaEncrypt(publicKey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, RsaPublicKeyError
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, RsaPrivateKeyError
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func RsaEncryptString(str string, pub string) (string, error) {
	if strings.Index(pub, "BEGIN") == -1 {
		pub = "-----BEGIN RSA PUBLIC KEY-----\n" + pub + "\n-----END RSA PUBLIC KEY-----"
	}
	bytes := []byte(str)
	pubBytes := []byte(pub)
	fmt.Println(pubBytes)
	res, err := RsaEncrypt(pubBytes, bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

//TODO 优化
type Rsa struct {
	publicKey  []byte
	privateKey []byte
}

func (rsa *Rsa) SetPublicKey(pub string) error {
	if strings.Index(pub, "BEGIN") == -1 {
		pub = "-----BEGIN RSA PUBLIC KEY-----\n" + pub + "\n-----END RSA PUBLIC KEY-----"
	}
	rsa.publicKey = []byte(pub)
	return nil
}

func (rsa *Rsa) Encrypt(str string) (string, error) {
	bytes := []byte(str)
	res, err := RsaEncrypt(rsa.publicKey, bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}
func (rsa *Rsa) Decrypt(str string) (string, error) {
	bytes := []byte(str)
	res, err := RsaDecrypt(rsa.privateKey, bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func NewRsa(pub string, pri string) *Rsa {
	rs := &Rsa{}
	_ = rs.SetPublicKey(pub)
	return rs
}
