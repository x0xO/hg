package hg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an HFloat for hashing.
type fhash struct{ hfloat HFloat }

// Hash returns a fhash struct wrapping the given HFloat.
func (hf HFloat) Hash() fhash { return fhash{hf} }

// MD5 computes the MD5 hash of the wrapped HFloat and returns the hash as an HString.
func (fh fhash) MD5() HString { return floatHasher(md5.New(), fh.hfloat) }

// SHA1 computes the SHA1 hash of the wrapped HFloat and returns the hash as an HString.
func (fh fhash) SHA1() HString { return floatHasher(sha1.New(), fh.hfloat) }

// SHA256 computes the SHA256 hash of the wrapped HFloat and returns the hash as an HString.
func (fh fhash) SHA256() HString { return floatHasher(sha256.New(), fh.hfloat) }

// SHA512 computes the SHA512 hash of the wrapped HFloat and returns the hash as an HString.
func (fh fhash) SHA512() HString { return floatHasher(sha512.New(), fh.hfloat) }

// floatHasher a helper function that computes the hash of the given HFloat using the specified
// hash.Hash algorithm and returns the hash as an HString.
func floatHasher(algorithm hash.Hash, fh HFloat) HString {
	algorithm.Write(fh.Bytes())
	return HString(hex.EncodeToString(algorithm.Sum(nil)))
}
