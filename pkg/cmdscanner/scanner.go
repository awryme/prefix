package cmdscanner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var Stop = fmt.Errorf("stop")

type ScanFn func(text string) error

func Scan(input io.Reader, prompt string, scanFn ScanFn) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		err := scanFn(text)
		if err == Stop {
			return nil
		}
		if err != nil {
			return fmt.Errorf("scan func: %w", err)
		}
		fmt.Print(prompt)
	}
	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	return nil
}
