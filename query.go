package browser_cookie_query

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/pbkdf2"
	"net/http"
	"os"
	"time"
)

type BrowserCookieQuery struct {
	browserName string
	db          *sql.DB
	aesBlock    cipher.Block
}

func NewBrowserCookieQuery(browserName string) *BrowserCookieQuery {
	return &BrowserCookieQuery{browserName: browserName}
}

func (q *BrowserCookieQuery) Init() error {
	f, err := getBrowserCookieFile(q.browserName)
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", f)
	if err != nil {
		return err
	}
	q.db = db
	key, err := getBrowserStorageKey(q.browserName)
	if err != nil {
		return err
	}
	encKey := pbkdf2.Key([]byte(key), []byte("saltysalt"), 1003, 16, sha1.New)
	q.aesBlock, err = aes.NewCipher(encKey)
	if err != nil {
		return err
	}
	return nil
}

func (q *BrowserCookieQuery) Query(host string) ([]*http.Cookie, error) {
	rows, err := q.db.Query("SELECT name,encrypted_value,path,expires_utc,is_secure,is_httponly FROM cookies WHERE host_key  = ? and has_expires = 1", host)
	if err != nil {
		return nil, err
	}
	iv := bytes.Repeat([]byte(" "), 16)
	defer rows.Close()
	var cookies []*http.Cookie

	for rows.Next() {
		var name, path string
		var value []byte
		var expiresUtc int64
		var isSecure, isHttponly bool
		if err = rows.Scan(&name, &value, &path, &expiresUtc, &isSecure, &isHttponly); err != nil {
			return nil, err
		}

		// value前三位是版本号 暂时不做校验
		ciphertext := value[3:]
		mode := cipher.NewCBCDecrypter(q.aesBlock, iv)
		mode.CryptBlocks(ciphertext, ciphertext)

		val := bytes.TrimFunc(ciphertext, func(r rune) bool {
			//valid rule copy from net/http/cookie.go validCookieValueByte
			return !(0x20 <= r && r < 0x7f && r != '"' && r != ';' && r != '\\')
		})

		c := &http.Cookie{
			Name:     name,
			Value:    string(val),
			Raw:      string(ciphertext),
			Path:     path,
			Domain:   host,
			Expires:  time.Unix((expiresUtc/1000000)-11644473600, 0),
			Secure:   isSecure,
			HttpOnly: isHttponly,
		}
		cookies = append(cookies, c)
	}
	return cookies, nil
}

func getBrowserCookieFile(browserName string) (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if browserName == "Chrome" {
		browserName = "Google/Chrome"
	}
	return fmt.Sprintf("%s/Library/Application Support/%s/Default/Cookies", homedir, browserName), nil
}

func getBrowserStorageKey(browserName string) (string, error) {
	service := fmt.Sprintf("%s Safe Storage", browserName)
	secret, err := keyring.Get(service, browserName)
	if err != nil {
		return "", err
	}
	return secret, nil
}
