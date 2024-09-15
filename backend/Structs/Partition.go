package structs

type Partition struct {
	Part_status      [1]byte  // 1 byte
	Part_type        [1]byte  // 1 byte
	Part_fit         [1]byte  // 1 byte
	Part_start       int32    // 4 bytes
	Part_size        int32    // 4 bytes
	Part_name        [16]byte // 16 bytes
	Part_correlative int32    // 4 bytes
	Part_id          [4]byte  // 4 bytes
	// Total = 35 bytess
}
