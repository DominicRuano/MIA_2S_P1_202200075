package structs

type MBR struct {
	Mbr_size       int32        // 4 bytes
	Mbr_date       float64      // 8 bytes
	Mbr_signature  int32        // 4 bytes
	Mbr_fit        byte         // 1 byte
	Mbr_partitions [4]Partition // 4 * (35) = 140 bytes
	// Total = 157 bytes
}

type Partition struct {
	Part_status      byte     // 1 byte
	Part_type        byte     // 1 byte
	Part_fit         byte     // 1 byte
	Part_start       int32    // 4 bytes
	Part_size        int32    // 4 bytes
	Part_name        [16]byte // 16 bytes
	Part_correlative int32    // 4 bytes
	Part_id          int32    // 4 bytes
	// Total = 35 bytes
}
