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

	//Borrar de aqui...
	// Crea un archivo para escribir la imagen
	file, err := os.Create("imagen.jpg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Copia el contenido de la variable 'body' al archivo
	_, err = file.Write(body)
	if err != nil {
		return nil, err
	}
	//A aqui si no se quiere guardar la imagen
	// Imprime el tamaño de body
	//fmt.Printf("Tamaño de body: %d bytes\n", len(body))

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
		io.WriteString(w, err.Error())               //Muestro el error que dió al intentar descargar la imagen
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
