package api

import (
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

func saveAndCrop(filename string, w int, h int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { err = file.Close() }()

	// Decodifica l'immagine in un oggetto image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return "", err
	}

	resizedImg := resize.Resize(1000, 0, img, resize.NearestNeighbor)
	filename = strings.Join(strings.Split(filename, "/")[0:2], "/") + "/"
	newFilename := filename + "profilePic_" + fmt.Sprint(w) + "x" + fmt.Sprint(h) + ".jpeg"
	// Salva l'immagine croppata su disco
	out, err := os.Create(newFilename)
	if err != nil {
		return "", err
	}
	defer func() { err = out.Close() }()
	if err := jpeg.Encode(out, resizedImg, nil); err != nil {
		return "", err
	}

	return newFilename, err
}
