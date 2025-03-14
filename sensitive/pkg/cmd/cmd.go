package cmd

import "flag"

type CommandArgs struct {
	Dict   string
	Config string
}

var Args *CommandArgs

func init() {
	dict := flag.String("dict", "dict.txt", "敏感词汇词库")
	config := flag.String("config", "config.yaml", "配置文件")
	flag.Parse()
	Args = &CommandArgs{}
	Args.Dict = *dict
	Args.Config = *config
}
