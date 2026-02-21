package diff

import (
	"os/exec"
)

func BuildContext(files []string) (string, error) {

	var ctx string

	for _, f := range files {
		out, err := exec.Command("git", "diff", "HEAD", f).CombinedOutput()
		if err != nil {
			return "", err
		}

		ctx += "\nFile: " + f + "\n" + string(out)
	}

	return ctx, nil
}
