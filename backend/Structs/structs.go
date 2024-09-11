package structs

type MBR struct {
	Mbr_size      int32   // 4 bytes
	Mbr_date      float64 // 8 bytes
	Mbr_signature int32   // 4 bytes
	// Total = 16 bytes
}
