package conf

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Server struct {
		Mode string
		Port string
	}
	Log struct {
		Level   string
		Fpath   string
		Msize   int
		Mage    int
		Mbackup int
	}
}

// level = "debug" # debug or info
// fpath = "./logs/go-loger" # 로그가 생성될 경로 : ./logs, 로그파일명 go-loger_xxx.log
// msize = 2000    # 2g : megabytes
// mage = 7        # 7days
// mbackup = 5    # number of log files
func NewConfig(fpath string) *Config {
	c := new(Config)

	if file, err := os.Open(fpath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		//toml 파일 디코딩
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			return c
		}
	}
}
