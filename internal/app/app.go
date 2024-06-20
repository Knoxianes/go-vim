package app

var mode = 0
var copyRegister = ""

const NormalMode = 0
const InsertMode = 1
const VisualMode = 2

func GetMode() int {
	return mode
}
func SetMode(m int) {
	if m < 0 || m > 2 {
		panic("Invalid mode")
	}
	mode = m
}
func GetCopyRegister() string {
	return copyRegister
}
func SetCopyRegister(r string) {
	copyRegister = r
}
