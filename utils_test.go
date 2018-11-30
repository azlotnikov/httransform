package httransform

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type ExtractHostTestSuite struct {
	suite.Suite
}

func (suite *ExtractHostTestSuite) TestOK() {
	host, err := ExtractHost("hostname:80")

	suite.Equal(host, "hostname")
	suite.Nil(err)
}

func (suite *ExtractHostTestSuite) TestIPv4() {
	host, err := ExtractHost("127.0.0.1:80")

	suite.Equal(host, "127.0.0.1")
	suite.Nil(err)
}

func (suite *ExtractHostTestSuite) TestIPv6() {
	host, err := ExtractHost("[::1]:80")

	suite.Equal(host, "::1")
	suite.Nil(err)
}

func (suite *ExtractHostTestSuite) TestErr() {
	_, err := ExtractHost("hostname")

	suite.NotNil(err)
}

func (suite *ExtractHostTestSuite) TestEmpty() {
	_, err := ExtractHost("")

	suite.NotNil(err)
}

type MakeSimpleResponseTestSuite struct {
	suite.Suite
	resp *fasthttp.Response
}

func (suite *MakeSimpleResponseTestSuite) SetupTest() {
	suite.resp = fasthttp.AcquireResponse()
}

func (suite *MakeSimpleResponseTestSuite) TearDownTest() {
	fasthttp.ReleaseResponse(suite.resp)
}

func (suite *MakeSimpleResponseTestSuite) TestOverrideValues() {
	suite.resp.SetStatusCode(fasthttp.StatusMovedPermanently)
	suite.resp.SetBodyString("HELLO")
	suite.resp.SetConnectionClose()

	MakeSimpleResponse(suite.resp, "overriden", fasthttp.StatusOK)

	suite.Equal(suite.resp.StatusCode(), fasthttp.StatusOK)
	suite.False(suite.resp.ConnectionClose())
	suite.True(bytes.Equal(suite.resp.Body(), []byte("overriden")))
	suite.True(bytes.Equal(suite.resp.Header.ContentType(), []byte("text/plain")))
}

type ExtractAuthenticationTestSuite struct {
	suite.Suite
}

func (suite *ExtractAuthenticationTestSuite) TestEmpty() {
	_, _, err := ExtractAuthentication([]byte{})

	suite.NotNil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestGarbage() {
	_, _, err := ExtractAuthentication([]byte("adfjsfsjfhaskfsjfsjh"))

	suite.NotNil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestDigest() {
	_, _, err := ExtractAuthentication([]byte("Digest dXNlcjpwYXNz"))

	suite.NotNil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestLowerCasedType() {
	_, _, err := ExtractAuthentication([]byte("basic dXNlcjpwYXNz"))

	suite.NotNil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestIncorrectPayload() {
	_, _, err := ExtractAuthentication([]byte("Basic XXXXXXXXX"))

	suite.NotNil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestNoPassword() {
	_, _, err := ExtractAuthentication([]byte("Basic dXNlcg=="))

	suite.NotNil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestEmptyPassword() {
	user, password, err := ExtractAuthentication([]byte("Basic dXNlcjo="))

	suite.True(bytes.Equal(user, []byte("user")))
	suite.Len(password, 0)
	suite.Nil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestEmptyUser() {
	user, password, err := ExtractAuthentication([]byte("Basic OnBhc3M="))

	suite.True(bytes.Equal(password, []byte("pass")))
	suite.Len(user, 0)
	suite.Nil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestAllEmpty() {
	user, password, err := ExtractAuthentication([]byte("Basic Og=="))

	suite.Len(user, 0)
	suite.Len(password, 0)
	suite.Nil(err)
}

func (suite *ExtractAuthenticationTestSuite) TestUserPass() {
	user, password, err := ExtractAuthentication([]byte("Basic dXNlcjpwYXNz"))

	suite.True(bytes.Equal(user, []byte("user")))
	suite.True(bytes.Equal(password, []byte("pass")))
	suite.Nil(err)
}

type MakeProxyAuthorizationHeaderValueTestSuite struct {
	suite.Suite
}

func (suite *MakeProxyAuthorizationHeaderValueTestSuite) TestEmpty() {
	result := MakeProxyAuthorizationHeaderValue(&url.URL{})

	suite.Len(result, 0)
}

func (suite *MakeProxyAuthorizationHeaderValueTestSuite) TestUserOnly() {
	result := MakeProxyAuthorizationHeaderValue(&url.URL{
		User: url.User("username"),
	})

	suite.True(bytes.Equal(result, []byte("Basic dXNlcm5hbWU6")))
}

func (suite *MakeProxyAuthorizationHeaderValueTestSuite) TestUserPass() {
	result := MakeProxyAuthorizationHeaderValue(&url.URL{
		User: url.UserPassword("username", "password"),
	})

	suite.True(bytes.Equal(result, []byte("Basic dXNlcm5hbWU6cGFzc3dvcmQ=")))
}

func TestExtractHost(t *testing.T) {
	suite.Run(t, &ExtractHostTestSuite{})
}

func TestMakeSimpleResponse(t *testing.T) {
	suite.Run(t, &MakeSimpleResponseTestSuite{})
}

func TestExtractAuthentication(t *testing.T) {
	suite.Run(t, &ExtractAuthenticationTestSuite{})
}

func TestMakeProxuAuthorizationHeaderValue(t *testing.T) {
	suite.Run(t, &MakeProxyAuthorizationHeaderValueTestSuite{})
}