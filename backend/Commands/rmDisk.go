package commands

import (
	utils "Backend/Utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RMDisk(tokens []string) string {
	// Expresión regular para encontrar los parámetros del comando mkdisk
	Regex := `(?i)-path="[^"]+"|-path=[^\s]+`

	if len(tokens) < 1 {
		return "Error: Comando rmdisk requiere al menos un parametro (path).\n"
	}

	tokens = utils.ParseParametros(tokens, Regex)
	for _, token := range tokens {
		partes := strings.SplitN(token, "=", 2) // Separa el token en partes

		// Obtener solo el nombre del archivo
		fileName := filepath.Base(strings.ReplaceAll(partes[1], "\"", ""))

		// Intentar eliminar el archivo
		err := os.Remove(strings.ReplaceAll(partes[1], "\"", ""))
		if err != nil {
			return fmt.Sprintln("Error al eliminar el archivo:", err)
		} else {
			return fmt.Sprintf("Archivo %s eliminado exitosamente. \n", fileName)
		}
	}

	return "Ocurrio un Problema al detectar rmdisk.\n"
}
