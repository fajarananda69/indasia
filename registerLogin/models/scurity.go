package models

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ToMd5(string string) string {
	data := []byte(string)
	b := md5.Sum(data)
	pass := hex.EncodeToString(b[:])
	return pass
}

func Encrypt(msg string) string {
	hexEncode := hex.EncodeToString([]byte(msg))
	// fmt.Println("enconding : " + hexEncode)
	base64Encode := base64.URLEncoding.EncodeToString([]byte(hexEncode))
	return base64Encode
}

func Decrypt(msg string) string {

	base64Decode, _ := base64.URLEncoding.DecodeString(msg)
	hexDecode, _ := hex.DecodeString(string(base64Decode))

	return string(hexDecode)
}

// 2FA Authentification ===================
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Append extra 0s if the length of otp is less than 6
//If otp is "1234", it will return it as "001234"
func prefix0(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}

func getHOTPToken(secret string, interval int64) string {

	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	check(err)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	//Signing the value using HMAC-SHA1 Algorithm
	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	o := (h[19] & 15)

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)

	check(err)
	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	h12 := (int(header) & 0x7fffffff) % 1000000

	//Converts number as a string
	otp := strconv.Itoa(int(h12))

	return prefix0(otp)
}

// get token 2fa
func GetTOTPToken(email string) string {
	var secret string
	var pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(email, "")

	if len(processedString) >= 16 {
		secret = processedString[:16]
	} else {
		a := 16 - len(processedString)
		for i := 0; i < a; i++ {
			rand.Seed(time.Now().UnixNano())
			c := pool[rand.Intn(len(pool))]
			processedString += string(c)
		}
		secret = processedString
	}
	//The TOTP token is just a HOTP token seeded with every 30 seconds.
	interval := time.Now().Unix() / 30
	return getHOTPToken(secret, interval)
}

// ====================
