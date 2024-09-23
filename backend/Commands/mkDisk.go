package commands

import (
	structs "Backend/Structs"
	utils "Backend/Utils"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func MkDisk(tokens []string) string {
	var path, unit, fit string
	var size, sizeBytes int

	Regex := `(?i)-size=\d+|-unit=[^\s]|-fit=[^\s]{2}|-path="[^"]+"|-path=[^\s]+|-param=[^\s]+`

	tokens = utils.ParseParametros(tokens, Regex)

	if len(tokens) > 4 || len(tokens) < 2 { // Si no hay cuatro parametros, no hace nada
		return "Error: Comando mkdisk requiere minimo 2 parametros (path, size)  y maximo 4 (path, size, fit, unit).\n"
	}

	for _, token := range tokens { // Itera sobre los tokens para obtener los parametros
		partes := strings.SplitN(token, "=", 2) // Separa el token en partes
		if len(partes) != 2 {
			return fmt.Sprintf("formato de parámetro inválido: %s", token)
		}

		switch strings.ToLower(partes[0]) { // Switch para manejar los parametros
		case "-size":
			size, _ = strconv.Atoi(partes[1])

			if size < 0 {
				return "Error: El tamaño debe ser mayor a cero.\n"
			}
		case "-unit":
			unit = strings.ToUpper(partes[1])

			if strings.ToUpper(unit) != "K" && strings.ToUpper(unit) != "M" {
				return "Error: Unidad de medida no reconocida, debe ser K o M.\n"
			}
		case "-fit":
			fit = strings.ToUpper(partes[1])

			if fit != "BF" && fit != "FF" && fit != "WF" {
				return "Error: Tipo de ajuste no reconocido, el ajuste debe ser BF, FF o WF.\n"
			}
		case "-path":
			path = strings.ReplaceAll(partes[1], "\"", "")
		default:
			return fmt.Sprintf("Error: Parametro [" + partes[0] + "] no reconocido.\n")
		}

	}

	// Manejo de erorres en los parametros
	if size == 0 {
		return "Error: Tamaño no especificado.\n"
	}

	if path == "" {
		return "Error: Path no especificado.\n"
	}

	if fit == "" {
		fit = "FF"
	}

	if unit == "" {
		unit = "M"
	}

	// Calcular el tamaño en bytes
	sizeBytes = utils.CalcularTamaño(size, unit)
	if sizeBytes == -1 { // Si el tamaño es -1, la unidad no es reconocida
		return "Error: Unidad de medida no reconocida.\n"
	}

	// Crear el MBR
	var disco structs.MBR

	// Inicializar el MBR
	disco.Mbr_size = int32(sizeBytes)           // Tamaño del disco en bytes
	disco.Mbr_date = float64(time.Now().Unix()) // Obtiene la fecha actual en formato Unix
	disco.Mbr_signature = rand.Int31()          // Genera un número aleatorio de tipo int32

	disco.Mbr_fit = [2]byte{fit[0], fit[1]} // Tipo de ajuste
	disco.Mbr_partitions = [4]structs.Partition{
		{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: 0, Part_id: [4]byte{'0'}},
		{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: 0, Part_id: [4]byte{'0'}},
		{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: 0, Part_id: [4]byte{'0'}},
		{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: 0, Part_id: [4]byte{'0'}},
	}

	// Extraer el directorio de la ruta del archivo
	dir := filepath.Dir(path)

	// Crear todas las carpetas necesarias, si no existen
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Sprintf("Error al crear directorios: %v\n", err)
	}

	// Crear el archivo
	file, err := os.Create(path)
	if err != nil {
		return fmt.Sprintf("Error al crear el archivo: %v\n", err)
	}
	defer file.Close()

	// Escribir ceros en el archivo para alcanzar el tamaño deseado
	_, err = file.Write(make([]byte, sizeBytes))
	if err != nil {
		return fmt.Sprintf("Error al completar el archivo con ceros: %v\n", err)
	}

	disco.SerializeMBR(path)

	return "Disco creado con exito.\n"
}
