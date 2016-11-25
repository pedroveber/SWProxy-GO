package main

import (
	"bytes"

	"testing"
)

var (
	versions = []int{1, 2}

	testCases = []struct {
		input []byte
	}{
		{[]byte("message")},
		{[]byte("a message which is longer than one block size")},
		{[]byte("a block size msg")},
	}
)

func TestSymmetric(t *testing.T) {
	for _, test := range testCases {
		for _, version := range versions {
			ciphertext, err := Encrypt(test.input, version)
			if err != nil {
				t.Error(err)
			}
			plaintext, err := Decrypt(ciphertext, version)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(plaintext, test.input) != 0 {
				t.Errorf("Plaintext: [%s], Result: [%s]\n", test.input, plaintext)
			}
		}
	}
}

func TestRequestAndResponse(t *testing.T) {
	funcs := []struct {
		en, de func([]byte, int) ([]byte, error)
	}{
		{EncryptResponse, DecryptResponse},
		{EncryptRequest, DecryptRequest},
	}
	for _, test := range testCases {
		for _, version := range versions {
			for _, funcPair := range funcs {
				ciphertext, err := funcPair.en(test.input, version)
				if err != nil {
					t.Error(err)
				}
				plaintext, err := funcPair.de(ciphertext, version)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(plaintext, test.input) != 0 {
					t.Errorf("Plaintext: [%s], Result: [%s]\n", test.input, plaintext)
				}
				reEncoded, err := funcPair.en(plaintext, version)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(reEncoded, ciphertext) != 0 {
					t.Error("Encoded messages do not match")
				}
			}
		}
	}
}
