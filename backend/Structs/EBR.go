package structs

import (
	"encoding/binary"
	"fmt"
	"os"
)

type EBR struct {
	Ebr_mount [1]byte  // 1 byte
	Ebr_fit   [1]byte  // 1 byte
	Ebr_start int32    // 4 bytes
	Ebr_size  int32    // 4 bytes
	Ebr_next  int32    // 4 bytes
	Ebr_name  [16]byte // 16 bytes
	// Total = 30 bytes
}

// DeserializeEBR lee la estructura EBR desde un offset en un archivo binario
func (ebr *EBR) DeserializeEBR(path string, offset int64) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del EBR
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al buscar el offset para leer EBR: %v", err)
	}

	// Leer directamente la estructura EBR desde el archivo
	if err := binary.Read(file, binary.LittleEndian, ebr); err != nil {
		return fmt.Errorf("error al deserializar EBR: %v", err)
	}

	return nil
}

// SerializeEBR escribe la estructura EBR en un offset en un archivo binario
func (ebr *EBR) SerializeEBR(path string, offset int64) error {

	// Abrir el archivo en modo lectura y escritura
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del EBR
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al mover el puntero al offset %d: %v", offset, err)
	}

	// Escribir la estructura EBR en el archivo
	if err := binary.Write(file, binary.LittleEndian, ebr); err != nil {
		return fmt.Errorf("error al escribir la estructura EBR en el archivo: %v", err)
	}

	return nil
}
