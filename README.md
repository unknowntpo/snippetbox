# snippetbox

Use following command to generate self-signed tls certificate.

$ go run /usr/local/opt/go/libexec/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

Then run the app:

$ go run ./cmd/web/
