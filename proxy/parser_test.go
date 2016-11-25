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
				cipherresp, err := funcPair.en(test.input, version)
				if err != nil {
					t.Error(err)
				}
				plainresp, err := funcPair.de(cipherresp, version)
				if err != nil {
					t.Error(err)
				}
				if bytes.Compare(plainresp, test.input) != 0 {
					t.Errorf("Plaintext: [%s], Result: [%s]\n", test.input, plainresp)
				}
			}
		}
	}
}
