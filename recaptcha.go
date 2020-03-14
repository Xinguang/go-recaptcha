package recaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

const (
	// recaptcha API to ensure the token is valid
	reCAPTCHALink = "https://www.google.com/recaptcha/api/siteverify"
	// DefaultTreshold Default minimin score when using V3 api
	DefaultTreshold = 0.5
)

var (
	logger = logrus.WithFields(logrus.Fields{"package": "ReCAPTCHA"})
)

type reCAPTCHAResponse struct {
	Success        bool      `json:"success"`                    // whether this request was a valid reCAPTCHA token for your site
	ChallengeTS    time.Time `json:"challenge_ts"`               // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Score          *float64  `json:"score,omitempty"`            // the score for this request (0.0 - 1.0)
	Action         *string   `json:"action,omitempty"`           // the action name for this request (important to verify)
	Hostname       string    `json:"hostname,omitempty"`         // the hostname of the site where the reCAPTCHA was solved
	ApkPackageName *string   `json:"apk_package_name,omitempty"` // the package name of the app where the reCAPTCHA was solved
	ErrorCodes     []string  `json:"error-codes,omitempty"`      // optional
}

// VerifyOption verification options expected for the challenge
type VerifyOption struct {
	Threshold      float64 // ignored in v2 recaptcha
	Action         string  // ignored in v2 recaptcha
	Hostname       string
	ApkPackageName string
	ResponseTime   float64
	RemoteIP       string
}

// ReCAPTCHA recpatcha holder struct, make adding mocking code simpler
type ReCAPTCHA struct {
	Secret string
}

// NewReCAPTCHA new ReCAPTCHA instance
func NewReCAPTCHA() (*ReCAPTCHA, error) {
	return NewReCAPTCHAWithSecert(os.Getenv("ReCAPTCHA_SECRET"))
}

// NewReCAPTCHAWithSecert new ReCAPTCHA instance
// get your secret from https://www.google.com/recaptcha/admin
func NewReCAPTCHAWithSecert(secret string) (*ReCAPTCHA, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("recaptcha secret cannot be blank")
	}
	return &ReCAPTCHA{
		Secret: secret,
	}, nil
}

// Verify returns `nil` if no error and the client solved the challenge correctly
func (r *ReCAPTCHA) Verify(token string) error {
	return r.VerifyWithOptions(token, VerifyOption{})
}

// VerifyWithOptions returns `nil` if no error and the client solved the challenge correctly and all options are natching
// `Threshold` and `Action` are ignored when using V2 version
func (r *ReCAPTCHA) VerifyWithOptions(token string, options VerifyOption) error {
	res, err := r.fetch(token, options.RemoteIP)
	if err != nil {
		logger.Error("confirm:", err)
		return err
	}
	return r.confirm(res, options)
}

func (r *ReCAPTCHA) fetch(token, remoteip string) (res reCAPTCHAResponse, err error) {
	var req http.Request
	req.ParseForm()
	req.Form.Add("secret", r.Secret)
	req.Form.Add("response", token)
	if len(remoteip) > 0 {
		req.Form.Add("remoteip", remoteip)
	}
	body := strings.NewReader(req.Form.Encode())
	logger.Info("fetch:", body)
	resp, err := http.Post(reCAPTCHALink, "application/x-www-form-urlencoded", body)
	defer resp.Body.Close()
	if err != nil {
		logger.Error("fetch: ", err)
		return
	}
	err = unmarshal(resp.Body, &res)
	// debug info
	r.showDebug(res)
	return
}

func (r *ReCAPTCHA) confirm(res reCAPTCHAResponse, options VerifyOption) (err error) {
	if res.ErrorCodes != nil {
		err = fmt.Errorf("remote error codes: %v", res.ErrorCodes)
		return
	}

	if !res.Success {
		err = fmt.Errorf("invalid challenge solution")
		return
	}
	// the hostname of the site where the reCAPTCHA was solved
	if len(options.Hostname) > 0 && options.Hostname != res.Hostname {
		err = fmt.Errorf("invalid response hostname '%s', while expecting '%s'", res.Hostname, options.Hostname)
		return
	}

	if options.ResponseTime != 0 {
		duration := time.Since(res.ChallengeTS)
		if options.ResponseTime < duration.Seconds() {
			err = fmt.Errorf("time spent in resolving challenge '%fs', while expecting maximum '%fs'", duration.Seconds(), options.ResponseTime)
			return
		}
	}
	// the package name of the app where the reCAPTCHA was solved
	if res.ApkPackageName != nil && len(options.ApkPackageName) > 0 && options.ApkPackageName != *res.ApkPackageName {
		err = fmt.Errorf("invalid response ApkPackageName '%s', while expecting '%s'", *res.ApkPackageName, options.ApkPackageName)
		return
	}
	// V3 api
	err = r.confirmV3(res, options)
	return
}

// V3 api
func (r *ReCAPTCHA) confirmV3(res reCAPTCHAResponse, options VerifyOption) (err error) {
	// ignored in v2 recaptcha
	if res.Score == nil && res.Action == nil {
		return
	}
	// the action name for this request
	if res.Action != nil && len(options.Action) > 0 && options.Action != *res.Action {
		err = fmt.Errorf("invalid response action '%s', while expecting '%s'", *res.Action, options.Action)
		return
	}
	// the score for this request (0.0 - 1.0)
	if res.Score == nil {
		return
	}
	threshold := DefaultTreshold
	if options.Threshold != 0 {
		threshold = options.Threshold
	}
	if threshold >= *res.Score {
		err = fmt.Errorf("received score '%f', while expecting minimum '%f'", *res.Score, threshold)
		return
	}
	return

}
func (r *ReCAPTCHA) showDebug(res reCAPTCHAResponse) {
	logger.Debug("res.Success:", res.Success)
	logger.Debug("res.ChallengeTS:", res.ChallengeTS)
	logger.Debug("res.Hostname:", res.Hostname)
	logger.Debug("res.ErrorCodes:", res.ErrorCodes)
	if res.Score != nil {
		logger.Debug("res.Score:", *res.Score)
	}
	if res.Action != nil {
		logger.Debug("res.Action:", *res.Action)
	}
	if res.ApkPackageName != nil {
		logger.Debug("res.ApkPackageName:", *res.ApkPackageName)
	}
}

func unmarshal(body io.Reader, v interface{}) error {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		logger.Errorf("ioutil.ReadAll: %s", err)
		return err
	}
	bodyBytes = bytes.TrimPrefix(bodyBytes, []byte("\xef\xbb\xbf"))

	var test interface{}
	err = json.Unmarshal(bodyBytes, &test)
	logger.Debugf("test: %s", test)

	err = json.Unmarshal(bodyBytes, &v)
	if err != nil {
		logger.Errorf("unmarshal: %s", err)
		return err
	}
	return nil
}
