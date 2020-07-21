package config

import (
	"bytes"
	"fmt"
	"strings"
)

type ApolloCfg struct {
	AppID      string   `json:"appid",yaml:"appid"`
	Cluster    string   `json:"cluster",yaml:"cluster"`
	Namespaces []string `json:"namespaces",yaml:"namespaces"`
	IP         string   `json:"ip",yaml:"ip"`
	Secret     string   `json:"secret",yaml:"secret"`
}

/* --------------------------------    -------------------------------- */
func Apollo2Yaml(cfgs map[string]string) string {
	root := newElemNode(cfgs, "apollo-cfg")
	var iob bytes.Buffer

	switch root.nodeType {
	case typeObject:
		mapChild := root.value.(map[string]*elemNode)
		for _, v := range mapChild {
			print(v, "", &iob)
		}
	case typeString:
		strChild := root.value.(string)
		iob.WriteString(fmt.Sprintf("%s: %s\n", root.name, strChild))
	}
	return string(iob.Bytes())
}

const (
	typeString = 0
	typeObject = 1
)

// ElemNode 元素节点
// typeString value=>string("")
// typeObject value=>map[string]*ElemNode  // 对key进行了排序，Print数据可能顺序有点差别
// typeArray  value=>[]*ElemNode
type elemNode struct {
	nodeType int
	name     string
	parent   *elemNode
	value    interface{}
}

func newElemNode(cfg map[string]string, appName string) *elemNode {
	root := &elemNode{
		nodeType: typeObject,
		name:     appName,
		parent:   nil,
		value:    make(map[string]*elemNode),
	}
	for k, v := range cfg {
		attributeNames := strings.Split(k, ".")
		insert(root, attributeNames, v)
	}
	return root
}

func insert(root *elemNode, attributeNames []string, val string) {
	attributeSize := len(attributeNames)
	if root == nil || attributeSize == 0 {
		return
	}
	typeNode := typeString
	attributeName := attributeNames[0]
	if attributeSize > 1 {
		typeNode = typeObject
	}
	if attributeSize > 1 {
		switch root.nodeType {
		case typeObject:
			mapChild := root.value.(map[string]*elemNode)
			child, ok := mapChild[attributeName]
			if !ok {
				child = &elemNode{
					nodeType: typeNode,
					name:     attributeName,
					parent:   root,
					value:    make(map[string]*elemNode),
				}
			}
			mapChild[attributeName] = child
			insert(child, attributeNames[1:], val)
		}
	} else {
		switch root.nodeType {
		case typeObject:
			mapChild := root.value.(map[string]*elemNode)
			child, ok := mapChild[attributeName]
			if !ok {
				child = &elemNode{
					nodeType: typeNode,
					name:     attributeName,
					parent:   root,
					value:    val,
				}
			}
			mapChild[attributeName] = child
		}
	}
}

func print(root *elemNode, prefixIndent string, iob *bytes.Buffer) {
	stepIndent := "  "
	switch root.nodeType {
	case typeObject:
		if root.name[0] == '[' {
			iob.WriteString(fmt.Sprintf("%s-\n", prefixIndent))
		} else {
			iob.WriteString(fmt.Sprintf("%s%s:\n", prefixIndent, root.name))
		}
		mapChild := root.value.(map[string]*elemNode)
		for _, v := range mapChild {
			print(v, prefixIndent+stepIndent, iob)
		}
	case typeString:
		strChild := root.value.(string)
		if root.name[0] == '[' {
			iob.WriteString(fmt.Sprintf("%s- %s\n", prefixIndent, strChild))
		} else {
			iob.WriteString(fmt.Sprintf("%s%s: %s\n", prefixIndent, root.name, strChild))
		}
	}
}

/**
  TODO 目前解析apollo处理逻辑有点乱，有待重构 [2020-07-21]
*/
