package filter

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/importcjj/sensitive"
)

var filter *sensitiveFilter

func InitFilter(pathToDict string) {
	if pathToDict == "" {
		panic("请指定敏感词库路径")
	}
	_, err := os.Stat(pathToDict)
	if os.IsNotExist(err) {
		panic("敏感词库文件不存在")
	}
	f := sensitive.New()
	f.UpdateNoisePattern("")
	f.LoadWordDict(pathToDict)
	filter = &sensitiveFilter{
		filter: f,
	}
}

func OverwriteDict(pathToDict string) error {
	file, err := os.Open(pathToDict)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`\p{Han}+`)
	newContent := ""
	kwMp := make(map[string]struct{}, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 去重
		if _, ok := kwMp[line]; ok {
			continue
		}
		kwMp[line] = struct{}{}
		match := re.FindString(line)
		if match == "" {
			newContent += " " + strings.Trim(line, " ") + " \n"
		} else {
			newContent += strings.Trim(line, " ") + "\n"
		}
	}
	newContent = strings.Trim(newContent, "\n")
	file.Close()
	os.Remove(pathToDict)
	file, err = os.OpenFile(pathToDict, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(newContent)
	if err != nil {
		panic(err)
	}
	file.Close()
	return nil
}

func (sf *sensitiveFilter) Validate(text string) (bool, string) {
	text = " " + strings.Trim(text, " ") + " "
	ok, word := sf.filter.Validate(text)
	word = strings.Trim(word, " ")
	return ok, word
}

type ISensitiveFilter interface {
	Validate(text string) (bool, string)
}

type sensitiveFilter struct {
	filter *sensitive.Filter
}

func GetFilter() ISensitiveFilter {
	return filter
}
