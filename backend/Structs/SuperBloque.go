package structs

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type SuperBloque struct {
	SB_filesystem_type   int32   // 4 bytes
	SB_inodes_count      int32   // 4 bytes
	SB_blocks_count      int32   // 4 bytes
	SB_free_blocks_count int32   // 4 bytes
	SB_free_inodes_count int32   // 4 bytes
	SB_mtime             float64 // 8 bytes
	SB_umtime            float64 // 8 bytes
	SB_mnt_count         int16   // 2 bytes
	SB_magic             uint16  // 2 bytes
	SB_inode_size        int16   // 2 bytes
	SB_block_size        int16   // 2 bytes
	SB_firs_ino          int32   // 4 bytes
	SB_first_blo         int32   // 4 bytes
	SB_bm_inode_start    int32   // 4 bytes
	SB_bm_block_start    int32   // 4 bytes
	SB_inode_start       int32   // 4 bytes
	SB_block_start       int32   // 4 bytes
	// Totla = 68 bytes + 4 bytes de padding = 72 bytes
}

// DeserializeSB lee la estructura SuperBloque desde un offset en un archivo binario
func (sb *SuperBloque) DeserializeSB(path string, offset int64) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del SuperBloque
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al buscar el offset para leer SuperBloque: %v", err)
	}

	// Leer directamente la estructura SuperBloque desde el archivo
	if err := binary.Read(file, binary.LittleEndian, sb); err != nil {
		return fmt.Errorf("error al deserializar SuperBloque: %v", err)
	}

	return nil
}

// SerializeSB escribe la estructura SuperBloque en un offset en un archivo binario
func (sb *SuperBloque) SerializeSB(path string, offset int64) error {

	// Abrir el archivo en modo lectura y escritura
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del SuperBloque
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al mover el puntero al offset %d: %v", offset, err)
	}

	// Escribir la estructura SuperBloque en el archivo
	if err := binary.Write(file, binary.LittleEndian, sb); err != nil {
		return fmt.Errorf("error al escribir la estructura SuperBloque en el archivo: %v", err)
	}

	return nil
}

// Crear users.txt
func (sb *SuperBloque) CreateUsersFile(path string) error {
	// ----------- Creamos / -----------
	// Creamos el inodo raíz
	rootInode := &Inodo{
		I_uid:   1,
		I_gid:   1,
		I_size:  0,
		I_atime: float64(time.Now().Unix()),
		I_ctime: float64(time.Now().Unix()),
		I_mtime: float64(time.Now().Unix()),
		I_block: [15]int32{sb.SB_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'0'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Serializar el inodo raíz
	err := rootInode.SerializeInodo(path, int64(sb.SB_firs_ino))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.SB_inodes_count++
	sb.SB_free_inodes_count--
	sb.SB_firs_ino += int32(sb.SB_inode_size)

	// Creamos el bloque del Inodo Raíz
	rootBlock := &FolderBlock{
		B_content: [4]FolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.SerializeFolderBlock(path, int64(sb.SB_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.SB_blocks_count++
	sb.SB_free_blocks_count--
	sb.SB_first_blo += int32(sb.SB_block_size)

	// ----------- Creamos /users.txt -----------
	usersText := "1,G,root\n1,U,root,123\n"

	// Deserializar el inodo raíz
	err = rootInode.DeserializeInodo(path, int64(sb.SB_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Actualizamos el inodo raíz
	rootInode.I_atime = float64(time.Now().Unix())

	// Serializar el inodo raíz
	err = rootInode.SerializeInodo(path, int64(sb.SB_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Deserializar el bloque de carpeta raíz
	err = rootBlock.DeserializeFolderBlock(path, int64(sb.SB_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Actualizamos el bloque de carpeta raíz
	rootBlock.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.SB_inodes_count}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.SerializeFolderBlock(path, int64(sb.SB_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Creamos el inodo users.txt
	usersInode := &Inodo{
		I_uid:   1,
		I_gid:   1,
		I_size:  int32(len(usersText)),
		I_atime: float64(time.Now().Unix()),
		I_ctime: float64(time.Now().Unix()),
		I_mtime: float64(time.Now().Unix()),
		I_block: [15]int32{sb.SB_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'1'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Serializar el inodo users.txt
	err = usersInode.SerializeInodo(path, int64(sb.SB_firs_ino))
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.SB_inodes_count++
	sb.SB_free_inodes_count--
	sb.SB_firs_ino += int32(sb.SB_inode_size)

	// Creamos el bloque de users.txt
	usersBlock := &FileBlock{
		B_content: [64]byte{},
	}
	// Copiamos el texto de usuarios en el bloque
	copy(usersBlock.B_content[:], usersText)

	// Serializar el bloque de users.txt
	err = usersBlock.SerializeFileBlock(path, int64(sb.SB_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.SB_blocks_count++
	sb.SB_free_blocks_count--
	sb.SB_first_blo += int32(sb.SB_block_size)

	return nil
}

func (sb *SuperBloque) Print() {

	// Convertir los tiempos a formato legible
	mtime := time.Unix(int64(sb.SB_mtime), 0)
	umtime := time.Unix(int64(sb.SB_umtime), 0)

	// Imprimir el superbloque
	fmt.Printf("SuperBloque:\n")
	fmt.Printf("SB_filesystem_type: %d\n", sb.SB_filesystem_type)
	fmt.Printf("SB_inodes_count: %d\n", sb.SB_inodes_count)
	fmt.Printf("SB_blocks_count: %d\n", sb.SB_blocks_count)
	fmt.Printf("SB_free_blocks_count: %d\n", sb.SB_free_blocks_count)
	fmt.Printf("SB_free_inodes_count: %d\n", sb.SB_free_inodes_count)
	fmt.Printf("SB_mtime: %s\n", mtime)
	fmt.Printf("SB_umtime: %s\n", umtime)
	fmt.Printf("SB_mnt_count: %d\n", sb.SB_mnt_count)
	fmt.Printf("SB_magic: %d\n", sb.SB_magic)
	fmt.Printf("SB_inode_size: %d\n", sb.SB_inode_size)
	fmt.Printf("SB_block_size: %d\n", sb.SB_block_size)
	fmt.Printf("SB_firs_ino: %d\n", sb.SB_firs_ino)
	fmt.Printf("SB_first_blo: %d\n", sb.SB_first_blo)
	fmt.Printf("SB_bm_inode_start: %d\n", sb.SB_bm_inode_start)
	fmt.Printf("SB_bm_block_start: %d\n", sb.SB_bm_block_start)
	fmt.Printf("SB_inode_start: %d\n", sb.SB_inode_start)
	fmt.Printf("SB_block_start: %d\n", sb.SB_block_start)

}
