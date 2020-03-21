<pre align="center" type="ascii-art">
                                            _       _         
  __ _  ___        _ __ ___  ___ __ _ _ __ | |_ ___| |__   __ _ 
 / _` |/ _ \ _____| '__/ _ \/ __/ _` | '_ \| __/ __| '_ \ / _` |
| (_| | (_) |_____| | |  __/ (_| (_| | |_) | || (__| | | | (_| |
 \__, |\___/      |_|  \___|\___\__,_| .__/ \__\___|_| |_|\__,_|
 |___/                               |_|                       
</pre>
# go-reCaptcha

This package handles [reCaptcha](https://www.google.com/recaptcha) (API versions [2](https://developers.google.com/recaptcha/docs/display) and [3](https://developers.google.com/recaptcha/docs/v3)) form submissions in [Go](http://golang.org/).

[![Build Status](https://travis-ci.org/Xinguang/go-recaptcha.svg)](https://travis-ci.org/Xinguang/go-recaptcha)
[![Go Report Card](https://goreportcard.com/badge/github.com/xinguang/go-recaptcha)](https://goreportcard.com/report/github.com/xinguang/go-recaptcha)
[![GoDoc](https://godoc.org/github.com/xinguang/go-recaptcha?status.svg)](https://godoc.org/github.com/xinguang/go-recaptcha)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/xinguang/go-recaptcha?tab=doc)

[![Sourcegraph](https://sourcegraph.com/github.com/Xinguang/go-recaptcha/-/badge.svg)](https://sourcegraph.com/github.com/Xinguang/go-recaptcha?badge)
[![Release](https://img.shields.io/github/release/Xinguang/go-recaptcha.svg?style=flat-square)](https://github.com/Xinguang/go-recaptcha/releases)

## Getting setup

* [Usage](#Usage)
* [Get Google reCaptcha Site Key And Secret Key](docs/reCaptcha.md)
* [Integrate Google reCAPTCHA in your website](#website)


### Usage <a name="Usage" />

* Install the package in your environment:
  
``` sh
go get github.com/xinguang/go-recaptcha
```

* To use it within your own code, import github.com/xinguang/go-recaptcha and call:

```go
package main

import (
	"github.com/xinguang/go-recaptcha"
)
func main() {
  // get your secret from https://www.google.com/recaptcha/admin
  // Using environment variables
  // export ReCAPTCHA_SECRET="reCaptcha Secret Key"
  recaptcha, err := New()

  // OR 
  const Secret = "reCaptcha Secret Key"
  recaptcha, err := NewWithSecert(Secret)
  // .....
}
```

* Now everytime you need to verify a API client with no special options request use

```go
  // the recaptchaResponse corresponds to 
  // the value of g-recaptcha-response sent by the reCaptcha server.
  recaptchaResponse := "g-recaptcha-response"
  err := recaptcha.Verify(recaptchaResponse)
  if err != nil {
      // do something with err
  }
  // proceed

```

For specific options use the `VerifyWithOptions` method  
Available options:

```go
  Threshold      float64 // ignored in v2 reCaptcha
  Action         string  // ignored in v2 reCaptcha
  Hostname       string
  ApkPackageName string
  ResponseTime   float64
  RemoteIP       string
```

Note that as reCaptcha v3 use score for challenge validation, if no threshold option is set the default value is 0.5


```go
err := recaptcha.VerifyWithOptions(recaptchaResponse, VerifyOption{RemoteIP: "127.0.0.1"})
if err != nil {
    // do something with err
}
// proceed
```

```go
err := recaptcha.VerifyWithOptions(recaptchaResponse, VerifyOption{Action: "homepage", Threshold: 0.8})
if err != nil {
    // do something with err
}
// proceed
```

Both recaptcha.Verify and recaptcha.VerifyWithOptions return a error or nil if successful

Use the error to check for issues with the secret, connection with the server, options mismatches and incorrect solution.


### Integrate Google reCAPTCHA in your website <a name="website" />

* [reCAPTCHA v3](https://developers.google.com/recaptcha/docs/v3)

To integrate it into your website you need to put it in the client side as well as in Server side. In client HTML page you need to integrate this line before the tag.

```html
<script src="https://www.google.com/recaptcha/api.js?render=put your site key here"></script>
```

Google reCAPTCHA v3 is invisible. You won’t see a captcha form of any sort on your web page. You need to capture the google captcha response in your JavaScript code. Here is a small snippet.


```html
<script>
  grecaptcha.ready(function() {
      grecaptcha.execute('put your site key here', {action: 'homepage'}).then(function(token) {
        // pass the token to the backend script for verification
      });
  });
</script>
```


* [reCAPTCHA v2](https://developers.google.com/recaptcha/docs/display)

To integrate it into your website you need to put it in client side as well as in Server side. In client HTML page you need to integrate this line before <HEAD> tag.

```html
<script src='https://www.google.com/recaptcha/api.js' async defer></script>
```
And to show the widget into your form you need to put this below contact form, comment form etc.

```html
<div class="g-recaptcha" data-sitekey="== Your site Key =="></div>
```

When the form get submit to Server, this script will send ‘g-recaptcha-response’ as a POST data. You need to verify it in order to see whether user has checked the Captcha or not.

Let’s look at the [sample code](./example/README.md) to understand it better.
**Replace `your_site_key` with your reCAPTCHA site key**