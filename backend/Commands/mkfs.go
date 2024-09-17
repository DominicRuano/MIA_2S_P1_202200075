package commands

import (
	global "Backend/Global"
	structs "Backend/Structs"
	utils "Backend/Utils"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

type mkfs struct {
	Id   string
	Type string
	Path string
}

func (m *mkfs) Print() {
	fmt.Printf("Id: %s\n", m.Id)
	fmt.Printf("Type: %s\n", m.Type)
}

func Mkfs(tokens []string) string {

	// Crea un nuevo objeto mkfs
	Cmd := &mkfs{}

	// Expressión regular para validar los parametros
	Regex := `(?i)-id="[^"]+"|-id=[^\s]+|-type="[^"]+"|-type=[^\s]+`

	// Parsea los parametros
	tokens = utils.ParseParametros(tokens, Regex)

	if len(tokens) > 2 || len(tokens) < 1 { // Si no un parametro, no hace nada
		return "Error: Comando MKFS requiere minimo 1 parametros (id)  y maximo 2 (id, type).\n"
	}

	for _, token := range tokens { // Itera sobre los tokens para obtener los parametros
		partes := strings.SplitN(token, "=", 2) // Separa el token en partes
		if len(partes) != 2 {
			return fmt.Sprintf("formato de parámetro inválido: %s", token)
		}

		switch strings.ToLower(partes[0]) { // Switch para manejar los parametros
		case "-id":
			if PathOrigen, exists := global.MountedPartitions[partes[1]]; exists {
				Cmd.Id = partes[1]
				Cmd.Path = PathOrigen
			} else {
				return fmt.Sprintf("Error: La particion [%s] no esta montada.\n", partes[1])
			}
		case "-type":
			Cmd.Type = partes[1]
		default:
			return fmt.Sprintf("MKFS: Parametro %s no reconocido.\n", partes[0])
		}
	}

	// Verifica los parametros
	if Cmd.Id == "" {
		return "Error: Parametro -id requerido.\n"
	}
	if Cmd.Type == "" { // Si no se especifica el tipo, se asume full
		Cmd.Type = "full"
	}

	// Verifica el tipo de formateo
	switch strings.ToLower(Cmd.Type) {
	case "full":
		return FullFformat(Cmd)
	case "fast":
		return fmt.Sprintf("mkfs: Particion %s (aun no implementado el mkfs para fast).\n", Cmd.Id)
	default:
		return fmt.Sprintf("Error: Tipo de formateo %s no reconocido.\n", Cmd.Type)
	}

}

func FullFformat(Cmd *mkfs) string {

	// Deserializa el mbr
	mbr := &structs.MBR{}

	// Lee el archivo binario
	err := mbr.DeserializeMBR(Cmd.Path)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Busca la particion
	Particion := mbr.GetPartitionById(Cmd.Id)
	if Particion == nil {
		return fmt.Sprintf("Error: Particion %s no encontrada.\n", Cmd.Id)
	}

	// Borra la particion y la formatea
	erro := BorrarPart(Cmd.Path, Particion)
	if erro != "" {
		return erro
	}

	// Calcula el n
	n := utils.CalculateN(Particion)
	fmt.Printf("N vale: %d \n", n)

	// Crea el superbloque
	SuperBloque := createSuperBlock(Particion, n)

	SuperBloque.Print()

	// Crear los bitmaps
	err = SuperBloque.CreateBitMaps(Cmd.Path)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Crear archivo users.txt
	err = SuperBloque.CreateUsersFile(Cmd.Path)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Serializar el superbloque
	err = SuperBloque.SerializeSB(Cmd.Path, int64(Particion.Part_start))
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	return fmt.Sprintf("mkfs: Particion %s Formateada con exito (aun no implementado).\n", Cmd.Id)
}

// createSuperBlock crea un nuevo superbloque
func createSuperBlock(partition *structs.Partition, n int32) *structs.SuperBloque {
	// Calcular punteros de las estructuras
	// Bitmaps
	bm_inode_start := partition.Part_start + int32(binary.Size(structs.SuperBloque{}))
	bm_block_start := bm_inode_start + n // n indica la cantidad de inodos, solo la cantidad para ser representada en un bitmap
	// Inodos
	inode_start := bm_block_start + (3 * n) // 3*n indica la cantidad de bloques, se multiplica por 3 porque se tienen 3 tipos de bloques
	// Bloques
	block_start := inode_start + (int32(binary.Size(structs.Inodo{})) * n) // n indica la cantidad de inodos, solo que aquí indica la cantidad de estructuras Inode

	// Crear un nuevo superbloque
	superBlock := &structs.SuperBloque{
		SB_filesystem_type:   2,
		SB_inodes_count:      0,
		SB_blocks_count:      0,
		SB_free_inodes_count: int32(n),
		SB_free_blocks_count: int32(n * 3),
		SB_mtime:             float64(time.Now().Unix()),
		SB_umtime:            float64(time.Now().Unix()),
		SB_mnt_count:         1,
		SB_magic:             0xEF53,
		SB_inode_size:        int16(binary.Size(structs.Inodo{})),
		SB_block_size:        int16(binary.Size(structs.FileBlock{})),
		SB_firs_ino:          inode_start,
		SB_first_blo:         block_start,
		SB_bm_inode_start:    bm_inode_start,
		SB_bm_block_start:    bm_block_start,
		SB_inode_start:       inode_start,
		SB_block_start:       block_start,
	}
	return superBlock
}

func BorrarPart(path string, Particion *structs.Partition) string {
	// Abrimos el archivo
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}
	defer file.Close()

	// Mover el puntero del archivo al inicio de la particion
	_, err = file.Seek(int64(Particion.Part_start), 0)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	// Crear un buffer de n '0'
	buffer := make([]byte, Particion.Part_size)
	for i := range buffer {
		buffer[i] = 0
	}

	// Escribir el buffer en el archivo
	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", err)
	}

	return ""
}
