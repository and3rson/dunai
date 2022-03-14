all:
	false

# certs:
# 	openssl genrsa -out dun.ai.key 1024
# 	openssl req -new -x509 -key dun.ai.key -out dun.ai.crt -days 50000
# 	openssl pkcs12 -export -out dun.ai.pfx -inkey dun.ai.key -in dun.ai.crt

certs:
	go run ../cert.go dun.ai 1234567h
