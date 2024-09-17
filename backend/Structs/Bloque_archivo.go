package structs

import (
	"encoding/binary"
	"fmt"
	"os"
)

/*
Esta estructura guardará la información sobre contenido de un archivo.

# Observaciones:

1. El bloque archivo guarda el contenido del archivo (64 caracteres).
*/
type FileBlock struct {
	B_content [64]byte // 64 bytes
	// Total: 64 bytes
}

// DeserializeFileBlock lee la estructura FileBlock desde un offset en un archivo binario
func (fb *FileBlock) DeserializeFileBlock(path string, offset int64) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del FileBlock
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al buscar el offset para leer FileBlock: %v", err)
	}

	// Leer directamente la estructura FileBlock desde el archivo
	if err := binary.Read(file, binary.LittleEndian, fb); err != nil {
		return fmt.Errorf("error al deserializar FileBlock: %v", err)
	}

	return nil
}

// SerializeFileBlock escribe la estructura FileBlock en un offset en un archivo binario
func (fb *FileBlock) SerializeFileBlock(path string, offset int64) error {

	// Abrir el archivo en modo lectura y escritura
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del FileBlock
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al mover el puntero al offset %d: %v", offset, err)
	}

	// Escribir la estructura FileBlock en el archivo
	if err := binary.Write(file, binary.LittleEndian, fb); err != nil {
		return fmt.Errorf("error al escribir la estructura FileBlock en el archivo: %v", err)
	}

	return nil
}

func (fb *FileBlock) Print() {
	fmt.Printf("FileBlock:\n")
	fmt.Printf("B_content: %s\n", fb.B_content)
}
