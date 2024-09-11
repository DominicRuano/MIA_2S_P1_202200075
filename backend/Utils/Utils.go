package utils

func CalcularTamaño(size int, unit string) int {
	if unit == "K" {
		return size * 1024
	} else if unit == "M" {
		return size * 1024 * 1024
	} else {
		return -1
	}
}
