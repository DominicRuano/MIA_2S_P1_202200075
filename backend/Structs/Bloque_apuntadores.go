package structs

import (
	"encoding/binary"
	"fmt"
	"os"
)

/*
Esta estructura guardará la información de los apuntadores indirectos (simples,
dobles y triples)

# Observaciones:

Bloque Simple Indirecto[0-12]: Inodo → Bloque apuntadores → bloque de datos.

Bloque Doble Indirecto[13-14]: Inodo → Bloque de apuntadores → Bloque
de apuntadores → bloque de datos.

Bloque Triple Indirecto[16]: Inodo → Bloque de apuntadores → Bloque
de apuntadores → Bloque de apuntadores → bloque de datos.
*/
type PointerBlock struct {
	B_content [16]int32 // 16 * 4 = 64 bytes
	// Total: 64 bytes
}

// DeserializePointerBlock lee la estructura PointerBlock desde un offset en un archivo binario
func (pb *PointerBlock) DeserializePointerBlock(path string, offset int64) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del PointerBlock
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al buscar el offset para leer PointerBlock: %v", err)
	}

	// Leer directamente la estructura PointerBlock desde el archivo
	if err := binary.Read(file, binary.LittleEndian, pb); err != nil {
		return fmt.Errorf("error al deserializar PointerBlock: %v", err)
	}

	return nil
}

// SerializePointerBlock escribe la estructura PointerBlock en un offset en un archivo binario
func (pb *PointerBlock) SerializePointerBlock(path string, offset int64) error {

	// Abrir el archivo en modo lectura y escritura
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del PointerBlock
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al mover el puntero al offset %d: %v", offset, err)
	}

	// Escribir la estructura PointerBlock en el archivo
	if err := binary.Write(file, binary.LittleEndian, pb); err != nil {
		return fmt.Errorf("error al escribir la estructura PointerBlock en el archivo: %v", err)
	}

	return nil
}
