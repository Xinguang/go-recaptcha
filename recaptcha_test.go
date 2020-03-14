package recaptcha

import (
	// "fmt"
	// "github.com/sirupsen/logrus"
	// "io/ioutil"
	// "net/http"
	// "net/url"
	"strings"
	"testing"
	// "time"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type ReCaptchaSuite struct{}

var _ = Suite(&ReCaptchaSuite{})

func (s *ReCaptchaSuite) TestNewReCAPTCHA(c *C) {
	captcha, err := NewReCAPTCHAWithSecert("my secret")
	c.Assert(err, IsNil)
	c.Check(captcha.Secret, Equals, "my secret")

	captcha, err = NewReCAPTCHAWithSecert("")
	c.Assert(err, NotNil)

	captcha, err = NewReCAPTCHA()
	c.Assert(err, NotNil)
}

func (s *ReCaptchaSuite) TestVerifyInvalidSolutionNoRemoteIp(c *C) {
	captcha := &ReCAPTCHA{
		Secret: "my secret",
	}

	err := captcha.Verify("mycode")
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "remote error codes: .*")
}
func (s *ReCaptchaSuite) TestVerifyInvalidSolutionWithRemoteIp(c *C) {
	captcha := &ReCAPTCHA{
		Secret: "my secret",
	}

	err := captcha.VerifyWithOptions("mycode", VerifyOption{RemoteIP: "localhost"})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "remote error codes: .*")
}

func (s *ReCaptchaSuite) TestVerifyUnmarshalFatal(c *C) {
	responseJson := `{
		"success": true,
	}`
	_, err := getReCAPTCHAResponse(responseJson)
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "invalid character '}' looking for beginning of object key string")
}

func (s *ReCaptchaSuite) TestVerifyUnmarshal(c *C) {
	responseJson := `{
		"success": true,
		"challenge_ts": "2020-03-15T03:41:29+00:00",
		"score": 0.5,
		"action": "homepage",
		"hostname": "localhost",
		"apk_package_name": "app name",
		"ErrorCodes": ["bad-request"]
	}`
	res, err := getReCAPTCHAResponse(responseJson)
	c.Assert(err, IsNil)

	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	captcha.showDebug(res)
}

func getReCAPTCHAResponse(responseJson string) (res reCAPTCHAResponse, err error) {
	body := strings.NewReader(responseJson)
	err = unmarshal(body, &res)
	return
}

func (s *ReCaptchaSuite) TestVerifyConfirmUnsuccess(c *C) {
	responseJson := `{
		"success": false
	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "invalid challenge solution")
}

func (s *ReCaptchaSuite) TestVerifyConfirmHostname(c *C) {
	responseJson := `{
		"success": true,
		"hostname": "localhost"

	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{Hostname: "localhost1"})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "invalid response hostname 'localhost', while expecting 'localhost1'")
}

func (s *ReCaptchaSuite) TestVerifyConfirmResponseTime(c *C) {
	responseJson := `{
		"success": true,
		"challenge_ts": "2020-03-14T18:21:29+09:00"

	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{ResponseTime: 1})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "time spent in resolving challenge .*")
}

func (s *ReCaptchaSuite) TestVerifyConfirmApkPackageName(c *C) {
	responseJson := `{
		"success": true,
		"apk_package_name": "app name"

	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{ApkPackageName: "app name1"})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "invalid response ApkPackageName 'app name', while expecting 'app name1'")
}

func (s *ReCaptchaSuite) TestVerifyConfirmScore(c *C) {
	responseJson := `{
		"success": true,
		"score": 0.4

	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{Threshold: 0.6})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "received score '0.400000', while expecting minimum '0.600000'")
	err = captcha.confirm(res, VerifyOption{})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "received score '0.400000', while expecting minimum '0.500000'")
}

func (s *ReCaptchaSuite) TestVerifyConfirmAction(c *C) {
	responseJson := `{
		"success": true,
		"action": "homepage"

	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{Action: "homepage1"})
	c.Assert(err, NotNil)
	c.Check(err, ErrorMatches, "invalid response action 'homepage', while expecting 'homepage1'")
}

func (s *ReCaptchaSuite) TestVerifyConfirmWithoutActionAndScore(c *C) {
	responseJson := `{
		"success": true

	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{})
	c.Assert(err, IsNil)
}

func (s *ReCaptchaSuite) TestVerifyConfirmWithoutAction(c *C) {
	responseJson := `{
		"success": true,
		"action": "homepage"
	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{})
	c.Assert(err, IsNil)
}
func (s *ReCaptchaSuite) TestVerifyConfirmScoreWithoutAction(c *C) {
	responseJson := `{
		"success": true,
		"score": 1.0
	}`
	res, err := getReCAPTCHAResponse(responseJson)
	captcha := ReCAPTCHA{
		Secret: "my secret",
	}
	err = captcha.confirm(res, VerifyOption{})
	c.Assert(err, IsNil)
}
