package structs

import (
	"encoding/binary"
	"fmt"
	"os"
)

/*
Esta estructura guardará la información sobre el nombre de de los archivos que
contiene y a que Inodo apuntan.

En cada inodo de carpeta, en el primer apuntador directo, en los primeros dos
registros se guardará el nombre de la carpeta y su padre
*/
type FolderBlock struct {
	B_content [4]FolderContent // 4 * 16 = 64 bytes
	// Total: 64 bytes
}

/*
 */
type FolderContent struct {
	B_name  [12]byte // Nombre de la carpeta o archivo
	B_inodo int32    // Apuntador hacia un inodo asociado al archivo o carpeta
}

// DeserializeFolderBlock lee la estructura FolderBlock desde un offset en un archivo binario
func (fb *FolderBlock) DeserializeFolderBlock(path string, offset int64) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del FolderBlock
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al buscar el offset para leer FolderBlock: %v", err)
	}

	// Leer directamente la estructura FolderBlock desde el archivo
	if err := binary.Read(file, binary.LittleEndian, fb); err != nil {
		return fmt.Errorf("error al deserializar FolderBlock: %v", err)
	}

	return nil
}

// SerializeSB escribe la estructura FolderBlock en un offset en un archivo binario
func (fb *FolderBlock) SerializeFolderBlock(path string, offset int64) error {

	// Abrir el archivo en modo lectura y escritura
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Mover el puntero al inicio del FolderBlock
	if _, err := file.Seek(offset, 0); err != nil {
		return fmt.Errorf("error al mover el puntero al offset %d: %v", offset, err)
	}

	// Escribir la estructura FolderBlock en el archivo
	if err := binary.Write(file, binary.LittleEndian, fb); err != nil {
		return fmt.Errorf("error al escribir la estructura FolderBlock en el archivo: %v", err)
	}

	return nil
}

func (fb *FolderBlock) Print() {
	fmt.Printf("FolderBlock:\n")
	for i, content := range fb.B_content {
		fmt.Printf("Content %d:\n", i)
		fmt.Printf("B_name: %s\n", string(content.B_name[:]))
		fmt.Printf("B_inodo: %d\n", content.B_inodo)
	}
}
