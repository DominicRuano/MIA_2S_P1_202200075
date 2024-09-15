package commands

import (
	structs "Backend/Structs"
	utils "Backend/Utils"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type Fdisk struct {
	Fdisk_size int
	Fdisk_unit string
	Fdisk_path string
	Fdisk_type string
	Fdisk_fit  string
	Fdisk_name string
}

func (fdisk *Fdisk) Print() {
	fmt.Println("Fdisk_size:", fdisk.Fdisk_size)
	fmt.Println("Fdisk_unit:", fdisk.Fdisk_unit)
	fmt.Println("Fdisk_path:", fdisk.Fdisk_path)
	fmt.Println("Fdisk_type:", fdisk.Fdisk_type)
	fmt.Println("Fdisk_fit:", fdisk.Fdisk_fit)
	fmt.Println("Fdisk_name:", fdisk.Fdisk_name)
}

func FDisk(tokens []string) string {
	Command := Fdisk{}

	Regex := `(?i)-size=\d+|-unit=[^\s]|-fit=[^\s]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`

	tokens = utils.ParseParametros(tokens, Regex)

	if len(tokens) > 6 || len(tokens) < 3 { // Si no hay cuatro parametros, no hace nada
		return "Error: Comando FDISK requiere minimo 3 parametros (path, size, name)  y maximo 6 (path, size, name, fit, type, unit ).\n"
	}

	for _, token := range tokens { // Itera sobre los tokens para obtener los parametros
		partes := strings.SplitN(token, "=", 2) // Separa el token en partes
		if len(partes) != 2 {
			return fmt.Sprintf("formato de parámetro inválido: %s", token)
		}

		switch strings.ToLower(partes[0]) { // Switch para manejar los parametros
		case "-size":
			Command.Fdisk_size, _ = strconv.Atoi(partes[1])

			if Command.Fdisk_size <= 0 {
				return "Error: El tamaño debe ser mayor a cero.\n"
			}
		case "-unit":
			Command.Fdisk_unit = strings.ToUpper(partes[1])

			if strings.ToUpper(Command.Fdisk_unit) != "K" && strings.ToUpper(Command.Fdisk_unit) != "M" && strings.ToUpper(Command.Fdisk_unit) != "B" {
				return "Error: Unidad de medida no reconocida, debe ser K o M.\n"
			}
		case "-fit":
			Command.Fdisk_fit = strings.ToUpper(partes[1])

			if Command.Fdisk_fit != "BF" && Command.Fdisk_fit != "FF" && Command.Fdisk_fit != "WF" {
				return "Error: Tipo de ajuste no reconocido, el ajuste debe ser BF, FF o WF.\n"
			}
		case "-path":
			Command.Fdisk_path = strings.ReplaceAll(partes[1], "\"", "")
		case "-name":
			Command.Fdisk_name = strings.ReplaceAll(partes[1], "\"", "")
		case "-type":
			Command.Fdisk_type = strings.ToUpper(partes[1])
		default:
			return fmt.Sprintf("Error: Parametro [" + partes[0] + "] no reconocido.\n")
		}

	}

	Command.Print()
	println("")

	// Manejo de erorres en los parametros
	if Command.Fdisk_size == 0 {
		return "Error: Tamaño no especificado.\n"
	}

	if Command.Fdisk_path == "" {
		return "Error: Path no especificado.\n"
	}

	if Command.Fdisk_name == "" {
		return "Error: Nombre no especificado.\n"
	}

	if Command.Fdisk_unit == "" {
		Command.Fdisk_unit = "K"
	}

	if Command.Fdisk_type == "" {
		Command.Fdisk_type = "P"
	}

	if Command.Fdisk_fit == "" {
		Command.Fdisk_fit = "WF"
	}

	// Se crea una estructura MBR para manejar las particiones
	var mbr = &structs.MBR{}
	mbr.DeserializeMBR(Command.Fdisk_path) // Se carga el MBR del disco

	//Dependiendo de que tipo de particion sea, se ejecuta una funcion diferente
	if strings.ToUpper(Command.Fdisk_type) == "P" {
		return Command.CrearParticionPrimaria(mbr)
	} else if strings.ToUpper(Command.Fdisk_type) == "E" {
		return Command.CrearParticionExtendida(mbr)
	} else if strings.ToUpper(Command.Fdisk_type) == "L" {
		return Command.CrearParticionLogica(mbr)
	} else {
		return "Error: Tipo de particion no reconocido, debe ser P, E o L.\n"
	}
}

func (fdisk *Fdisk) CrearParticionPrimaria(mbr *structs.MBR) string {

	//Verificar si ya existe una particion con ese nombre
	IfNameExist := mbr.Verifyname(fdisk.Fdisk_name)
	if IfNameExist == -1 {
		return "Error: Ya existe una particion con ese nombre.\n"
	}

	//Obtener el indice de la particion
	PartIndex, err := mbr.GetPartitionIndex(fdisk.Fdisk_fit)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	start := mbr.CalcularStart(PartIndex)
	size := utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit)

	//Verificar si la particion cabe en el disco
	if !mbr.CabeParticion(utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit), start) {
		return fmt.Sprintln("Error: La particion no cabe en el disco (", utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit)+start, "B /", mbr.Mbr_size, "B ).")
	}

	//Crear particion primaria
	mbr.Mbr_partitions[PartIndex].Part_status = [1]byte{0}
	mbr.Mbr_partitions[PartIndex].Part_type = [1]byte{fdisk.Fdisk_type[0]}
	mbr.Mbr_partitions[PartIndex].Part_fit = [1]byte{fdisk.Fdisk_fit[0]}
	mbr.Mbr_partitions[PartIndex].Part_start = int32(start)
	mbr.Mbr_partitions[PartIndex].Part_size = int32(size)
	copy(mbr.Mbr_partitions[PartIndex].Part_name[:], fdisk.Fdisk_name)

	mbr.SerializeMBR(fdisk.Fdisk_path)

	return "Particion primaria creada.\n"
}

func (fdisk *Fdisk) CrearParticionExtendida(mbr *structs.MBR) string {

	// Verificar si ya existe una particion extendida
	if mbr.ExistExtended() {
		return "Error: Ya existe una particion extendida.\n"
	}

	//Verificar si ya existe una particion con ese nombre
	IfNameExist := mbr.Verifyname(fdisk.Fdisk_name)
	if IfNameExist == -1 {
		return "Error: Ya existe una particion con ese nombre.\n"
	}

	//Obtener el indice de la particion
	PartIndex, err := mbr.GetPartitionIndex(fdisk.Fdisk_fit)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	start := mbr.CalcularStart(PartIndex)
	size := utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit)

	//Verificar si la particion cabe en el disco
	if !mbr.CabeParticion(utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit), start) {
		return fmt.Sprintln("Error: La particion no cabe en el disco (", utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit)+start, "B /", mbr.Mbr_size, "B ).")
	}

	//Crear particion primaria
	mbr.Mbr_partitions[PartIndex].Part_status = [1]byte{0}
	mbr.Mbr_partitions[PartIndex].Part_type = [1]byte{fdisk.Fdisk_type[0]}
	mbr.Mbr_partitions[PartIndex].Part_fit = [1]byte{fdisk.Fdisk_fit[0]}
	mbr.Mbr_partitions[PartIndex].Part_start = int32(start)
	mbr.Mbr_partitions[PartIndex].Part_size = int32(size)
	copy(mbr.Mbr_partitions[PartIndex].Part_name[:], fdisk.Fdisk_name)

	mbr.SerializeMBR(fdisk.Fdisk_path)

	ebr := &structs.EBR{}

	//Crear el EBR
	ebr.Ebr_mount = [1]byte{'N'}
	ebr.Ebr_fit = [1]byte{0}
	ebr.Ebr_start = int32(start)
	ebr.Ebr_size = int32(-1)
	ebr.Ebr_next = -1
	ebr.Ebr_name = [16]byte{'N'}

	err = ebr.SerializeEBR(fdisk.Fdisk_path, int64(start))
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	return "Particion extendida creada.\n"
}

func (fdisk *Fdisk) CrearParticionLogica(mbr *structs.MBR) string {

	// Verificar si ya existe una particion extendida
	if !mbr.ExistExtended() {
		return "Error: No existe una particion extendida, no es posible crear una logica.\n"
	}

	//Verificar si ya existe una particion P o E con ese nombre
	IfNameExist := mbr.Verifyname(fdisk.Fdisk_name)
	if IfNameExist == -1 {
		return "Error: Ya existe una particion P o E con ese nombre.\n"
	}

	//Obtener el indice de la particion Extendida
	PartIndex := mbr.GetextendedPartitionIndex()
	if PartIndex == -1 {
		return "Error: No se pudo encontrar una particion extendida.\n"
	}

	//Crear el EBR
	ebr := &structs.EBR{}

	//deserialize EBR
	err := ebr.DeserializeEBR(fdisk.Fdisk_path, int64(mbr.Mbr_partitions[PartIndex].Part_start))
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	for ebr.Ebr_next != -1 {
		// Verificar si ya existe una particion con ese nombre
		if string(bytes.Trim(ebr.Ebr_name[:], "\x00")) == fdisk.Fdisk_name {
			return "Error: Ya existe una particion Logica con ese nombre.\n"
		}

		err = ebr.DeserializeEBR(fdisk.Fdisk_path, int64(ebr.Ebr_next))
		if err != nil {
			return fmt.Sprintf("Error: %s\n", err)
		}
	}

	//Verificar si hay espacio suficiente para crear la particion logica
	if int32(ebr.Ebr_start+ebr.Ebr_size)+int32(utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit))+int32(unsafe.Sizeof(structs.EBR{})) >= mbr.Mbr_partitions[PartIndex].Part_start+mbr.Mbr_partitions[PartIndex].Part_size {
		return "Error: No hay espacio suficiente para crear la particion logica.\n"
	}

	//Modificar el EBR anterior
	ebr.Ebr_mount = [1]byte{0}
	ebr.Ebr_fit = [1]byte{fdisk.Fdisk_fit[0]}
	ebr.Ebr_size = int32(utils.CalcularTamaño(fdisk.Fdisk_size, fdisk.Fdisk_unit))
	ebr.Ebr_next = int32(ebr.Ebr_start + ebr.Ebr_size)
	copy(ebr.Ebr_name[:], fdisk.Fdisk_name)

	err = ebr.SerializeEBR(fdisk.Fdisk_path, int64(ebr.Ebr_start))
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	//Crear el nuevo EBR
	New_ebr := &structs.EBR{}

	//Inicializa el nuevo EBR
	New_ebr.Ebr_mount = [1]byte{'N'}
	New_ebr.Ebr_fit = [1]byte{0}
	New_ebr.Ebr_start = int32(ebr.Ebr_next)
	New_ebr.Ebr_size = int32(-1)
	New_ebr.Ebr_next = -1
	New_ebr.Ebr_name = [16]byte{'N'}

	err = New_ebr.SerializeEBR(fdisk.Fdisk_path, int64(ebr.Ebr_next))
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	return "Particion logica Creada con exito.\n"
}
