package markdown

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

type parsedMarkdown map[string][]string

const itemTemplate = `<a href="https://raw.githubusercontent.com/prplx/svg-logos/master/svg/%s.svg"><img src="svg/%s.svg" alt="%s" width="40px" /></a> [%s](https://raw.githubusercontent.com/prplx/svg-logos/master/svg/%s.svg)`

func AddFilesToMarkdown(filepaths []string) (io.Reader, error) {
	var pm = make(parsedMarkdown)
	for _, fp := range filepaths {
		fileName := filepath.Base(fp)
		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		key := strings.ToUpper(string([]rune(fileName)[0]))
		if _, ok := pm[key]; !ok {
			pm[key] = []string{}
		}

		if !isStringInSlice(pm[key], fileNameWithoutExt) {
			pm[key] = pasteToSliceInAlphabeticOrder(pm[key], fileNameWithoutExt)
		}
	}

	keys := make([]string, 0, len(pm))
	for k := range pm {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf := new(strings.Builder)
	for _, k := range keys {
		if buf.Len() == 0 {
			buf.WriteString(fmt.Sprintf("## %s\n\n", k))
		} else {
			buf.WriteString(fmt.Sprintf("\n\n## %s\n\n", k))
		}

		newValue := []string{}

		for _, value := range pm[k] {
			newValue = append(newValue, fmt.Sprintf(itemTemplate, value, value, removeDash(value), removeDash(value), value))
		}

		buf.WriteString((strings.Join(newValue, " | ")))
	}
	return strings.NewReader(buf.String()), nil
}

func pasteToSliceInAlphabeticOrder(sl []string, s string) []string {
	sl = append(sl, s)
	sort.Slice(sl, func(i, j int) bool {
		return strings.ToLower(sl[i]) < strings.ToLower(sl[j])
	})
	return sl
}

func removeDash(s string) string {
	return strings.ReplaceAll(s, "-", " ")
}

func isStringInSlice(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
