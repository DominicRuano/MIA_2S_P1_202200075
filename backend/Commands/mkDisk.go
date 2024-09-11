package commands

import (
	structs "Backend/Structs"
	utils "Backend/Utils"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func MkDisk(tokens []string) string {
	var path, unit, fit string
	var size, sizeBytes int

	Regex := `(?i)-size=\d+|-unit=[kKmM]|-fit=[bBfFwW]{2}|-path="[^"]+"|-path=[^\s]+`

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
			unit = partes[1]

			if strings.ToUpper(unit) != "K" && strings.ToUpper(unit) != "M" {
				return "Error: Unidad de medida no reconocida, debe ser K o M.\n"
			}
		case "-fit":
			fit = partes[1]

			if strings.ToUpper(fit) != "BF" && strings.ToUpper(fit) != "FF" && strings.ToUpper(fit) != "WF" {
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
	disco.Mbr_size = int32(sizeBytes)
	disco.Mbr_date = float64(time.Now().Unix())
	disco.Mbr_signature = int32(rand.Intn(1000))

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

	// Escribir el MBR en el archivo
	err = binary.Write(file, binary.LittleEndian, &disco)
	if err != nil {
		return fmt.Sprintf("Error al escribir el MBR en el archivo: %v\n", err)
	}

	// Calcular cuántos bytes ocupa el MBR
	sizeOfMBR := int(unsafe.Sizeof(disco))

	// Escribir ceros en el archivo para alcanzar el tamaño deseado
	// Aquí restamos el tamaño del MBR del total para llenar el espacio restante
	_, err = file.Write(make([]byte, sizeBytes-sizeOfMBR))
	if err != nil {
		return fmt.Sprintf("Error al completar el archivo con ceros: %v\n", err)
	}

	return "Disco creado con exito.\n"
}
