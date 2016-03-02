package main

import "fmt"
import "encoding/base64"

import "os"
import "io"
import "strings"
import "path/filepath"

//import "bufio"
import "crypto/sha1"
import "regexp"

import "io/ioutil"
import "github.com/davidgamba/go-getoptions"
import "github.com/davidgamba/go-utils/fileutils"

var user_file_dir = string(filepath.Separator) + "tmp" + string(filepath.Separator) + "CSI"
var user_service_file = "user-service.properties"

func main() {

	opt := getoptions.GetOptions()
	var list bool
	var dir string
	var user string
	var password string
	var role string
	var delete string

	opt.BoolVar(&list, "list", false, "l", "panda")
	opt.StringVar(&dir, "dir", "")
	opt.StringVar(&user, "user", "", "u")
	opt.StringVar(&password, "password", "", "p")
	opt.StringVar(&role, "role", "enduser", "r")
	opt.StringVar(&delete, "delete", "", "d")

	remaining, err := opt.Parse(os.Args[1:])

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

	if remaining != nil {
		fmt.Printf("Remaining: %s\n, remaining")
	}

	if opt.Called["dir"] {
		user_file_dir = dir
	}

	writer, err := os.OpenFile(user_file_dir+string(filepath.Separator)+user_service_file, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println(err)
		if os.IsNotExist(err) {
			fmt.Printf("[INFO] File Doesn't Exist\n")
			os.MkdirAll(user_file_dir, 0777)
			writer, err := os.Create(user_file_dir + string(filepath.Separator) + user_service_file)
			check(err)
			writer.Close()
		}
	}
	writer.Close()

	fmt.Println("[INFO] Using the following user file: " + user_file_dir + string(filepath.Separator) + user_service_file)

	if opt.Called["user"] {
		// check that password is also passed through
		if opt.Called["password"] {
			//fmt.Println(user_string)
			add_user_to_user_service_file(user, password, role, true)
		} else {
			fmt.Println("[INFO] Password must also be set.")
		}
	}

	if opt.Called["delete"] {
		remove_user_from_user_service_file(delete)
	}
	if opt.Called["list"] {
		get_users_from_user_service_file()
	}
}

func get_users_from_user_service_file() {
	input, err := ioutil.ReadFile(user_file_dir + string(filepath.Separator) + user_service_file)
	check(err)
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		username := strings.Split(string(line), "=")
		fmt.Println(string(username[0]))
	}
}

func remove_user_from_user_service_file(username string) {

	var user_found bool = false
	input, err := ioutil.ReadFile(user_file_dir + "/" + user_service_file)
	check(err)

	lines := strings.Split(string(input), "\n")
	result := []string{}

	for i, line := range lines {
		if !strings.Contains(line, username) {
			result = append(result, lines[i])
		} else {
			user_found = true
		}
	}

	if user_found {
		output := strings.Join(result, "\n")
		err = ioutil.WriteFile("/tmp/CSI/user-service.properties2", []byte(output), 0644)
		check(err)

		err = fileutils.CopyFile("/tmp/CSI/user-service.properties2", "/tmp/CSI/user-service.properties")
		check(err)
	} else {
		fmt.Println("[WARNING] user " + username + " was not found in user file")
	}

}

func add_user_to_user_service_file(username string, password string, user_type string, encrypt bool) {
	buf, err := ioutil.ReadFile(user_file_dir + string(filepath.Separator) + user_service_file)
	check(err)
	s := string(buf)
	regex, err := regexp.Compile(username)
	check(err)
	if regex.MatchString(s) {
		panic("User already exists: " + username)
	}

	if encrypt {
		password = generate_encrypted_password(password)
	}

	var role_admin = "ROLE_PUREWEB_SERVER_ADMIN,ROLE_PUREWEB_USER_ADMIN"
	var role_monitor = "ROLE_PUREWEB_SERVER_MONITOR"
	var role_user = "ROLE_PUREWEB_USER"
	var status = "enabled"

	var user_string string
	switch user_type {
	case "admin":
		user_string = username + "=" + password + "," + role_admin + "," + role_user + "," + status + "\n"
	case "monitor":
		user_string = username + "=" + password + "," + role_monitor + "," + status + "\n"
	case "enduser":
		user_string = username + "=" + password + "," + role_user + "," + status + "\n"
	default:
		panic("Invalid user type: " + user_type)
	}
	fmt.Println("[INFO] Adding user: " + username)
	file, err := os.OpenFile(user_file_dir+string(filepath.Separator)+user_service_file, os.O_RDWR|os.O_APPEND, 0600)
	check(err)
	n3, err := io.WriteString(file, user_string)
	fmt.Printf("wrote %d bytes\n", n3)
	file.Sync()
	defer file.Close()
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func generate_encrypted_password(password string) string {
	hash := sha1.New()
	io.WriteString(hash, password)
	input := []byte(hash.Sum(nil))
	encoded := base64.StdEncoding.EncodeToString(input)
	return encoded
}
