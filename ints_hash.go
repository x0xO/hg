package hg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an HInt for hashing.
type ihash struct{ hint HInt }

// Hash returns a ihash struct wrapping the given HInt.
func (hi HInt) Hash() ihash { return ihash{hi} }

// MD5 computes the MD5 hash of the wrapped HInt and returns the hash as an HString.
func (ih ihash) MD5() HString { return intHasher(md5.New(), ih.hint) }

// SHA1 computes the SHA1 hash of the wrapped HInt and returns the hash as an HString.
func (ih ihash) SHA1() HString { return intHasher(sha1.New(), ih.hint) }

// SHA256 computes the SHA256 hash of the wrapped HInt and returns the hash as an HString.
func (ih ihash) SHA256() HString { return intHasher(sha256.New(), ih.hint) }

// SHA512 computes the SHA512 hash of the wrapped HInt and returns the hash as an HString.
func (ih ihash) SHA512() HString { return intHasher(sha512.New(), ih.hint) }

// intHasher a helper function that computes the hash of the given HInt using the specified hash.Hash algorithm and returns the hash as an HString.
func intHasher(algorithm hash.Hash, ih HInt) HString {
	algorithm.Write(ih.Bytes())
	return HString(hex.EncodeToString(algorithm.Sum(nil)))
}
