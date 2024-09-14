package structs

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type MBR struct {
	Mbr_size       int32        // 4 bytes
	Mbr_date       float64      // 8 bytes
	Mbr_signature  int32        // 4 bytes
	Mbr_fit        [2]byte      // 2 byte
	Mbr_partitions [4]Partition // 4 * (35) = 140 bytes
	// Total = 158 bytes
}

// DeserializeMBR lee la estructura MBR desde el inicio de un archivo binario
func (mbr *MBR) DeserializeMBR(path string) error {
	// Abrir el archivo en modo lectura
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Leer directamente la estructura MBR desde el archivo
	err = binary.Read(file, binary.LittleEndian, mbr)
	if err != nil {
		return fmt.Errorf("error al deserializar el MBR: %v", err)
	}

	return nil
}

// SerializeMBR escribe la estructura MBR al inicio de un archivo binario
func (mbr *MBR) SerializeMBR(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

func (mbr *MBR) Print() {
	// Convertir Mbr_creation_date a time.Time
	creationTime := time.Unix(int64(mbr.Mbr_date), 0)

	// Convertir Mbr_disk_fit a char
	diskFit := rune(mbr.Mbr_fit[0])

	fmt.Printf("MBR Size: %d\n", mbr.Mbr_size)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_signature)
	fmt.Printf("Disk Fit: %c\n", diskFit)
}

func (mbr *MBR) PrintPartitions() {
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.Part_status[0])
		partType := rune(partition.Part_type[0])
		partFit := rune(partition.Part_fit[0])

		// Convertir Part_name a string
		partName := string(partition.Part_name[:])

		fmt.Printf("Partition %d:\n", i+1)
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.Part_start)
		fmt.Printf("  Size: %d\n", partition.Part_size)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.Part_correlative)
		fmt.Printf("  ID: %d\n", partition.Part_id)
	}
}
