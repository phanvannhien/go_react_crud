package validator

import "os/exec"

func RunAll() []string {

	var errors []string

	if err := exec.Command("go", "vet", "./...").Run(); err != nil {
		errors = append(errors, "go vet failed")
	}

	if err := exec.Command("go", "test", "./...").Run(); err != nil {
		errors = append(errors, "go test failed")
	}

	// SECURITY SCAN
	issues := ScanSecurity(".")
	for _, i := range issues {
		errors = append(errors,
			i.Rule+" in "+i.File,
		)
	}

	// OPTIONAL GOSEC
	gosecErrors := RunGoSec()
	errors = append(errors, gosecErrors...)

	return errors
}
