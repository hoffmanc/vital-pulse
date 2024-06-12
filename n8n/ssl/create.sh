# Generate a private key
openssl genpkey -algorithm RSA -out server.key -aes256

# Generate a certificate signing request (CSR)
openssl req -new -key server.key -out server.csr

# Generate a self-signed certificate
openssl req -x509 -key server.key -in server.csr -out server.crt -days 365

# (Optional) Remove passphrase from private key
openssl rsa -in server.key -out server.key
