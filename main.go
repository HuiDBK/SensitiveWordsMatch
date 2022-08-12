package main

import (
	"fmt"
	"github.com/chain-zhang/pinyin"
	"regexp"
	"strings"
)

// HansCovertPinyin 中文汉字转拼音
func HansCovertPinyin(contents []string) []string {
	pinyinContents := make([]string, 0)
	for _, content := range contents {
		chineseReg := regexp.MustCompile("[\u4e00-\u9fa5]")
		if !chineseReg.Match([]byte(content)) {
			continue
		}

		// 只有中文才转
		pin := pinyin.New(content)
		pinStr, err := pin.Convert()
		fmt.Println(content, "->", pinStr)
		if err == nil {
			pinyinContents = append(pinyinContents, pinStr)
		}
	}
	return pinyinContents
}

func runeTest() {
	fmt.Println("a -> ", rune('a'))
	fmt.Println("A -> ", rune('A'))
	fmt.Println("晖 -> ", rune('晖'))
	fmt.Println("霞 -> ", rune('霞'))
	fmt.Println("晖霞 -> ", []rune("晖霞"))
}

// 暴力匹配
func normalDemo(sensitiveWords []string, matchContents []string) {

	for _, text := range matchContents {
		srcText := text
		for _, word := range sensitiveWords {
			replaceChar := ""

			for i, wordLen := 0, len([]rune(word)); i < wordLen; i++ {
				// 根据敏感词的长度构造和谐字符
				replaceChar += "*"
			}

			text = strings.Replace(text, word, replaceChar, -1)
		}
		fmt.Println("srcText     -> ", srcText)
		fmt.Println("replaceText -> ", text)
		fmt.Println()
	}

}

// 正则匹配敏感词
func regDemo(sensitiveWords []string, matchContents []string) {

	banWords := make([]string, 0) // 收集匹配到的敏感词

	// 构造正则匹配字符
	regStr := strings.Join(sensitiveWords, "|")
	wordReg := regexp.MustCompile(regStr)
	println("regStr -> ", regStr)

	for _, text := range matchContents {
		textBytes := wordReg.ReplaceAllFunc([]byte(text), func(bytes []byte) []byte {
			banWords = append(banWords, string(bytes))
			textRunes := []rune(string(bytes))
			replaceBytes := make([]byte, 0)
			for i, runeLen := 0, len(textRunes); i < runeLen; i++ {
				replaceBytes = append(replaceBytes, byte('*'))
			}
			return replaceBytes
		})
		fmt.Println("srcText        -> ", text)
		fmt.Println("replaceText    -> ", string(textBytes))
		fmt.Println("sensitiveWords -> ", banWords)
		fmt.Println()
	}
}

// 前缀树匹配敏感词
func trieDemo(sensitiveWords []string, matchContents []string) {

	// 汉字转拼音
	pinyinContents := HansCovertPinyin(sensitiveWords)
	trie := NewSensitiveTrie()
	trie.AddWords(sensitiveWords)
	trie.AddWords(pinyinContents) // 添加拼音敏感词

	//trie.AddWords(pinyinContents)
	//for _, content := range contents {
	//	trie.AddWord(content)
	//}

	for _, srcText := range matchContents {
		matchSensitiveWords, replaceText := trie.Match(srcText)
		fmt.Println("srcText        -> ", srcText)
		fmt.Println("replaceText    -> ", replaceText)
		fmt.Println("sensitiveWords -> ", matchSensitiveWords)
		fmt.Println()
	}

	// 动态添加
	trie.AddWord("牛大大")
	content := "今天，牛大大签发军令"
	matchSensitiveWords, replaceText := trie.Match(content)
	fmt.Println("srcText        -> ", content)
	fmt.Println("replaceText    -> ", replaceText)
	fmt.Println("sensitiveWords -> ", matchSensitiveWords)
}

func main() {
	fmt.Println("--------- rune测试 ---------")
	runeTest()

	sensitiveWords := []string{
		"傻逼",
		"傻叉",
		"垃圾",
		"妈的",
		"sb",
	}

	fmt.Println("\n--------- 汉字转拼音 ---------")
	pinyinContents := HansCovertPinyin(sensitiveWords)
	fmt.Println(pinyinContents)

	matchContents := []string{
		"你是一个大傻&逼，大傻 叉",
		"你是傻☺叉",
		"shabi东西",
		"他made东西",
		"什么垃圾打野，傻逼一样，叫你来开龙不来，SB",
		"正常的内容☺",
	}

	fmt.Println("\n--------- 普通暴力匹配敏感词 ---------")
	normalDemo(sensitiveWords, matchContents)

	fmt.Println("\n--------- 正则匹配敏感词 ---------")
	regDemo(sensitiveWords, matchContents)

	fmt.Println("\n--------- 前缀树匹配敏感词 ---------")
	trieDemo(sensitiveWords, matchContents)

}
