package structs

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Inodo struct {
	I_uid   int32     // 4 bytes
	I_gid   int32     // 4 bytes
	I_size  int32     // 4 bytes
	I_atime float64   // 8 bytes
	I_ctime float64   // 8 bytes
	I_mtime float64   // 8 bytes
	I_block [15]int32 // 60 bytes
	I_type  [1]byte   // 1 byte
	I_perm  [3]byte   // 3 bytes
	// Total = 100 bytes + 4 bytes de padding = 104 bytes
}

// DeserializeInodo lee la estructura Inodo desde un offset en un archivo binario
func (in *Inodo) DeserializeInodo(path string, offset int64) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del SuperBloque
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al buscar el offset para leer Inodo: %v", err)
	}

	// Leer directamente la estructura SuperBloque desde el archivo
	if err := binary.Read(file, binary.LittleEndian, in); err != nil {
		return fmt.Errorf("error al deserializar Inodo: %v", err)
	}

	return nil
}

// SerializeInodo escribe la estructura Inodo en un offset en un archivo binario
func (in *Inodo) SerializeInodo(path string, offset int64) error {

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
	if err := binary.Write(file, binary.LittleEndian, in); err != nil {
		return fmt.Errorf("error al escribir la estructura Inodo en el archivo: %v", err)
	}

	return nil
}

func (in *Inodo) Print() {
	fmt.Printf("Inodo:\n")
	fmt.Printf("I_uid: %d\n", in.I_uid)
	fmt.Printf("I_gid: %d\n", in.I_gid)
	fmt.Printf("I_size: %d\n", in.I_size)
	fmt.Printf("I_atime: %f\n", in.I_atime)
	fmt.Printf("I_ctime: %f\n", in.I_ctime)
	fmt.Printf("I_mtime: %f\n", in.I_mtime)
	fmt.Printf("I_block: %v\n", in.I_block)
	fmt.Printf("I_type: %v\n", in.I_type)
	fmt.Printf("I_perm: %v\n", in.I_perm)
}
