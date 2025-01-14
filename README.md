# lmdrouterv2

[![](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue&style=flat-square)](https://godoc.org/github.com/shodgson/lmdrouterv2) [![](https://img.shields.io/github/license/shodgson/lmdrouterv2?style=flat-square)](LICENSE) 

**Go HTTP router library for AWS API Gateway-invoked Lambda Functions**

**Full credit to [aquasecurity](https://github.com/aquasecurity/lmdrouter) for developing the router library for the AWS REST API (V1)**

## Table of Contents

* [Overview](#overview)
    * [Use Case](#use-case)
    * [Features](#features)
* [Status](#status)
* [Installation](#installation)
* [Usage](#usage)
* [License](#license)

## Overview

`lmdrouterv2` is a simple-to-use library for writing AWS Lambda functions in Go
that listen to events of type API Gateway Proxy Request. It allows creating a
lambda function that can match requests based on their URI, just like an HTTP
server would.

The library provides an interface not unlike the standard `net/http.Mux` type
or community libraries such as [httprouter](https://github.com/julienschmidt/httprouter)
and [chi](https://github.com/go-chi/chi).

### Use Case

When building large cloud-native applications, there's a certain balance to
strike when it comes to deployment of APIs. On one side of the scale, each API
endpoint has its own lambda function. This provides the greatest flexibility,
but is extremely difficult to maintain. On the other side of the scale, there
can be one lambda function for the entire API. This provides the least flexibility,
but is the easiest to maintain. Both are probably not a good idea.

With `lmdrouter`, one can create small lambda functions for different aspects of
the API. For example, if your application model contains multiple domains (e.g.
articles, authors, topics, etc...), you can create one lambda function for each
domain, and deploy these independently (e.g. everything below "/api/articles" is
one lambda function, everything below "/api/authors" is another function). This
is also useful for applications where different teams are in charge of different
parts of the API.

### Features

* Supports all HTTP methods.
* Supports middleware at a global and per-resource level.
* Supports path parameters with a simple ":<name>" format (e.g. "/posts/:id").
* Provides ability to automatically "unmarshal" an API Gateway request to an
  arbitrary Go struct, with data coming either from path and query string
  parameters, or from the request body (only JSON requests are currently
  supported).
* Provides ability to automatically "marshal" responses of any type to an API
  Gateway response (only JSON responses are currently generated).

## Status

This is a very early, alpha release. API is subject to change.

## Installation

```shell
go get github.com/shodgson/lmdrouterv2
```

## Usage

`lmdrouterv2` is meant to be used inside Go Lambda functions.

```go
package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/shodgson/lmdrouterv2"
)

var router *lmdrouterv2.Router

func init() {
    router = lmdrouterv2.NewRouter("/api", loggerMiddleware, authMiddleware)
    router.Route("GET", "/", listSomethings)
    router.Route("POST", "/", postSomething, someOtherMiddleware)
    router.Route("GET", "/:id", getSomething)
    router.Route("PUT", "/:id", updateSomething)
    router.Route("DELETE", "/:id", deleteSomething)
}

func main() {
    lambda.Start(router.Handler)
}

// the rest of the code is a redacted example, it will probably reside in a
// separate package inside your project

type listSomethingsInput struct {
    ID                string   `lambda:"path.id"`                // a path parameter declared as :id
    ShowSomething     bool     `lambda:"query.show_something"`   // a query parameter named "show_something"
    AcceptedLanguages []string `lambda:"header.Accept-Language"` // a multi-value header parameter
}

type postSomethingInput struct {
    Title   string    `json:"title"`
    Date    time.Time `json:"date"`
}

func listSomethings(ctx context.Context, req events.APIGatewayV2HTTPRequest) (
    res events.APIGatewayProxyResponse,
    err error,
) {
    // parse input from request and path parameters
    var input listSomethingsInput
    err = lmdrouterv2.UnmarshalRequest(req, false, &input)
    if err != nil {
        return lmdrouterv2.HandleError(err)
    }

    // call some business logic that generates an output struct
    // ...

    return lmdrouterv2.MarshalResponse(http.StatusOK, nil, output)
}

func postSomethings(ctx context.Context, req events.APIGatewayV2HTTPRequest) (
    res events.APIGatewayProxyResponse,
    err error,
) {
    // parse input from request body
    var input postSomethingsInput
    err = lmdrouterv2.UnmarshalRequest(req, true, &input)
    if err != nil {
        return lmdrouterv2.HandleError(err)
    }

    // call some business logic that generates an output struct
    // ...

    return lmdrouterv2.MarshalResponse(http.StatusCreated, nil, output)
}

func loggerMiddleware(next lmdrouterv2.Handler) lmdrouterv2.Handler {
    return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (
        res events.APIGatewayProxyResponse,
        err error,
    ) {
        // [LEVEL] [METHOD PATH] [CODE] EXTRA
        format := "[%s] [%s %s] [%d] %s"
    	level := "INF"
    	var code int
    	var extra string

    	res, err = next(ctx, req)
    	if err != nil {
    	    level = "ERR"
    	    code = http.StatusInternalServerError
    	    extra = " " + err.Error()
    	} else {
    	    code = res.StatusCode
    	    if code >= 400 {
    	        level = "ERR"
    	    }
        }

        log.Printf(format, level, req.HTTPMethod, req.Path, code, extra)

        return res, err
    }
}
```

## License

This library is distributed under the terms of the [Apache License 2.0](LICENSE).
