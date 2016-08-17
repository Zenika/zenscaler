dockerDNS=localhost
dockerPublicIP=10.10.10.20

clean: clean-csr
	rm *.crt *.key *.cnf
	rm *.pem

all: ca.crt server.crt client.crt links

%.key:
	openssl genrsa -out $@ 4096

ca.crt: ca.key
	openssl req -new -x509 -days 365 -key ca.key -sha256 -out ca.crt \
		-subj "/C=FR/O=zenika/CN=${dockerDNS}/emailAdress=security@zenika.com"

server.csr: server.key
	openssl req -subj "/CN=${dockerDNS}" -sha256 -new -key server.key -out server.csr

extfile.cnf:
	echo subjectAltName = IP:${dockerPublicIP},IP:127.0.0.1,DNS:${dockerDNS} > extfile.cnf

# sign server key with CA
server.crt: server.csr extfile.cnf
	openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.crt -CAkey ca.key \
	  -CAcreateserial -out server.crt -extfile extfile.cnf

client.csr: client.key
	openssl req -subj '/CN=client' -new -key client.key -out client.csr

extfile.client.cnf:
	echo extendedKeyUsage = clientAuth > extfile.client.cnf

client.crt: client.csr extfile.client.cnf
	openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.crt -CAkey ca.key \
		  -CAcreateserial -out client.crt -extfile extfile.client.cnf

clean-csr:
	rm -v client.csr server.csr

links:
	ln -s ca.crt ca.pem
	ln -s client.crt cert.pem
	ln -s client.key key.pem
