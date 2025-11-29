# IMGA - URL Shortener

A simple URL shortener service built with Go.

## Building and Running

### Basic Build

```bash
go build
go run main.go
```

### Vanilla Docker

```bash
docker build -t imga .
docker run -p 8080:8080 imga
```

### With Docker Compose

```bash
docker-compose up
```

## Testing

### Run Tests

```bash
go test ./...
```

### Manual Testing

#### Create Short Code

```bash
curl -X POST -H "Content-Type: application/json" -d '{"URL":"http://example.com"}' localhost:8080/shorten
```

Should return something like the following:

```json
{
  "ShortCode": "CTv8Ax",
  "ShortURL": "http://localhost:8080/CTv8Ax",
  "OriginalURL": "http://example.com"
}
```

#### Test Created Short Code

You can run curl in the terminal:

```bash
curl http://localhost:8080/CTv8Ax
```

This should produce output like this in the terminal window running the program:

```
2025/11/29 20:28:12 Redirecting to original URL: http://example.com
```

Or enter the URL in the browser, and it should redirect you to example.com. You can open devtools and inspect the network tab to see the 302 redirect with the location header.
