package crypto

import (
	"crypto"
	"net/http"
	"net/url"
	"time"

	"github.com/go-fed/httpsig"
)

// SignResponse will sign a response using the provided response body and public key.
func SignResponse(w http.ResponseWriter, body []byte, publicKey PublicKey) error {
	privateKey := GetPrivateKey()

	return signResponse(privateKey, *publicKey.ID, body, w)
}

func signResponse(privateKey crypto.PrivateKey, pubKeyID url.URL, body []byte, w http.ResponseWriter) error {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256

	// The "Date" and "Digest" headers must already be set on r, as well as r.URL.
	headersToSign := []string{}
	if body != nil {
		headersToSign = append(headersToSign, "digest")
	}

	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, 0)

	if err != nil {
		return err
	}

	// If r were a http.ResponseWriter, call SignResponse instead.
	return signer.SignResponse(privateKey, pubKeyID.String(), w, body)
}

// SignRequest will sign an ounbound request given the provided body.
func SignRequest(req *http.Request, body []byte, actorIRI *url.URL) error {
	publicKey := GetPublicKey(actorIRI)
	privateKey := GetPrivateKey()

	return signRequest(privateKey, publicKey.ID.String(), body, req)
}

func signRequest(privateKey crypto.PrivateKey, pubKeyID string, body []byte, r *http.Request) error {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256

	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	r.Header["Date"] = []string{date}
	r.Header["Host"] = []string{r.URL.Host}
	r.Header["Accept"] = []string{`application/ld+json; profile="https://www.w3.org/ns/activitystreams"`}

	// The "Date" and "Digest" headers must already be set on r, as well as r.URL.
	headersToSign := []string{httpsig.RequestTarget, "host", "date"}
	if body != nil {
		headersToSign = append(headersToSign, "digest")
	}

	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, 0)

	if err != nil {
		return err
	}

	// If r were a http.ResponseWriter, call SignResponse instead.
	return signer.SignRequest(privateKey, pubKeyID, r, body)
}
