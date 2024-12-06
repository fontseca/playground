<div align="center">
  <h1>fontseca.dev's playground</h1>
</div>

<figure>
  <img src="https://github.com/user-attachments/assets/87fea3d1-833e-4cdf-b091-28ec4b9693bc" alt="playground"/>
  <figcaption>
    <p><i>A screenshot of the initial version of the fontseca.dev's playground running with a collection exported from Postman.</i></p>
  </figcaption>
</figure>

The fontseca.dev's playground is a web-based HTTP client designed for testing APIs. It was originally built with the
intent to have a dedicated space to showcase my work as a back-end engineer, but it can be used with other third-party
APIs.

## Features

The Playground currently supports HTTP/1.1 servers and allows importing Postman collections in JSON format. It handles

The playground supports the following HTTP methods:

- GET,
- POST,
- PUT,
- PATCH,
- and DELETE.

### HTTP request features

- Query parameters
- HTTP headers
- HTTP body

The HTTP request body is passed as a raw string, with content differentiated by the `Content-Type` header.

### HTTP response features

- Response body
- HTTP headers
- Cookies

## Getting Started

To get started working with the Playground, clone this repository and simply run the application. You'll need to install
the [templ](https://github.com/a-h/templ) command to generate the website code.

By default, the Playground is intended to be used as a plug-in within a wrapper project. I use it as a Git submodule in
my [website's codebase](https://github.com/fontseca/.dev). If you want to run it separately, you'll have to create a
`go.mod` file.

1. Clone the repository:

    ```shell
    git clone git@github.com:fontseca/playground.git
    cd playground/
    ```

2. Generate the `website.templ` file:

    ```shell
    ~/go/bin/templ generate
    Processing path: /home/fontseca/playground
   (!) templ version check failed: failed to read go.mod file: open /go.mod: no such file or directory
   Generating production code: /home/fontseca/playground
   (✓) Generated code for "/home/fontseca/playground/website.templ" in 26.710246ms
   (✓) Generated code for 1 templates with 0 errors in 27.004517ms
    ```

3. Start a Go module called `playground`:

   ```shell
   go mod init playground
   go mod tidy
   ```

4. And finally, run the Playground:

   ```shell
   go run cmd/playground/main.go
   running fontseca.dev/playground server at 172.17.0.1:40273
   ```
