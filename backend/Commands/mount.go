package commands

import (
	global "Backend/Global"
	structs "Backend/Structs"
	utils "Backend/Utils"
	"fmt"
	"strconv"
	"strings"
)

type Mount_st struct {
	Fdisk_path string
	Fdisk_name string
}

func (mount *Mount_st) Print() {
	fmt.Println("Fdisk_path:", mount.Fdisk_path)
	fmt.Println("Fdisk_name:", mount.Fdisk_name)
}

func Mount(tokens []string) string {

	// Crear un nuevo comando
	Cmd := &Mount_st{}

	// Expresión regular para los parametros
	Regex := `(?i)-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`

	// Parsear los parametros
	tokens = utils.ParseParametros(tokens, Regex)

	// Verificar si hay dos parametros
	if len(tokens) != 2 { // Si no hay cuatro parametros, no hace nada
		return "Error: Comando FDISK requiere 2 parametros (path, name).\n"
	}

	// Iterar sobre los tokens para obtener los parametros
	for _, token := range tokens {
		// Separar el token en partes
		partes := strings.SplitN(token, "=", 2)
		if len(partes) != 2 {
			return fmt.Sprintf("formato de parámetro inválido: %s", token)
		}

		// Switch para manejar los parametros
		switch strings.ToLower(partes[0]) {
		case "-path":
			Cmd.Fdisk_path = strings.ReplaceAll(partes[1], "\"", "")
		case "-name":
			Cmd.Fdisk_name = partes[1]
		default:
			return fmt.Sprintf("Parámetro %s no reconocido\n", partes[0])
		}
	}

	// Verifica los parametros
	if Cmd.Fdisk_path == "" {
		return "Error: Falta el parametro obligatorio -path.\n"
	}
	if Cmd.Fdisk_name == "" {
		return "Error: Falta el parametro obligatorio -name.\n"
	}

	// Crea el objeto mbr
	mbr := &structs.MBR{}

	// Cargar el mbr
	err := mbr.DeserializeMBR(Cmd.Fdisk_path)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Verificar si la particion Existe
	index := mbr.GetIndexByName(Cmd.Fdisk_name)
	if index == -1 {
		return fmt.Sprintf("Error: No se encontró la partición %s\n", Cmd.Fdisk_name)
	}

	// Verificar si la particion ya esta montada
	if mbr.Mbr_partitions[index].Part_status == [1]byte{'1'} {
		return fmt.Sprintf("Error: La partición %s ya está montada\n", Cmd.Fdisk_name)
	}

	// Verificar que la particion no sea extendida
	if mbr.Mbr_partitions[index].Part_type[0] == 'E' {
		return fmt.Sprintf("Error: La partición %s es una partición extendida\n", Cmd.Fdisk_name)
	}

	// Marcar la particion como montada
	mbr.Mbr_partitions[index].Part_status = [1]byte{1}

	// Obtiene la letra de la particion.
	letter, err := global.GetLetter(Cmd.Fdisk_path)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Obtiene el numero de la particion.
	NumPartition := global.GetNumPartition(Cmd.Fdisk_path)

	ID_Part := global.Carnet + strconv.Itoa(NumPartition) + letter
	mbr.Mbr_partitions[index].Part_id = [4]byte{ID_Part[0], ID_Part[1], ID_Part[2], ID_Part[3]}

	// Marca el Correlativo de la particion
	mbr.Mbr_partitions[index].Part_correlative = int32(index + 1)

	// Guarda los cambios
	err = mbr.SerializeMBR(Cmd.Fdisk_path)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Guarda la particion montada en memoria
	global.MountedPartitions[ID_Part] = Cmd.Fdisk_path

	return fmt.Sprintf("Partición %s con Id: %s montada con exito.\n", Cmd.Fdisk_name, ID_Part)
}
