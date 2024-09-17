package utils

import (
	structs "Backend/Structs"
	"encoding/binary"
	"math"
	"regexp"
	"strings"
)

func CalcularTamaño(size int, unit string) int {
	if unit == "B" {
		return size
	} else if unit == "K" {
		return size * 1024
	} else if unit == "M" {
		return size * 1024 * 1024
	} else {
		return -1
	}
}

func ParseParametros(tokens []string, RE string) []string {
	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")

	// Expresión regular para encontrar los parámetros del comando mkdisk
	re := regexp.MustCompile(RE)

	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	tokens = re.FindAllString(args, -1)

	return tokens
}

func CalculateN(partition *structs.Partition) int32 {
	/*
		numerador = (partition_montada.size - sizeof(Structs::Superblock)
		denrominador base = (4 + sizeof(Structs::Inodes) + 3 * sizeof(Structs::Fileblock))
		n = floor(numerador / denrominador)
	*/

	numerator := int(partition.Part_size) - binary.Size(structs.SuperBloque{})
	denominator := 4 + binary.Size(structs.Inodo{}) + 3*binary.Size(structs.FileBlock{}) // No importa que bloque poner, ya que todos tienen el mismo tamaño
	n := math.Floor(float64(numerator) / float64(denominator))

	return int32(n)
}
