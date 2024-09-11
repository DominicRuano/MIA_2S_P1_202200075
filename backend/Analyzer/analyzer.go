package analyzer

import (
	Commands "Backend/Commands"
	"strings"
)

// Función De analizador.
func Analyzer(text string) string {
	Lineas := strings.Split(text, "\n")
	processed := ""

	for _, linea := range Lineas {
		linea = strings.TrimSpace(linea)
		tokens := strings.Split(linea, " ")

		// Verificar si la linea es una linea vacia
		if len(tokens) == 0 || tokens[0] == "" {
			processed += "\n"
			continue
		}

		// Verificar si la linea es un comentario
		if tokens[0][0] == '#' {
			processed += linea[1:] + "\n"
			continue
		}

		// Verificar si la linea es un comando

		switch tokens[0] {
		case "#": // Comentario
			processed += strings.Join(tokens[1:], " ") + "\n"
		case "MKDISK": // Comando MKDISK
			processed += Commands.MkDisk(tokens)
		case "RMDISK": // Comando RMDISK
			processed += Commands.RMDisk(tokens)
		default: // Comando no reconocido
			processed += "Comando [" + tokens[0] + "] no reconocido\n"
		}
	}

	return processed
}
