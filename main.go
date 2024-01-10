package main

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
func getKey(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /getkey request\n")

	var myString = "1234567890123456" // 16 bytes

	// Convertir la cadena de texto en un io.Reader
	stringReader := strings.NewReader(myString)
	uuid, err := uuid.NewRandomFromReader(stringReader)

	if err == nil {
		io.WriteString(w, uuid.String())
	} else {
		//io.WriteString(w, "error")
		io.WriteString(w, err.Error())
	}
}
*/

func DownloadImage(url string) ([]byte, error) {
	// Realiza una solicitud GET a la URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Lee la respuesta en una variable
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getKey(w http.ResponseWriter, r *http.Request) {
	var hash [32]byte

	image, err := DownloadImage("https://thispersondoesnotexist.com")
	if err != nil { //Si algo ha ido mal, genero igualmente un hash a partir de un número aleatorio
		// Generar un número aleatorio de 256 bytes
		randomBytes := make([]byte, 256)
		_, err := rand.Read(randomBytes)
		if err != nil {
			panic(err)
		}
		hash = sha256.Sum256(randomBytes)
	} else {
		hash = sha256.Sum256(image)
	}

	fmt.Printf("%x\n", hash)

	if err == nil {
		io.WriteString(w, fmt.Sprintf("%x\n", hash)) //convierto  el array de bytes en string
	} else {
		io.WriteString(w, err.Error())
		io.WriteString(w, fmt.Sprintf("%x\n", hash)) //convierto  el array de bytes en string
	}
}

func main() {
	http.HandleFunc("/getkey", getKey)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
