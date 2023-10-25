package utils

import "fmt"

func GetFormatText(texts map[string]string, key string, values ...any) string {
	return fmt.Sprintf(texts[key], values...)
}
