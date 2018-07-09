package common

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

func Obj2Map(obj interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data, nil
}

func Obj2Bytes(bs interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(bf)

	if err := enc.Encode(bs); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func Bytes2Obj(data []byte, bs interface{}) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b) //will read from b
	return dec.Decode(bs)
}

func Gob_Obj2Bytes(t interface{}) []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

func Gob_Bytes2Obj(data []byte, bs interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(bs)
	if err != nil {
		log.Panic(err)
	}
	return err
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func ToJSON(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func DoEvery(d time.Duration, f func()) {
	for x := range time.Tick(d) {
		fmt.Printf("%s: Scheduled mining\n", x.Format(time.RFC3339))
		f()
	}
}

func GenerateRandomString() string {
	randData := make([]byte, 20)
	_, err := rand.Read(randData)
	if err != nil {
		log.Panic(err)
	}

	return fmt.Sprintf("%x", randData)
}
