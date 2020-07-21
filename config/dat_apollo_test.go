package config_test

import (
	"fmt"
	"testing"

	"github.com/ydbt/devtool/v1/config"
)

func TestApolloString(t *testing.T) {
	cfg := make(map[string]string)
	cfg["log.appname"] = "config-log"
	cfg["log.path"] = "ut-apollo-cfg.log"
	cfg["log.level"] = "debug"
	actualYaml := config.Apollo2Yaml(cfg)
	expectYaml := `ApolloCfgUt:
  log:
    appname: config-log
    path: ut-apollo-cfg.log
    level: debug
`
	if actualYaml != expectYaml {
		t.Error("parse failed:", actualYaml)
		actualByte := []byte(actualYaml)
		expectByte := []byte(expectYaml)
		for i, _ := range actualByte {
			if actualByte[i] != expectByte[i] {
				fmt.Println(string(actualYaml[i:]))
				fmt.Println(string(expectYaml[i:]))
				break
			}
		}
	}
}

/*
 */
func TestApolloArray01(t *testing.T) {
	cfg := make(map[string]string)
	cfg["services.ip.[0]"] = "192.168.1.1:100861"
	cfg["services.ip.[1]"] = "127.0.0.1:10010"
	actualYaml := config.Apollo2Yaml(cfg)
	expectYaml := `services:
  ip:
    - 192.168.1.1:100861
    - 127.0.0.1:10010
`
	fmt.Println(actualYaml)
	fmt.Println(expectYaml)
	//	if actualYaml != expectYaml {
	//		t.Error("parse failed:\n", actualYaml)
	//		actualByte := []byte(actualYaml)
	//		expectByte := []byte(expectYaml)
	//		for i, _ := range actualByte {
	//			if actualByte[i] != expectByte[i] {
	//				fmt.Printf("%v != %v\n", actualByte[i], expectByte[i])
	//				fmt.Println(string(actualYaml[i:]))
	//				fmt.Println(string(expectYaml[i:]))
	//				break
	//			}
	//		}
	//	}
}

func TestApolloArray02(t *testing.T) {
	cfg := make(map[string]string)
	cfg["services.[0].host"] = "192.168.1.1"
	cfg["services.[0].port"] = "100861"
	actualYaml := config.Apollo2Yaml(cfg)
	expectYaml := `services:
  -
    port: 100861
    host: 192.168.1.1
`
	fmt.Println(actualYaml)
	fmt.Println(expectYaml)
	//	if actualYaml != expectYaml {
	//		t.Error("parse failed:", actualYaml)
	//		actualByte := []byte(actualYaml)
	//		expectByte := []byte(expectYaml)
	//		for i, _ := range actualByte {
	//			if actualByte[i] != expectByte[i] {
	//				fmt.Printf("%v != %v\n", actualByte[i], expectByte[i])
	//				fmt.Println(string(actualYaml[i:]))
	//				fmt.Println(string(expectYaml[i:]))
	//				break
	//			}
	//		}
	//	}
}
