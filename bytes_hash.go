package hg

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// A struct that wraps an HBytes for hashing.
type bhash struct{ hbytes HBytes }

// Hash returns a bhash struct wrapping the given HBytes.
func (hbs HBytes) Hash() bhash { return bhash{hbs} }

// MD5 computes the MD5 hash of the wrapped HBytes and returns the hash as an HBytes.
func (bh bhash) MD5() HBytes { return bytesHasher(md5.New(), bh.hbytes) }

// SHA1 computes the SHA1 hash of the wrapped HBytes and returns the hash as an HBytes.
func (bh bhash) SHA1() HBytes { return bytesHasher(sha1.New(), bh.hbytes) }

// SHA256 computes the SHA256 hash of the wrapped HBytes and returns the hash as an HBytes.
func (bh bhash) SHA256() HBytes { return bytesHasher(sha256.New(), bh.hbytes) }

// SHA512 computes the SHA512 hash of the wrapped HBytes and returns the hash as an HBytes.
func (bh bhash) SHA512() HBytes { return bytesHasher(sha512.New(), bh.hbytes) }

// bytesHasher a helper function that computes the hash of the given HBytes using the specified
// hash.Hash algorithm and returns the hash as an HBytes.
func bytesHasher(algorithm hash.Hash, hbs HBytes) HBytes {
	algorithm.Write(hbs)
	return HBytes(hex.EncodeToString(algorithm.Sum(nil)))
}
