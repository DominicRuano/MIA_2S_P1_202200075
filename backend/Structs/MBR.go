package structs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
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

	fmt.Printf("MBR Size: %d\n", mbr.Mbr_size)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_signature)
	fmt.Printf("Disk Fit: %c%c\n", mbr.Mbr_fit[0], mbr.Mbr_fit[0])
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

/*
GetPartitionIndex devuelve el índice de la partición que cumple con el tipo de ajuste especificado

BF: Best Fit
FF: First Fit
WF: Worst Fit

Si no se encuentra una partición disponible, devuelve -1 y un error
*/
func (mbr *MBR) GetPartitionIndex(Type string) (int, error) {
	if mbr.PartitionAvailable() == -1 {
		return -1, fmt.Errorf("no hay particiones disponibles en el disco")
	}

	if strings.ToUpper(Type) == "BF" {
		return mbr.getBestFit()
	} else if strings.ToUpper(Type) == "FF" {
		return mbr.getFirstFit()
	} else if strings.ToUpper(Type) == "WF" {
		return mbr.getWorstFit()
	}
	return -1, fmt.Errorf("getPartitionIndex, no detecto un tipo de particion correcto")
}

/*
PartitionAvailable devuelve 1 si hay una partición disponible, -1 en caso contrario
*/
func (mbr *MBR) PartitionAvailable() int {
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_start == -1 && partition.Part_size == -1 {
			return 1
		}
	}
	return -1
}

/*
getBestFit devuelve el índice de la partición con el mejor ajuste

# El mejor ajuste es la partición con el tamaño más pequeño que sea mayor o igual al tamaño de la partición a crear

Si no se encuentra una partición disponible, devuelve -1 y un error
*/
func (mbr *MBR) getBestFit() (int, error) {

	for i, partition := range mbr.Mbr_partitions {
		if partition.Part_start == -1 {
			return i, nil
		}
	}

	return -1, fmt.Errorf("no se pudo encontrar una particion BF ")
}

/*
getFirstFit devuelve el índice de la primera partición disponible

Si no se encuentra una partición disponible, devuelve -1 y un error
*/
func (mbr *MBR) getFirstFit() (int, error) {
	for i, partition := range mbr.Mbr_partitions {
		if partition.Part_start == -1 {
			return i, nil
		}
	}

	return -1, fmt.Errorf("error: no se pudo encontrar una particion FF ")
}

/*
getWorstFit devuelve el índice de la partición con el peor ajuste

# El peor ajuste es la partición con el tamaño más grande que sea mayor o igual al tamaño de la partición a crear

Si no se encuentra una partición disponible, devuelve -1 y un error
*/
func (mbr *MBR) getWorstFit() (int, error) {
	for i, partition := range mbr.Mbr_partitions {
		if partition.Part_start == -1 {
			return i, nil
		}
	}

	return -1, fmt.Errorf("error: no se pudo encontrar una particion WF ")
}

/*
Funcion imcompleta, solo suma los valores de las particiones consecutivas.
*/
func (mbr *MBR) CalcularStart(index int) int {
	var start int = int(binary.Size(mbr))

	for i := 0; i < index; i++ {
		if mbr.Mbr_partitions[i].Part_start != -1 {
			start += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	return start
}

func (mbr *MBR) Verifyname(name string) int {
	for _, partition := range mbr.Mbr_partitions {
		// Convertimos el nombre de la partición quitando cualquier carácter nulo
		partitionName := string(bytes.Trim(partition.Part_name[:], "\x00"))

		if partitionName == name {
			return -1 // Existe
		}
	}
	return 1 // No existe
}

/*
Retorna true si existe una partición extendida en el disco
*/
func (mbr *MBR) ExistExtended() bool {
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] == 'E' {
			return true
		}
	}
	return false
}

/*
CabeParticion verifica si una partición cabe en el disco, solo implementada para FF.
*/
func (mbr *MBR) CabeParticion(Size int, Start int) bool {
	// Calcular el final de la partición
	end := Start + Size

	// Verificar si la partición cabe en el disco
	return end <= int(mbr.Mbr_size)
}

func (mbr *MBR) GetextendedPartitionIndex() int {
	for i, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] == 'E' {
			return i
		}
	}
	return -1
}
