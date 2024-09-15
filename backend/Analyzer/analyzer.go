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
		case "mkdisk": // Comando MKDISK
			processed += Commands.MkDisk(tokens[1:])
		case "rmdisk": // Comando RMDISK
			processed += Commands.RMDisk(tokens[1:])
		case "rep": // Comando REP
			processed += Commands.Rep(tokens[1:])
		case "fdisk": // Comando FDISK
			processed += Commands.FDisk(tokens[1:])
		case "mount": // Comando MOUNT
			processed += Commands.Mount(tokens[1:])
		default: // Comando no reconocido
			processed += "Comando [" + tokens[0] + "] no reconocido\n"
		}
	}

	return processed
}
