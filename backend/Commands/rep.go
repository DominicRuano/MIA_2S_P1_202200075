package commands

import (
	global "Backend/Global"
	structs "Backend/Structs"
	utils "Backend/Utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
)

type Rep_st struct {
	Rep_name         string
	Rep_pathDisk     string
	Rep_pathDestino  string
	Rep_id           string
	Rep_path_file_ls string
}

func (rep *Rep_st) Print() {
	fmt.Println("Rep_name:", rep.Rep_name)
	fmt.Println("Rep_pathDisk:", rep.Rep_pathDisk)
	fmt.Println("Rep_pathDestino:", rep.Rep_pathDestino)
	fmt.Println("Rep_id:", rep.Rep_id)
	fmt.Println("Rep_path_file_ls:", rep.Rep_path_file_ls)
}

func Rep(tokens []string) string {
	// Variables para almacenar los parametros del comando
	Cmd := &Rep_st{}

	// Expresión regular para encontrar los parámetros del comando rep
	Regex := `(?i)-name=[^\s]+|-name="[^"]+"|-path="[^"]+"|-path=[^\s]+|-id=[^\s]+|-id="[^"]+"|-path_file_ls="[^"]+"|-path_file_ls=[^\s]+`

	// Verificar si el comando tiene los parametros necesarios
	if len(tokens) < 3 {
		return "Error: Comando rep requiere al menos 3 parametros (name, path, id).\n"
	}
	if len(tokens) > 4 {
		return "Error: Comando rep requiere maximo 4 parametros (name, path, id, path_file_ls).\n"
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
			Cmd.Rep_name = partes[1]
		case "-path":
			Cmd.Rep_pathDestino = strings.ReplaceAll(partes[1], "\"", "")
		case "-id":
			if PathOrigen, exists := global.MountedPartitions[partes[1]]; exists {
				Cmd.Rep_id = partes[1]
				Cmd.Rep_pathDisk = PathOrigen
			} else {
				return fmt.Sprintf("Error: La particion %s no esta montada.\n", partes[1])
			}
		case "-path_file_ls":
			Cmd.Rep_path_file_ls = strings.ReplaceAll(partes[1], "\"", "")
		default:
			return fmt.Sprintf("Error: Parametro [" + partes[0] + "] no reconocido.\n")
		}
	}

	// Verificar los parametros
	if Cmd.Rep_name == "" {
		return "Error: Faltan el parametro nombre.\n"
	}
	if Cmd.Rep_pathDestino == "" {
		return "Error: Faltan el path.\n"
	}
	if Cmd.Rep_id == "" {
		return "Error: Faltan el id.\n"
	}

	// Ejecutar el comando rep
	switch strings.ToLower(Cmd.Rep_name) {
	case "mbr":
		return MBRReporte(Cmd)
	case "disk":
		return DISKReporte(Cmd)
	case "bm_inode":
		return BMInodeReporte(Cmd)
	case "bm_block":
		return BMBloqueReporte(Cmd)
	case "sb":
		return sbReporte(Cmd)
	default:
		return fmt.Sprintf("Error: Nombre de reporte %s no reconocido.\n", Cmd.Rep_name)
	}

}

func sbReporte(Cmd *Rep_st) string {
	// Deserializa el mbr
	mbr := &structs.MBR{}

	// Lee el archivo binario
	err := mbr.DeserializeMBR(Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Busca la particion
	Particion := mbr.GetPartitionById(Cmd.Rep_id)
	if Particion == nil {
		return fmt.Sprintf("Error: Particion %s no encontrada.\n", Cmd.Rep_id)
	}

	// Obtiene el ofcet del superbloque
	offsetSuperbloque := Particion.Part_start

	// Crea una superblock para almacenar los datos leídos
	superblock := &structs.SuperBloque{}

	// Lee el superbloque
	superblock.DeserializeSB(Cmd.Rep_pathDisk, int64(offsetSuperbloque))

	// Obtener la extensión del archivo
	ext := filepath.Ext(Cmd.Rep_pathDestino)

	// Path final del archivo sin la extensión
	finalPath := strings.TrimSuffix(Cmd.Rep_pathDestino, ext)

	// Obtener el directorio del path (sin el archivo)
	dir := filepath.Dir(Cmd.Rep_pathDestino)

	// Crear todas las carpetas necesarias si no existen
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Sprintf("Error al crear directorios: %v\n", err)
	}

	// Generar el archivo DOT
	err = GenerateDotFileSB(superblock, finalPath+".dot", Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Ejecutar el comando DOT
	err = ExecuteDot(finalPath+".dot", Cmd.Rep_pathDestino)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Intentar eliminar el archivo
	os.Remove(finalPath + ".dot")

	// Mensaje de éxito
	return fmt.Sprintf("REP: Reporte SB Para el Id %s generado con éxito.\n", Cmd.Rep_id)

}

// Función para generar el archivo DOT con la estructura del MBR
func GenerateDotFileSB(sb *structs.SuperBloque, filePath string, discPath string) error {

	// Crear el archivo .dot
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo DOT: %v", err)
	}
	defer file.Close()

	disconame := filepath.Base(discPath)

	// Escribir el contenido del archivo DOT
	_, err = file.WriteString(`
digraph G {
	node [shape=plaintext];
	A [label=< 
		<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4">
			<TR>
				<TD COLSPAN="2" BGCOLOR="green"><B>Reporte de SUPERBLOQUE</B></TD>
			</TR>
			<TR>
				<TD>Disco</TD>
				<TD>` + disconame + `</TD>
			</TR>
			<TR>
				<TD>s_filesystem_type</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_filesystem_type) + `</TD>
			</TR>
			<TR>
				<TD>s_inodes_count</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_inodes_count) + `</TD>
			</TR>
			<TR>
				<TD>s_blocks_count</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_blocks_count) + `</TD>
			</TR>
			<TR>
				<TD>s_free_blocks_count</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_free_blocks_count) + `</TD>
			</TR>
			<TR>
				<TD>s_free_inodes_count</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_free_inodes_count) + `</TD>
			</TR>
			<TR>
				<TD>s_mtime</TD>
				<TD>` + time.Unix(int64(sb.SB_mtime), 0).Format("2006-01-02 15:04") + `</TD>
			</TR>
			<TR>
				<TD>s_umtime</TD>
				<TD>` + time.Unix(int64(sb.SB_umtime), 0).Format("2006-01-02 15:04") + `</TD>
			</TR>
			<TR>
				<TD>s_mnt_count</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_mnt_count) + `</TD>
			</TR>
			<TR>
				<TD>s_magic</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_magic) + `</TD>
			</TR>
			<TR>
				<TD>s_inode_s</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_inode_size) + `</TD>
			</TR>
			<TR>
				<TD>s_block_s</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_block_size) + `</TD>
			</TR>
			<TR>
				<TD>s_first_ino</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_firs_ino) + `</TD>
			</TR>
			<TR>
				<TD>s_first_blo</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_first_blo) + `</TD>
			</TR>
			<TR>
				<TD>s_bm_inode_start</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_bm_inode_start) + `</TD>
			</TR>
			<TR>
				<TD>s_bm_block_start</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_bm_block_start) + `</TD>
			</TR>
			<TR>
				<TD>s_inode_start</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_inode_start) + `</TD>
			</TR>
			<TR>
				<TD>s_block_start</TD>
				<TD>` + fmt.Sprintf("%d", sb.SB_block_start) + `</TD>
			</TR>
			`)

	if err != nil {
		return fmt.Errorf("error al escribir en el archivo DOT: %v", err)
	}

	// Cerrar la tabla
	_, err = file.WriteString(`
        </TABLE>
        >];
}
    `)
	if err != nil {
		return fmt.Errorf("error al cerrar la tabla en el archivo DOT: %v", err)
	}

	return nil
}

func BMBloqueReporte(Cmd *Rep_st) string {
	// Deserializa el mbr
	mbr := &structs.MBR{}

	// Lee el archivo binario
	err := mbr.DeserializeMBR(Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Busca la particion
	Particion := mbr.GetPartitionById(Cmd.Rep_id)
	if Particion == nil {
		return fmt.Sprintf("Error: Particion %s no encontrada.\n", Cmd.Rep_id)
	}

	// Obtiene el ofcet del superbloque
	offsetSuperbloque := Particion.Part_start

	// Crea una superblock para almacenar los datos leídos
	superblock := &structs.SuperBloque{}

	// Lee el superbloque
	superblock.DeserializeSB(Cmd.Rep_pathDisk, int64(offsetSuperbloque))

	// Obtiene el ofcet del bitmap de inodos
	offsetBitmapInodos := superblock.SB_bm_block_start

	// Calcula el tamaño del bitmap
	size := superblock.SB_free_blocks_count + superblock.SB_blocks_count

	// Lee el bitmap de inodos
	// Crea un bitmap para almacenar los datos leídos
	bitmapInodos, err := utils.ReadBytes(Cmd.Rep_pathDisk, int64(offsetBitmapInodos), size)

	// Verifica si hubo un error al leer el bitmap
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Crea un archivo para almacenar el bitmap de inodos
	file, err := os.Create(Cmd.Rep_pathDestino)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	counter := 0
	for _, b := range bitmapInodos {
		var charToWrite string

		// Verificar si el byte es 0 o 1 y escribir el carácter correspondiente
		if fmt.Sprintf("%d", b) == "88" {
			charToWrite = "1"
		} else if fmt.Sprintf("%d", b) == "79" {
			charToWrite = "0"
		} else {
			// Manejar cualquier otro valor que no sea 0 o 1 si es necesario
			charToWrite = "?"
		}

		// Escribir el carácter en el archivo
		_, err = file.WriteString(charToWrite)
		if err != nil {
			return fmt.Sprintf("Error: %s\n", err)
		}

		counter++

		// Cada 20 bytes, agregar un salto de línea
		if counter%20 == 0 {
			_, err = file.WriteString("\n")
			if err != nil {
				return fmt.Sprintf("Error: %s\n", err)
			}
		}
	}

	// Añadir un salto de línea al final, si es necesario (opcional)
	if counter%20 != 0 {
		_, err = file.WriteString("\n")
		if err != nil {
			return fmt.Sprintf("Error: %s\n", err)
		}
	}

	return fmt.Sprintf("REP: Reporte BM Inode Para el Id %s generado con éxito.\n", Cmd.Rep_id)
}

func BMInodeReporte(Cmd *Rep_st) string {
	// Deserializa el mbr
	mbr := &structs.MBR{}

	// Lee el archivo binario
	err := mbr.DeserializeMBR(Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Busca la particion
	Particion := mbr.GetPartitionById(Cmd.Rep_id)
	if Particion == nil {
		return fmt.Sprintf("Error: Particion %s no encontrada.\n", Cmd.Rep_id)
	}

	// Obtiene el ofcet del superbloque
	offsetSuperbloque := Particion.Part_start

	// Crea una superblock para almacenar los datos leídos
	superblock := &structs.SuperBloque{}

	// Lee el superbloque
	superblock.DeserializeSB(Cmd.Rep_pathDisk, int64(offsetSuperbloque))

	// Obtiene el ofcet del bitmap de inodos
	offsetBitmapInodos := superblock.SB_bm_inode_start

	// Calcula el tamaño del bitmap
	size := superblock.SB_free_inodes_count + superblock.SB_inodes_count

	// Lee el bitmap de inodos
	// Crea un bitmap para almacenar los datos leídos
	bitmapInodos, err := utils.ReadBytes(Cmd.Rep_pathDisk, int64(offsetBitmapInodos), size)

	// Verifica si hubo un error al leer el bitmap
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Crea un archivo para almacenar el bitmap de inodos
	file, err := os.Create(Cmd.Rep_pathDestino)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	counter := 0
	for _, b := range bitmapInodos {
		var charToWrite string

		// Verificar si el byte es 0 o 1 y escribir el carácter correspondiente
		if fmt.Sprintf("%d", b) == "49" {
			charToWrite = "1"
		} else if fmt.Sprintf("%d", b) == "48" {
			charToWrite = "0"
		} else {
			// Manejar cualquier otro valor que no sea 0 o 1 si es necesario
			charToWrite = "?"
		}

		// Escribir el carácter en el archivo
		_, err = file.WriteString(charToWrite)
		if err != nil {
			return fmt.Sprintf("Error: %s\n", err)
		}

		counter++

		// Cada 20 bytes, agregar un salto de línea
		if counter%20 == 0 {
			_, err = file.WriteString("\n")
			if err != nil {
				return fmt.Sprintf("Error: %s\n", err)
			}
		}
	}

	// Añadir un salto de línea al final, si es necesario (opcional)
	if counter%20 != 0 {
		_, err = file.WriteString("\n")
		if err != nil {
			return fmt.Sprintf("Error: %s\n", err)
		}
	}

	return fmt.Sprintf("REP: Reporte BM Inode Para el Id %s generado con éxito.\n", Cmd.Rep_id)
}

func DISKReporte(Cmd *Rep_st) string {
	// Crea una variable para almacenar el disco.
	var mbr structs.MBR

	// Obtiene el mbr del archivo
	mbr.DeserializeMBR(Cmd.Rep_pathDisk)

	// Obtener la extensión del archivo
	ext := filepath.Ext(Cmd.Rep_pathDestino)

	// Path final del archivo sin la extensión
	finalPath := strings.TrimSuffix(Cmd.Rep_pathDestino, ext)

	// Obtener el directorio del path (sin el archivo)
	dir := filepath.Dir(Cmd.Rep_pathDestino)

	// Crear todas las carpetas necesarias si no existen
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Sprintf("Error al crear directorios: %v\n", err)
	}

	// Generar el archivo DOT
	err = GenerateDotFileDISK(mbr, finalPath+".dot", Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Ejecutar el comando DOT
	err = ExecuteDot(finalPath+".dot", Cmd.Rep_pathDestino)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Intentar eliminar el archivo
	os.Remove(finalPath + ".dot")

	// Mensaje de éxito
	return fmt.Sprintf("REP: Reporte Disk Para el Id %s generado con éxito.\n", Cmd.Rep_id)
}

func MBRReporte(Cmd *Rep_st) string {
	// Crea una variable para almacenar los datos leídos
	var mbr structs.MBR

	// Obtiene el mbr del archivo
	err := mbr.DeserializeMBR(Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintf("Error al obtener el MBR: %v\n", err)
	}

	// Obtener la extensión del archivo
	ext := filepath.Ext(Cmd.Rep_pathDestino)

	// Path final del archivo sin la extensión
	finalPath := strings.TrimSuffix(Cmd.Rep_pathDestino, ext)

	// Obtener el directorio del path (sin el archivo)
	dir := filepath.Dir(Cmd.Rep_pathDestino)

	// Crear todas las carpetas necesarias si no existen
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Sprintf("Error al crear directorios: %v\n", err)
	}

	// Generar el archivo DOT
	err = GenerateDotFileMBR(mbr, finalPath+".dot", Cmd.Rep_pathDisk)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Ejecutar el comando DOT
	err = ExecuteDot(finalPath+".dot", Cmd.Rep_pathDestino)
	if err != nil {
		return fmt.Sprintln("Error:", err)
	}

	// Intentar eliminar el archivo
	os.Remove(finalPath + ".dot")

	// Mensaje de éxito
	return fmt.Sprintf("REP: Reporte MBR Para el Id %s generado con éxito.\n", Cmd.Rep_id)
}

// Función para generar el archivo DOT con la estructura del MBR
func GenerateDotFileDISK(mbr structs.MBR, filePath string, discPath string) error {
	// Crear el archivo .dot
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo DOT: %v", err)
	}
	defer file.Close()

	disconame := filepath.Base(discPath)

	// Calcular el Porcentaje en BYTES que ocupa la partición en el disco
	mbrWeight := float64(unsafe.Sizeof(mbr)) * 100 / float64(mbr.Mbr_size)

	// Escribir el contenido del archivo DOT
	_, err = file.WriteString(`
digraph G {
    node [shape=plaintext];
    A [label=<
        <TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4">
            <TR>
                <TD COLSPAN="999" BORDER="0"><B>` + disconame + `</B></TD>
            </TR>
            <TR>
				<TD ROWSPAN="2">MBR<BR/><FONT POINT-SIZE="7" COLOR="blue">` + fmt.Sprintf("%.2f", mbrWeight) + `%` + ` del disco</FONT></TD>`)

	if err != nil {
		return fmt.Errorf("error el escribir el tamaño del MBR en el archivo DOT: %v", err)
	}

	// Indice de la partición extendida
	ExtendidaIndex := -1

	// Indice de la ultima partición
	LastIndex := -1

	// Escribir las particiones en el archivo DOT
	for i, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] == 'P' {

			// Actualiza el indice de la ultima partición
			LastIndex = i

			// Calcular el Porcentaje en BYTES que ocupa la partición en el disco
			partWeight := float64(partition.Part_size) * 100 / float64(mbr.Mbr_size)

			// Escribir la partición en el archivo DOT
			_, err := file.WriteString(`
				<TD ROWSPAN="2">Primaria<BR/><FONT POINT-SIZE="10" COLOR="red">` + fmt.Sprintf("%.2f", partWeight) + `%` + ` del disco</FONT></TD>`)
			if err != nil {
				return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
			}
		} else if partition.Part_type[0] == 'E' {

			// Actualiza el indice de la ultima partición
			LastIndex = i

			// Guarda el indice de la partición extendida
			ExtendidaIndex = i

			// Calcula el ROWSPAN de la particion extendida
			ROWSPAN := 2

			// Crear una variable para almacenar los datos leídos
			var ebr structs.EBR

			// Obtiene el ebr del archivo
			ebr.DeserializeEBR(discPath, int64(partition.Part_start))

			for ebr.Ebr_next != -1 {
				ROWSPAN += 2
				// Obtiene el proximo ebr del archivo
				ebr.DeserializeEBR(discPath, int64(ebr.Ebr_next))
			}

			// Calcular el Porcentaje en BYTES que ocupa la partición en el disco
			partWeight := float64(partition.Part_size) * 100 / float64(mbr.Mbr_size)

			// Escribir la partición en el archivo DOT
			_, err := file.WriteString(`
					<TD COLSPAN="` + fmt.Sprint(ROWSPAN) + `">Extendida<BR/><FONT POINT-SIZE="10" COLOR="red">` + fmt.Sprintf("%.2f", partWeight) + `%` + ` del disco</FONT></TD>`)
			if err != nil {
				return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
			}
		}
	}

	// Calcular el Porcentaje en BYTES que ocupa la ultima partición en el disco
	ultimaParticion := mbr.Mbr_size - (mbr.Mbr_partitions[LastIndex].Part_start + mbr.Mbr_partitions[LastIndex].Part_size)
	partWeight := float64(ultimaParticion) * 100 / float64(mbr.Mbr_size)

	if partWeight > 0.01 {
		// Escribir la partición en el archivo DOT
		_, err = file.WriteString(`
				<TD ROWSPAN="2">Libre<BR/><FONT POINT-SIZE="10" COLOR="red">` + fmt.Sprintf("%.2f", partWeight) + `%` + ` del disco</FONT></TD>`)
		if err != nil {
			return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
		}
	}

	// Escribir en el archivo DOT
	_, err = file.WriteString(`
		    </TR>`)
	if err != nil {
		return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
	}

	if ExtendidaIndex != -1 {
		// Escribir en el archivo DOT
		_, err = file.WriteString(`
			<TR>`)
		if err != nil {
			return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
		}
		// Crear una variable para almacenar los datos leídos
		var ebr structs.EBR

		// Obtiene el primer ebr del archivo.
		ebr.DeserializeEBR(discPath, int64(mbr.Mbr_partitions[ExtendidaIndex].Part_start))

		validarEBR := false
		lastnext := int32(-1)

		// Escribir las particiones logicas en el archivo DOT
		for ebr.Ebr_next != -1 {
			validarEBR = true
			// Calcular el Porcentaje en BYTES que ocupa la partición en el disco
			logicaweight := float64(ebr.Ebr_size) * 100 / float64(mbr.Mbr_size)

			// Calcular el Porcentaje en BYTES que ocupa el ebr en la partición extendida
			ebrWeight := float64(unsafe.Sizeof(ebr)) * 100 / float64(mbr.Mbr_partitions[ExtendidaIndex].Part_size)

			// Escribir la partición en el archivo DOT
			_, err := file.WriteString(`
					<TD>EBR<BR/><FONT POINT-SIZE="5" COLOR="blue">` + fmt.Sprintf("%.2f", ebrWeight) + `%` + ` del disco</FONT></TD>
					<TD>Lógica<BR/><FONT POINT-SIZE="10" COLOR="red">` + fmt.Sprintf("%.2f", logicaweight) + `%` + ` del disco</FONT></TD>`)
			if err != nil {
				return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
			}

			lastnext = ebr.Ebr_next

			// Obtiene el proximo ebr del archivo
			ebr.DeserializeEBR(discPath, int64(ebr.Ebr_next))
		}

		if !validarEBR {
			// Calcular el Porcentaje en BYTES que ocupa el ebr en la partición extendida
			ebrWeight := float64(unsafe.Sizeof(ebr)) * 100 / float64(mbr.Mbr_partitions[ExtendidaIndex].Part_size)

			// Escribir la partición en el archivo DOT
			_, err := file.WriteString(`
					<TD>EBR<BR/><FONT POINT-SIZE="5" COLOR="blue">` + fmt.Sprintf("%.2f", ebrWeight) + `%` + ` del disco</FONT></TD>`)
			if err != nil {
				return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
			}
		}

		if lastnext == -1 {
			lastnext = mbr.Mbr_partitions[ExtendidaIndex].Part_start
		}

		// Calcular el espacio libre en la ultima partición logica
		ultimaParticionLogica := float64((mbr.Mbr_partitions[ExtendidaIndex].Part_size + mbr.Mbr_partitions[ExtendidaIndex].Part_start - lastnext)) * 100 / float64(mbr.Mbr_size)

		if ultimaParticionLogica > 0.01 {
			_, err = file.WriteString(`
			<TD>Libre<BR/><FONT POINT-SIZE="10" COLOR="red">` + fmt.Sprintf("%.2f", ultimaParticionLogica) + `%` + ` del disco</FONT></TD>`)
			if err != nil {
				return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
			}
		}

		// Cerrar la tabla
		_, err = file.WriteString(`
	</TR>`)
		if err != nil {
			return fmt.Errorf("error al escribir el tamaño de la particion en el archivo DOT: %v", err)
		}

	}
	// Cerrar la tabla
	_, err = file.WriteString(`
        </TABLE>
    >];
}`)

	if err != nil {
		return fmt.Errorf("error al cerrar la tabla en el archivo DOT: %v", err)
	}

	return nil
}

// Función para generar el archivo DOT con la estructura del MBR
func GenerateDotFileMBR(mbr structs.MBR, filePath string, discPath string) error {
	// Crear el archivo .dot
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error al crear el archivo DOT: %v", err)
	}
	defer file.Close()

	fitString := ByteToString(mbr.Mbr_fit[:])

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
            <tr><td>mbr_disk_fit</td><td>` + fitString + `</td></tr>
    `)

	if err != nil {
		return fmt.Errorf("error al escribir en el archivo DOT: %v", err)
	}

	// Escribir las particiones en el archivo DOT
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_start == -1 {
			continue
		}

		if partition.Part_type[0] == 'E' {
			// Convertir Part_name a string ignorando los caracteres nulos (0)
			partName := ByteToString(partition.Part_name[:])

			partID := ByteToString(partition.Part_id[:])

			// Escribir la partición en el archivo DOT
			_, err := file.WriteString(`
            <tr><td colspan="2" bgcolor="purple"><font color="white">Particion</font></td></tr>
			<tr><td>part_status</td><td>` + fmt.Sprint(partition.Part_status[0]) + `</td></tr>
            <tr><td>part_type</td><td>` + string(partition.Part_type[0]) + `</td></tr>
            <tr><td>part_fit</td><td>` + string(partition.Part_fit[:]) + `</td></tr>
            <tr><td>part_start</td><td>` + fmt.Sprint(partition.Part_start) + `</td></tr>
            <tr><td>part_size</td><td>` + fmt.Sprint(partition.Part_size) + `</td></tr>
            <tr><td>part_name</td><td>` + partName + `</td></tr>
            <tr><td>part_correlative</td><td>` + fmt.Sprint(partition.Part_correlative) + `</td></tr>
            <tr><td>part_id</td><td>` + partID + `</td></tr>`)

			if err != nil {
				return fmt.Errorf("error al escribir particion primaria en el archivo DOT: %v", err)
			}

			// Crear una variable para almacenar los datos leídos
			var ebr structs.EBR

			// Obtiene el ebr del archivo
			ebr.DeserializeEBR(discPath, int64(partition.Part_start))

			// Escribir las particiones logicas en el archivo DOT
			for ebr.Ebr_next != -1 {
				// Convertir Part_name a string ignorando los caracteres nulos (0)
				EbrPartName := ByteToString(ebr.Ebr_name[:])

				partID := ByteToString(partition.Part_id[:])

				// Escribir la partición en el archivo DOT
				_, err = file.WriteString(`
            <tr><td colspan="2" bgcolor="pink"><font color="white">Particion Logica</font></td></tr>
			<tr><td>part_status</td><td>` + fmt.Sprint(ebr.Ebr_mount[0]) + `</td></tr>
            <tr><td>part_next</td><td>` + fmt.Sprint(ebr.Ebr_next) + `</td></tr>
            <tr><td>part_fit</td><td>` + string(ebr.Ebr_fit[:]) + `</td></tr>
            <tr><td>part_start</td><td>` + fmt.Sprint(ebr.Ebr_start) + `</td></tr>
            <tr><td>part_size</td><td>` + fmt.Sprint(ebr.Ebr_size) + `</td></tr>
            <tr><td>part_name</td><td>` + EbrPartName + `</td></tr>
            <tr><td>part_correlative</td><td>` + fmt.Sprint(partition.Part_correlative) + `</td></tr>
            <tr><td>part_id</td><td>` + partID + `</td></tr>`)

				if err != nil {
					return fmt.Errorf("error al escribir particion logica en el archivo DOT: %v", err)
				}

				// Obtiene el proximo ebr del archivo
				ebr.DeserializeEBR(discPath, int64(ebr.Ebr_next))
			}

		} else if partition.Part_type[0] != 'E' {

			// Convertir Part_name a string ignorando los caracteres nulos (0)
			partName := ByteToString(partition.Part_name[:])

			partID := ByteToString(partition.Part_id[:])

			// Escribir la partición en el archivo DOT
			_, err := file.WriteString(`
            <tr><td colspan="2" bgcolor="purple"><font color="white">Particion</font></td></tr>
			<tr><td>part_status</td><td>` + fmt.Sprint(partition.Part_status[0]) + `</td></tr>
            <tr><td>part_type</td><td>` + string(partition.Part_type[0]) + `</td></tr>
            <tr><td>part_fit</td><td>` + string(partition.Part_fit[:]) + `</td></tr>
            <tr><td>part_start</td><td>` + fmt.Sprint(partition.Part_start) + `</td></tr>
            <tr><td>part_size</td><td>` + fmt.Sprint(partition.Part_size) + `</td></tr>
            <tr><td>part_name</td><td>` + partName + `</td></tr>
            <tr><td>part_correlative</td><td>` + fmt.Sprint(partition.Part_correlative) + `</td></tr>
            <tr><td>part_id</td><td>` + partID + `</td></tr>`)

			if err != nil {
				return fmt.Errorf("error al escribir particion primaria en el archivo DOT: %v", err)
			}
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

func ExecuteDot(dotFile, outputImage string) error {
	// Obtener la extensión del archivo de salida (sin el punto inicial)
	ext := filepath.Ext(outputImage)[1:]

	// Si el formato es jpg, primero generamos un archivo png temporal
	tempImage := outputImage
	if ext == "jpg" {
		tempImage = outputImage[:len(outputImage)-len(ext)] + "png"
	}

	// Ejecutar el comando `dot` con el formato dinámico (basado en la extensión del archivo de salida o png si es jpg)
	cmd := exec.Command("dot", "-T"+filepath.Ext(tempImage)[1:], dotFile, "-o", tempImage)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando dot: %v", err)
	}

	// Si la extensión original es jpg, convertir el archivo png a jpg
	if ext == "jpg" {
		convertCmd := exec.Command("convert", tempImage, outputImage)
		err = convertCmd.Run()
		if err != nil {
			return fmt.Errorf("error al convertir png a jpg: %v", err)
		}

		// Opcionalmente, eliminar el archivo png temporal
		err = os.Remove(tempImage)
		if err != nil {
			return fmt.Errorf("error al eliminar el archivo temporal png: %v", err)
		}
	}

	return nil
}

// Convertir Part_name a string ignorando los caracteres nulos (0)
func ByteToString(partName []byte) string {
	var name []byte
	for _, b := range partName {
		if b == 0 {
			break // Ignorar el resto si encontramos un valor nulo
		}
		name = append(name, b)
	}
	return string(name)
}
