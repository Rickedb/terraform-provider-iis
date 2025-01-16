package agent

func toPascalCase(value bool) string {
	bVal := "False"
	if value {
		bVal = "True"
	}

	return bVal
}
