package commands

import (
	utils "Backend/Utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RMDisk(tokens []string) string {
	var response string
	// Expresión regular para encontrar los parámetros del comando mkdisk
	Regex := `(?i)-path="[^"]+"|-path=[^\s]+`

	// Verificar si el comando tiene al menos un parametro
	if len(tokens) < 1 {
		return "Error: Comando rmdisk requiere al menos un parametro (path).\n"
	}

	// Obtener los parametros del comando
	tokens = utils.ParseParametros(tokens, Regex)
	if len(tokens) < 1 {
		return "Error: No se detecto el path del archivo, por favor verifica la entrada.\n"
	}

	// Iterar sobre los tokens para ejecutar el comando
	for _, token := range tokens {
		partes := strings.SplitN(token, "=", 2) // Separa el token en partes

		// Obtener solo el nombre del archivo
		fileName := filepath.Base(strings.ReplaceAll(partes[1], "\"", ""))

		// Intentar eliminar el archivo
		err := os.Remove(strings.ReplaceAll(partes[1], "\"", ""))
		if err != nil {
			response += fmt.Sprintln("Error al eliminar el archivo:", err)
		} else {
			response += fmt.Sprintf("Archivo %s eliminado exitosamente. \n", fileName)
		}
	}

	return response
}
