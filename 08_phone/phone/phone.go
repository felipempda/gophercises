package phone

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
)

func SampleData() []string {
	// turn this into an array
	entries := `1234567890
	123 456 7891
	(123) 456 7892
	(123) 456-7893
	123-456-7894
	123-456-7890
	1234567892
	(123)456-7892`

	result := make([]string, 0)
	reader := bufio.NewReader(strings.NewReader(entries))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		line = strings.TrimSpace(line)
		result = append(result, line)
		//fmt.Println(line)
	}
	return result
}

func Normalize(number string) string {
	var result bytes.Buffer
	for _, r := range number {
		ascii := int(r)
		if ascii >= 48 && ascii <= 57 {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func NormalizeRegex(number string) string {
	re := regexp.MustCompile("[^0-9]") // or //D
	return re.ReplaceAllString(number, "")
}
