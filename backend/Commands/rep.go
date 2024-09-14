package commands

import (
	structs "Backend/Structs"
	utils "Backend/Utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Rep(tokens []string) string {
	var name, path string
	// Expresión regular para encontrar los parámetros del comando mkdisk
	Regex := `(?i)-name=[^\s]+|-path="[^"]+"|-path=[^\s]+`

	// Verificar si el comando tiene al menos un parametro
	if len(tokens) < 1 {
		return "Error: Comando rep requiere al menos X parametros (x,x,x,x).\n"
	}

	// Obtener los parametros del comando
	tokens = utils.ParseParametros(tokens, Regex)
	if len(tokens) < 1 {
		return "Error: No se detectaron los parametros, por favor verifica la entrada.\n"
	}

	for _, token := range tokens { // Itera sobre los tokens para obtener los parametros
		partes := strings.SplitN(token, "=", 2) // Separa el token en partes
		if len(partes) != 2 {
			return fmt.Sprintf("formato de parámetro inválido: %s", token)
		}

		switch strings.ToLower(partes[0]) { // Switch para manejar los parametros
		case "-name":
			name = partes[1]
		case "-path":
			path = strings.ReplaceAll(partes[1], "\"", "")
		default:
			return fmt.Sprintf("Error: Parametro [" + partes[0] + "] no reconocido.\n")
		}
	}

	if name == "" {
		return "Error: Faltan el parametro nombre.\n"
	}

	if path == "" {
		return "Error: Faltan el path.\n"
	}

	// Crea una variable para almacenar los datos leídos
	var mbr structs.MBR

	// Obtiene el mbr del archivo
	mbr.DeserializeMBR(path)

	// Path final del archivo sin la extensión
	finalPath := "../Reportes/" + name + "_reporte"

	// Generar el archivo DOT
	err := GenerateDotFile(mbr, finalPath+".dot")
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Ejecutar el comando DOT
	err = ExecuteDot(finalPath+".dot", finalPath+".png")
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Intentar eliminar el archivo
	//os.Remove(finalPath + ".dot")

	return "Comando REP ejecutado correctamente.\n"
}

// Función para generar el archivo DOT con la estructura del MBR
func GenerateDotFile(mbr structs.MBR, filePath string) error {
	// Crear el archivo .dot
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo DOT: %v", err)
	}
	defer file.Close()

	// Escribir el contenido del archivo DOT
	_, err = file.WriteString(`
digraph G {
    node [shape=plaintext];
    tabla_mbr [
        label=<
        <table border="0" cellborder="1" cellspacing="0" cellpadding="10">
            <tr><td colspan="2" bgcolor="purple"><font color="white">REPORTE DE MBR</font></td></tr>
            <tr><td>mbr_tamano</td><td>` + fmt.Sprint(mbr.Mbr_size) + `</td></tr>
            <tr><td>mbr_fecha_creacion</td><td>` + time.Unix(int64(mbr.Mbr_date), 0).Format("2006-01-02 15:04") + `</td></tr>
            <tr><td>mbr_disk_signature</td><td>` + fmt.Sprint(mbr.Mbr_signature) + `</td></tr>
            <tr><td>mbr_disk_fit</td><td>` + fmt.Sprint(mbr.Mbr_fit[0]) + `</td></tr>
    `)

	if err != nil {
		return fmt.Errorf("error al escribir en el archivo DOT: %v", err)
	}

	// Escribir las particiones en el archivo DOT
	for i, partition := range mbr.Mbr_partitions {
		partName := ConvertirNombre(partition.Part_name)
		_, err := file.WriteString(`
            <tr><td colspan="2" bgcolor="purple"><font color="white">Particion ` + fmt.Sprint(i+1) + `</font></td></tr>
			<tr><td>part_status</td><td>` + fmt.Sprint(partition.Part_status[0]) + `</td></tr>
            <tr><td>part_type</td><td>` + fmt.Sprint(partition.Part_type[0]) + `</td></tr>
            <tr><td>part_fit</td><td>` + fmt.Sprint(partition.Part_fit[0]) + `</td></tr>
            <tr><td>part_start</td><td>` + fmt.Sprint(partition.Part_start) + `</td></tr>
            <tr><td>part_size</td><td>` + fmt.Sprint(partition.Part_size) + `</td></tr>
            <tr><td>part_name</td><td>` + partName + `</td></tr>
            <tr><td>part_correlative</td><td>` + fmt.Sprint(partition.Part_correlative) + `</td></tr>
            <tr><td>part_id</td><td>` + fmt.Sprint(partition.Part_id) + `</td></tr>
        `)

		if err != nil {
			return fmt.Errorf("error al escribir las particiones en el archivo DOT: %v", err)
		}
	}

	// Cerrar la tabla
	_, err = file.WriteString(`
        </table>
        >];
}
    `)

	if err != nil {
		return fmt.Errorf("error al cerrar la tabla en el archivo DOT: %v", err)
	}

	return nil
}

// Función para ejecutar el comando dot para generar la imagen
func ExecuteDot(dotFile, outputImage string) error {
	// Ejecutar el comando `dot` para convertir .dot en .png
	cmd := exec.Command("dot", "-Tpng", dotFile, "-o", outputImage)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando dot: %v", err)
	}
	return nil
}

// Convertir Part_name a string ignorando los caracteres nulos (0)
func ConvertirNombre(partName [16]byte) string {
	var name []byte
	for _, b := range partName {
		if b == 0 {
			break // Ignorar el resto si encontramos un valor nulo
		}
		name = append(name, b)
	}
	return string(name)
}
