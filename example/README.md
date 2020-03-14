# Using this Example

```sh
go get github.com/xinguang/go-recaptcha
cd $GOPATH/src/github.com/xinguang/go-recaptcha/example
ReCAPTCHA_SECRET="reCaptcha Secret Key" go run main.go
```
You can access the page from http://localhost:8002/ in your browser.

For more information on client side setup and other configuration options, check the [official documentation](https://developers.google.com/recaptcha/intro).

 <sup>&#42;</sup> make sure ['localhost' is added to the list of domains allowed](https://developers.google.com/recaptcha/docs/domain_validation) for the site registered at reCaptcha.


## reCAPTCHA v3

- [example](./html/v3/index.html)

## reCAPTCHA v2

- [Automatically render the reCAPTCHA widget](https://developers.google.com/recaptcha/docs/display#automatically_render_the_recaptcha_widget)
    
    - [example](./html/v2/automatically.html)

- [Explicitly render the reCAPTCHA widget](https://developers.google.com/recaptcha/docs/display#explicitly_render_the_recaptcha_widget)

    - [example1](./html/v2/explicitly.html)
    - [example2](./html/v2/explicitly.verifyCallback.html)
