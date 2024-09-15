package utils

import (
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
