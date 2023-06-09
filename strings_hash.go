package hg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an HString for hashing.
type shash struct{ hstr HString }

// Hash returns a shash struct wrapping the given HString.
func (hs HString) Hash() shash { return shash{hs} }

// MD5 computes the MD5 hash of the wrapped HString and returns the hash as an HString.
func (sh shash) MD5() HString { return stringHasher(md5.New(), sh.hstr) }

// SHA1 computes the SHA1 hash of the wrapped HString and returns the hash as an HString.
func (sh shash) SHA1() HString { return stringHasher(sha1.New(), sh.hstr) }

// SHA256 computes the SHA256 hash of the wrapped HString and returns the hash as an HString.
func (sh shash) SHA256() HString { return stringHasher(sha256.New(), sh.hstr) }

// SHA512 computes the SHA512 hash of the wrapped HString and returns the hash as an HString.
func (sh shash) SHA512() HString { return stringHasher(sha512.New(), sh.hstr) }

// stringHasher a helper function that computes the hash of the given HString using the specified
// hash.Hash algorithm and returns the hash as an HString.
func stringHasher(algorithm hash.Hash, hstr HString) HString {
	algorithm.Write(hstr.Bytes())
	return HString(hex.EncodeToString(algorithm.Sum(nil)))
}
