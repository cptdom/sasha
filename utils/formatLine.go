package utils

func CreateLine() string {
	var line string
	for i := 0; i < 70; i++ {
		line = line + "-"
	}
	return line
}
