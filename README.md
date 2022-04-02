`req` is a no-brainer library for making HTTP requests. It is heavily inspired by Python's [requests](https://docs.python-requests.org/en/latest/) library.

`req` is a thin wrapper around the standard library's `http` and `json` packages. It provides a simple, consistent interface for sending HTTP requests and parsing responses. Responses are type-safe, thanks to Go 1.18's addition of generic type parameters.

Features:

- Less boilerplate when writing Go scripts and small applications
- Type-safety through generic type parameters
- Easy configuration and extensions through [functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).

# Disclaimer

This library is very much for my own use, and I don't advise anyone to use it for anything other than education, or small personal scripts. Please, feel free to add requests, but abstain from stating the ovious - this library is library is unfinished, untested, or when it will ever be production-ready (it is not mean for production use).

# Usage

Check the example under `/examples/reqbin/main.go`:

```go
type ReqBinResponse struct {
	Success string `json:"success,omitempty"`
}

func main() {
	res, err := req.Get[ReqBinResponse]("https://reqbin.com/echo/get/json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", res)
}
```

# Motivation

Go's standard library provides all the bits and pieces one needs, in order to make HTTP requests. If you are building production Go, you have no reason to abstract those away. When it comes to quick-and-dirty apps and scripts, however, we can all take a tiny shortcut in the name of seeing results quickly. Going through the obtaining of results, checkign the error, converting the JSON, and checking the error again can be a little boiler-plaity. This is why I started work on `req` - to speed up the bootstrapping of scrapbook projects, and not to replace the standard library.
