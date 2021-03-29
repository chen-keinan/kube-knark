package utils

//FindSeverityInt return severity int value
func FindSeverityInt(severity string) int {
	switch severity {
	case "CRITICAL":
		return 1
	case "MAJOR":
		return 2
	case "MINOR":
		return 3
	case "INFO":
		return 4
	}
	return 0
}
