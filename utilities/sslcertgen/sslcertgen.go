package sslcertgen

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// GenerateCert will generate an ssl certificate and place in a provided directory path
func GeneratePlaceAndExecute(certPath string) error {
	fmt.Sprintf("Searching for certificate in %s", certPath)
	if generateErr := generate(certPath); generateErr != nil {
		return generateErr
	}

	return nil
}

// generate will automatically create a script that will create the cert
func generate(path string) error {
	gensh := fmt.Sprintf(`openssl req -new -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -x509 -sha256 -days 365 -nodes -out ` + path + `/tls.crt -keyout ` + path + `/tls.key -subj '/C=US/ST=Oregon/L=Portland/CN=www.com.com'`)
	filePath := fmt.Sprintf("%v/%v", path, "gen.sh")

	execPath, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(execPath)

	wdpath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(wdpath) // for example /home/user

	if writeErr := ioutil.WriteFile(filePath, []byte(gensh), 0755); writeErr != nil {
		return writeErr
	}

	cmd := exec.Command("sh", filePath)
	execErr := cmd.Run()

	if execErr != nil {
		return execErr
	}

	return nil
}
