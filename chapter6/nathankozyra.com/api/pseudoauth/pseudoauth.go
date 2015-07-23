package pseudoauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Token struct {
	Valid       bool
	Created     int64
	Expires     int64
	ForUser     int
	AccessToken string
}

var nonces map[string]Token

func init() {
	nonces = make(map[string]Token)
}

func ValidateSignature(consumer_key string, consumer_secret string, timestamp string, nonce string, signature string, for_user int) (Token, error) {
	var hashKey []byte
	t := Token{}
	t.Created = time.Now().UTC().Unix()
	t.Expires = t.Created + 600
	t.ForUser = for_user

	qualifiedMessage := []string{consumer_key, consumer_secret, timestamp, nonce}
	fullyQualified := strings.Join(qualifiedMessage, " ")

	fmt.Println(fullyQualified)
	mac := hmac.New(sha1.New, hashKey)
	mac.Write([]byte(fullyQualified))
	generatedSignature := mac.Sum(nil)

	//nonceExists := nonces[nonce]

	if hmac.Equal([]byte(signature), generatedSignature) == true {

		t.Valid = true
		t.AccessToken = GenerateToken()
		nonces[nonce] = t
		return t, nil
	} else {
		err := errors.New("Unauthorized")
		t.Valid = false
		t.AccessToken = ""
		nonces[nonce] = t
		return t, err
	}

}

func GenerateToken() string {
	var token []byte
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 32; i++ {
		token = append(token, byte(rand.Int63n(74)+48))
	}
	return string(token)
}
