package main

import "fmt"
import "encoding/base64"

import "os"
import "io"
import "crypto/sha1"

//import "github.com/davidgamba/go-getoptions"

func main() {
	fmt.Printf("Hello world!.\n")
	fmt.Printf("%s", generate_encrypted_password("test"))
	_, err := os.OpenFile("/opt/firefox/", os.O_WRONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			//directory does not exist, use default.
			fmt.Print("File Does Not Exist: ")
		}
	}

}

func generate_encrypted_password(password string) string {
	hash := sha1.New()
	io.WriteString(hash, password)
	input := []byte(hash.Sum(nil))
	encoded := base64.StdEncoding.EncodeToString(input)
	return encoded
}
