package credential

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	lestrratjws "github.com/lestrrat-go/jwx/v2/jws"
	lestrratjwt "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"golang.org/x/oauth2/jws"
)

type CredentialPrivyProperty struct {
	Host     string
	Username string
	Password string
	Client   *http.Client
}

type CredentialPrivy struct {
	host      string
	username  string
	password  string
	requester *requester
}

func NewCredentialPrivy(prop CredentialPrivyProperty) *CredentialPrivy {
	if prop.Client == nil {
		prop.Client = http.DefaultClient
	}

	r := &requester{
		hc: prop.Client,
	}

	return &CredentialPrivy{
		host:      prop.Host,
		username:  prop.Username,
		password:  prop.Password,
		requester: r,
	}
}

func (c *CredentialPrivy) GenerateJwtTokenWithNode(ctx context.Context) (JWTToken, error) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:3000", nil)
	req.Header.Set("Content-Type", "application/json")

	jwtToken := JWTToken{}
	err := c.requester.Do(ctx, req, &jwtToken)
	if err != nil {
		return JWTToken{}, err
	}

	return jwtToken, nil
}

func (c *CredentialPrivy) getPrivateKey() (*rsa.PrivateKey, error) {
	secret := `-----BEGIN PRIVATE KEY-----
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
-----END PRIVATE KEY-----`
	// secret := "MIIG/gIBADANBgkqhkiG9w0BAQEFAASCBugwggbkAgEAAoIBgQC/3oTe8dtl/7UIwz6rHkOmSB0F6pxvzNdTcRA3Hr3gPBIv6NX+0JYE4y7dYv00FavDPgGKDcLYNGjv49u1Cvlq0nwGWQHRVFgLsI+ToHxoNrsnIj3MZWGpiVPRqjv7x7Jgxme9+MLHdrwuM9ohZvODMadRpgEWFHr74p2CvN2WsfgKnU2AvJGBjigiIqrbQfODzTr5rNY3t+XbygOAcStu5fo9ZVk4SUf0A+J1Pfr0OIWCEt3ZAyufttmX+FrDOtn/x76xWwGZqwPnZfKk0uQCYmQzm3govSqIE4imZgoJWwmWMtPSZQag7WQDZQKa/Lhdsucar+/r4/dv6zWheDnoKCSmm4cH6fr/wWT5skNvUpA9nZf/2DwVVNtx0MC8ZwIFa2XbCLCOIKXgmQQ0TMcwj1zrAEOolKoRbABLm092Uz55vaFY64irsMvnu8ZV43FD95ZTgACKa+RSBxCoddDqcJimnOfTp2PRCprU2DuR8rGyUHNS9unwHq+AZn4cxYkCAwEAAQKCAYB/UIk8YqIh0YZv9RZ9d4yOJuXTSjVZ3kO3Y5vN54E47MNotQhimEgjoBg14PyA9ixOVrOwxqbWzhgfrkPYoxqTrm2LzL3vCBeQUi1DWUeT41VWh1sYaOFgHPbYdixtSrpritvF6+5basc9pTyF04mcMXOEMzIfU6GzyFCaMvaaiyc669xEo1ut6wVoeTaEXQPYHnOWvwA6YrWMnUSrkuyuHr5oxPpp5f5vxbAb+e+u5F56zEgXOwpUZYGHXM3lJAPhHWfFWq20ZaE3wO7Qr7dVgJbWxb8BnFdWe2azzamDSJw3wFDm4YVRsQ1F6gXfIEqLl4GdmvVJxOwFAENKV7jgUwNNjm+rA4Xv6L7dnJsjYyD5+lPDyt7YcBX6o2GSr+xww8Aaw5s4tYx2n1C8KvYFyP9Ggei5Bq9609V4bRF8kQ3JGZy473u1Rx4MEix+Tr/CtBkR1Lof9CqUS7K7F90nt8HwPegI3fPwJfB1Z0phjXmSmjmEQoMgV/cdRywDIG0CgcEA/HOiqolSnIoNfeqsL0bCUNkoJl0Z3mmpcJAY/MC0KhZUkmig1keKLTDz+Sun3fb+vVthufjIPrEb65eoxv5GBm9NISfMWvft5GXq8MyW0UccurXr3aJWPY6hsc4jdo83KBGdK7cCbZI/1T4GadjFpLoOvz9jR4SZbSHb1msG1nSGBlW+07rKASnFc+6f9i5xiJXpxFoH7erlrMUYFSxqkTNrFhJ8jUywI9E1begWWOlJrqGp25uZfFN0DAwk5783AoHBAMKQ5b1lY6qhAkuWX/IIY/LGM8vZc7FU59dlT+/6uMtmT92KQLGuBiM5zbMLX7A0r6kf7LhfZp02970e5xbAXRQwvCfxB13KTb9ZJxljpdLyruigm7bv7eNpt7rDT+nsWrX6gzAcz4556tgToRgbocxz49QXfORvuotGUfMWN6Yo6nbb6L7MLrGTMn2SeoLkH11zy4Y1Gz1OdAgmETyHeY3JA/TjZwym3Fqdf3MB7MNlO+jvK57DS/YmqF725RaBPwKBwQDCnvMcvdkuTU+xbyVnHQnZAlDtooC68oJqAO0Cjh1XBPgWFwtHpsdjp4Wu5nT3rd8dZEGm+aCGlEuOCNFY99ZMR6oSkXf9+X0ww1GoalEq7cO8PVuk9e2+byNAzEaStD41ab7dYK1Cg6kqYDLZjwBvyfBsiBWloRgXBi9Q0hnnYtWgNKJ04F+zHdzXT3OKequUPN6HPVE3nguKcjfqut89KYK49W+ID0MLmdGy9WWlDdR8CK3GW+Kr8jpYv1QERWkCgcEAowsXamI3Zfos2Ti3SDRdxUjwiMe9moHjEm400Y5SIyimjqjXU83YGNbMmFhUpG1SMrCtB0fyzpYMfrARyNAEb/HzCqmBmcf45PuJt0343NA/YHOOaXuf5u1laJ1ZL1bAITU/kCbki6mA8fdpDLHDIXiQF+Bi6W7zbNjNvZ4Fnjk1WcsovBKQooAOVIpWHa+a1Q4/JEUGgZZnx5hW41lHtFgZ97JEXabKiyjmb3LSfF8uGCGsdQuFcU2t1H6jNPuzAoHAOkM7o62HJDc4jv6RAa57zcDLPGi5T+6yVYrQta5Ve46ZfBFHX4VjkFOljBBQE72OHKUYOOx1n5jbYpliTc9fArz0lWibBolUPBF2pKBb70n83D7qGgKUKZ19lYu/rVMdFHBWuLu+n8eqK8B+zBqe5Lm37HSnfRiFVbKp+LakZ7YkK9NW1qmqkcr8bf+Ue4wYiQgKRGHuwIRfC7JBpyfaKpbrd/qFaJDpAr2XRlXj5H2YmVvHY1rxkeFxkUapk7GP"
	// bts, err := base64.RawStdEncoding.DecodeString(secret)

	privPem, _ := pem.Decode([]byte(secret))
	if privPem == nil {
		return nil, errors.New("err: Something went wrong")
	}

	keys, _ := x509.ParsePKCS8PrivateKey(privPem.Bytes)

	prvKey, ok := keys.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("err: Not a private key")
	}

	return prvKey, nil
}

func (c *CredentialPrivy) GenerateJwtTokenGojose(ctx context.Context) (JWTToken, error) {
	jwtToken := JWTToken{}
	_ = lestrratjwt.NewBuilder().
		Issuer("f8b4fb43614c756a168b1d1367607569ae547e8ffb6001b0ec7526641627c112").
		Audience([]string{"https://tstdrv2245019.suitetalk.api.netsuite.com/services/rest/auth/oauth2/v1/token"})

	lestrratjws.WithHeaders(nil)

	return jwtToken, nil
}

func (c *CredentialPrivy) GenerateJwtToken(ctx context.Context) (JWTToken, error) {
	jwtToken := JWTToken{}

	hdr := &jws.Header{
		Algorithm: "PS256",
		Typ:       "JWT",
		KeyID:     "T17p9_H2k1WjJB3h3vQk0qFJlbFERxUZRx2w71TDw7k",
	}

	clmSet := &jws.ClaimSet{
		Iss:   "f8b4fb43614c756a168b1d1367607569ae547e8ffb6001b0ec7526641627c112",
		Scope: "restlets rest_webservices",
		Aud:   "https://tstdrv2245019.suitetalk.api.netsuite.com/services/rest/auth/oauth2/v1/token",
		Iat:   time.Now().Unix(),
		Exp:   time.Now().Add(1 * time.Hour).Unix(),
		Typ:   "JWT",
	}

	secret := `-----BEGIN PRIVATE KEY-----
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
-----END PRIVATE KEY-----`
	// secret := "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"
	// bts, err := base64.RawStdEncoding.DecodeString(secret)

	privPem, _ := pem.Decode([]byte(secret))
	if privPem == nil {
		return jwtToken, errors.New("err: Something went wrong")
	}

	keys, _ := x509.ParsePKCS8PrivateKey(privPem.Bytes)

	prvKey, ok := keys.(*rsa.PrivateKey)
	if !ok {
		return jwtToken, errors.New("err: Not a private key")
	}

	s, err := jws.Encode(hdr, clmSet, prvKey)
	if err != nil {
		return jwtToken, errors.New("err: Something went wrong jws encode")
	}

	jwtTokenx, _ := c.GenerateJwtTokenWithNode(ctx)

	fmt.Println("xxxxxx jwtTokenx", jwtTokenx.SignedJWT)
	fmt.Println("xxxxxx s", s)
	fmt.Println("xxxxxx s==jwtTokenx.SignedJWT", s == jwtTokenx.SignedJWT)

	jwtToken.SignedJWT = s
	return jwtToken, nil
}

func (c *CredentialPrivy) CreateCustomer(ctx context.Context, param CustomerParam) (CustomerResponse, error) {
	// get jwt
	isNode := true
	jwtToken := JWTToken{}
	var err error

	if isNode {
		jwtToken, err = c.GenerateJwtTokenWithNode(ctx)
	} else {
		jwtToken, err = c.GenerateJwtToken(ctx)
	}

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "JWTToken{}",
			}).
			Error(err)

		return CustomerResponse{}, err
	}

	// get access token
	accessTokenURL := c.host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	logrus.
		WithFields(logrus.Fields{
			"at":                    "CredentialPrivy.CreateCustomer",
			"src":                   "EnvelopeCustomer{}.beforeDo",
			"grant_type":            "client_credentials",
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"client_assertion":      jwtToken.SignedJWT,
		}).
		Info(accessTokenURL)

	req, _ := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.username, c.password)

	credential := CredentialResponse{}
	err = c.requester.Do(ctx, req, &credential)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "CredentialResponse{}",
			}).
			Error(err)

		return CustomerResponse{}, err
	}

	// post customer
	postCustomerURL := c.host + EndpointPostCustomer

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateCustomer",
			"src":  "EnvelopeCustomer{}.beforeDo",
			"host": postCustomerURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, postCustomerURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", credential.TokenType+" "+credential.AccessToken)

	q := req.URL.Query()
	q.Add("script", "9")
	q.Add("deploy", "1")

	req.URL.RawQuery = q.Encode()

	custResp := EnvelopeCustomer{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "EnvelopeCustomer{}",
			}).
			Error(err)

		return CustomerResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return CustomerResponse{}, rapperror.ErrNotFound(
			"",
			"Customer is not found",
			"",
			nil,
		)
	}

	return custResp.SuccessTransaction[0], nil
}

func (c *CredentialPrivy) CreateCustomerUsage(ctx context.Context, param CustomerUsageParam) (CustomerUsageResponse, error) {
	// get jwt
	isNode := true
	jwtToken := JWTToken{}
	var err error

	if isNode {
		jwtToken, err = c.GenerateJwtTokenWithNode(ctx)
	} else {
		jwtToken, err = c.GenerateJwtToken(ctx)
	}

	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "JWTToken{}",
			}).
			Error(err)

		return CustomerUsageResponse{}, err
	}

	// get access token
	accessTokenURL := c.host + EndpointGetAccessToken
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", jwtToken.SignedJWT)

	logrus.
		WithFields(logrus.Fields{
			"at":                    "CredentialPrivy.CreateCustomer",
			"src":                   "EnvelopeCustomer{}.beforeDo",
			"grant_type":            "client_credentials",
			"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
			"client_assertion":      jwtToken.SignedJWT,
		}).
		Info(accessTokenURL)

	req, _ := http.NewRequest(http.MethodPost, accessTokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.username, c.password)

	credential := CredentialResponse{}
	err = c.requester.Do(ctx, req, &credential)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "CredentialResponse{}",
			}).
			Error(err)

		return CustomerUsageResponse{}, err
	}

	// post customer
	postCustomerURL := c.host + EndpointPostCustomer

	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(param)

	logrus.
		WithFields(logrus.Fields{
			"at":   "CredentialPrivy.CreateCustomer",
			"src":  "EnvelopeCustomer{}.beforeDo",
			"host": postCustomerURL,
		}).
		Info(body.String())

	req, _ = http.NewRequest(http.MethodPost, postCustomerURL, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", credential.TokenType+" "+credential.AccessToken)

	q := req.URL.Query()
	q.Add("script", "10")
	q.Add("deploy", "1")

	req.URL.RawQuery = q.Encode()

	custResp := EnvelopeCustomerUsage{}
	err = c.requester.Do(ctx, req, &custResp)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "CredentialPrivy.CreateCustomer",
				"src": "EnvelopeCustomer{}",
			}).
			Error(err)

		return CustomerUsageResponse{}, err
	}

	if len(custResp.SuccessTransaction) == 0 {
		return CustomerUsageResponse{}, rapperror.ErrNotFound(
			"",
			"Customer is not found",
			"",
			nil,
		)
	}

	return custResp.SuccessTransaction[0], nil
}
