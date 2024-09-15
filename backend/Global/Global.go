package global

import (
	"errors"
)

var (
	MountedPartitions map[string]string = make(map[string]string)
)

// Ultimos 2 digitos del carnet
const Carnet string = "75"

// Lista de todo el abecedario
var Abecedario = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

// Mapa para guardar las particiones con letras.
var PathToletter = make(map[string]string)

// Indice para la siguiente letra.
var NextLetterindex int = 0

// GetLetter obtiene la letra asignada a un path
func GetLetter(path string) (string, error) {
	// Asignar una letra al path si no tiene una asignada
	if _, exists := PathToletter[path]; !exists {
		if NextLetterindex < len(Abecedario) {
			PathToletter[path] = Abecedario[NextLetterindex]
			NextLetterindex++
		} else {
			return "", errors.New("no hay más letras disponibles para asignar")
		}
	}

	return PathToletter[path], nil
}

// Mapa para guardar las particiones con numeros.
var PathToNumber = make(map[string]int)

// GetLetter obtiene el numero para el path
func GetNumPartition(path string) int {
	// Asignar una letra al path si no tiene una asignada
	if _, exists := PathToNumber[path]; !exists {
		PathToNumber[path] = 0
	} else {
		PathToNumber[path] += 1
	}

	return PathToNumber[path]
}
