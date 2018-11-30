package httransform

import (
	"bytes"
	"encoding/base64"
	"net"
	"net/url"

	"github.com/juju/errors"
	"github.com/valyala/fasthttp"
)

func ExtractHost(rawurl string) (string, error) {
	host, _, err := net.SplitHostPort(rawurl)
	if err != nil {
		return "", errors.Annotate(err, "Cannot split host/port")
	}

	return host, nil
}

func MakeSimpleResponse(resp *fasthttp.Response, msg string, statusCode int) {
	resp.Reset()
	resp.SetBodyString(msg)
	resp.SetStatusCode(statusCode)
	resp.Header.SetContentType("text/plain")
}

func ExtractAuthentication(text []byte) ([]byte, []byte, error) {
	pos := bytes.IndexByte(text, ' ')
	if pos < 0 {
		return nil, nil, errors.New("Malformed Proxy-Authorization header")
	}
	if !bytes.Equal(text[:pos], []byte("Basic")) {
		return nil, nil, errors.New("Incorrect authorization prefix")
	}

	for pos < len(text) && (text[pos] == ' ' || text[pos] == '\t') {
		pos++
	}

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(text[pos:])))
	n, err := base64.StdEncoding.Decode(decoded, text[pos:])
	decoded = decoded[:n]
	if err != nil {
		return nil, nil, errors.Annotate(err, "Incorrectly encoded authorization payload")
	}

	pos = bytes.IndexByte(decoded, ':')
	if pos < 0 {
		return nil, nil, errors.New("Cannot find a user/password delimiter in decoded authorization string")
	}

	return decoded[:pos], decoded[pos+1:], nil
}

func MakeProxyAuthorizationHeaderValue(proxyURL *url.URL) []byte {
	username := proxyURL.User.Username()
	password, ok := proxyURL.User.Password()
	if ok || username != "" {
		line := username + ":" + password
		return []byte("Basic " + base64.StdEncoding.EncodeToString([]byte(line)))
	}

	return nil
}
