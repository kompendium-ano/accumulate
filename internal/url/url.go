package url

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/AccumulateNetwork/accumulate/internal/encoding"
)

// ErrMissingHost means that a URL did not include a hostname.
var ErrMissingHost = errors.New("missing host")

// ErrWrongScheme means that a URL included a scheme other than the Accumulate
// scheme.
var ErrWrongScheme = errors.New("wrong scheme")

// URL is an Accumulate URL.
type URL struct {
	UserInfo  string
	Authority string
	Path      string
	Query     string
	Fragment  string
}

type Values = url.Values

// Parse parses the string as an Accumulate URL. The scheme may be omitted, in
// which case `acc://` will be added, but if present it must be `acc`. The
// hostname must be non-empty. RawPath, ForceQuery, and RawFragment are not
// preserved.
func Parse(s string) (*URL, error) {
	u, err := url.Parse(s)
	if err == nil && u.Scheme == "" {
		u, err = url.Parse("acc://" + s)
	}
	if err != nil {
		return nil, err
	}

	if u.Scheme != "acc" {
		return nil, ErrWrongScheme
	}

	if u.Host == "" || u.Host[0] == ':' {
		return nil, ErrMissingHost
	}

	v := new(URL)
	v.Authority = u.Host
	v.Path = u.Path
	v.Query = u.RawQuery
	v.Fragment = u.Fragment
	if u.User != nil {
		v.UserInfo = u.User.Username()
		if pw, _ := u.User.Password(); pw != "" {
			v.UserInfo += ":" + pw
		}
	}
	return v, nil
}

func splitColon(s string) (string, string) {
	t := strings.SplitN(s, ":", 2)
	if len(t) == 1 {
		return t[0], ""
	}
	return t[0], t[1]
}

// URL returns a net/url.URL.
func (u *URL) URL() *url.URL {
	v := new(url.URL)
	v.Scheme = "acc"
	v.Host = u.Authority
	v.Path = u.Path
	v.RawQuery = u.Query
	v.Fragment = u.Fragment
	if u.UserInfo != "" {
		user, pass := splitColon(u.UserInfo)
		if pass != "" {
			v.User = url.UserPassword(user, pass)
		} else {
			v.User = url.User(user)
		}
	}
	return v
}

// String reassembles the URL into a valid URL string. See net/url.URL.String().
func (u *URL) String() string {
	return u.URL().String()
}

// RawString concatenates all of the URL parts. Does not percent-encode
// anything. Primarily used for validation.
func (u *URL) RawString() string {
	s := "acc://"
	if u.UserInfo != "" {
		s += u.UserInfo + "@"
	}
	s += u.Authority + u.Path
	if u.Query != "" {
		s += "?" + u.Query
	}
	if u.Fragment != "" {
		s += "#" + u.Fragment
	}
	return s
}

// Hostname returns the hostname from the authority component.
func (u *URL) Hostname() string {
	s, _ := splitColon(u.Authority)
	return s
}

// Port returns the port from the authority component.
func (u *URL) Port() string {
	_, s := splitColon(u.Authority)
	return s
}

// Username returns the username from the user info component.
func (u *URL) Username() string {
	s, _ := splitColon(u.UserInfo)
	return s
}

// Password returns the password from the user info component.
func (u *URL) Password() string {
	_, s := splitColon(u.UserInfo)
	return s
}

// QueryValues parses Query and returns the corresponding values. It silently
// discards malformed value pairs. To check errors use net/url.ParseQuery.
func (u *URL) QueryValues() Values {
	v, _ := url.ParseQuery(u.Query)
	return v
}

func chain(s string) [32]byte {
	s = strings.ToLower(s)
	h := sha256.Sum256([]byte(s))
	return h
}

func ensurePath(s string) string {
	if s == "" || s[0] == '/' {
		return s
	}
	return "/" + s
}

// Identity returns a copy of the URL with an empty path.
func (u *URL) Identity() *URL {
	v := *u
	v.Path = ""
	return &v
}

// IdentityChain constructs a chain identifier from the lower case hostname. The
// port is not included.
//
//   Chain = Hash(LowerCase(u.Host()))
func (u *URL) IdentityChain() []byte {
	c := u.IdentityChain32()
	return c[:]
}

// IdentityChain32 returns IdentityChain as a [32]byte.
func (u *URL) IdentityChain32() [32]byte {
	return chain(u.Hostname())
}

// ResourceChain constructs a chain identifier from the lower case hostname and
// path. The port is not included. If the path does not begin with `/`, `/` is
// added between the hostname and the path.
//
//   Chain = Hash(LowerCase(Sprintf("%s/%s", u.Host(), u.Path)))
func (u *URL) ResourceChain() []byte {
	c := u.ResourceChain32()
	return c[:]
}

// ResourceChain32 returns ResourceChain as a [32]byte.
func (u *URL) ResourceChain32() [32]byte {
	return chain(u.Hostname() + ensurePath(u.Path))
}

// Routing returns the first 8 bytes of the identity chain as an integer.
//
//   Routing = uint64(u.IdentityChain()[:8])
func (u *URL) Routing() uint64 {
	return binary.BigEndian.Uint64(u.IdentityChain())
}

// Equal reports whether u and v, converted to strings and interpreted as UTF-8,
// are equal under Unicode case-folding.
func (u *URL) Equal(v *URL) bool {
	return strings.EqualFold(u.String(), v.String())
}

// JoinPath returns a copy of U with additional path elements.
func (u *URL) JoinPath(s ...string) *URL {
	v := *u
	v.Path = path.Join(append([]string{u.Path}, s...)...)
	return &v
}

// BinarySize returns the number of bytes the URL will marshal to.
func (u *URL) BinarySize() int {
	return encoding.StringBinarySize(u.String())
}

// MarshalBinary marshals the URL to binary.
func (u *URL) MarshalBinary() ([]byte, error) {
	return encoding.StringMarshalBinary(u.String()), nil
}

// UnmarshalBinary unmarshals the URL from binary.
func (u *URL) UnmarshalBinary(data []byte) error {
	s, err := encoding.StringUnmarshalBinary(data)
	if err != nil {
		return err
	}

	v, err := Parse(s)
	if err != nil {
		return err
	}

	*u = *v
	return nil
}

// MarshalJSON marshals the URL to JSON as a string.
func (u *URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON unmarshals the URL from JSON as a string.
func (u *URL) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	v, err := Parse(s)
	if err != nil {
		return err
	}

	*u = *v
	return nil
}
