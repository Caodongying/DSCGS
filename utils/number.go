package utils

import (
	"strconv"
	"sync"
)


var once sync.Once
var base62 [62]string

func initBase62() {
	// 初始化基底
	once.Do(func(){
		index := 0
		for i := 0; i <= 9; i++ {
			base62[index] = strconv.Itoa(i)
			index ++
		}
		for i:= 'a'; i <= 'z'; i++ {
			base62[index] = string(i)
			index++
		}
		for i:= 'A'; i <= 'Z'; i++ {
			base62[index] = string(i)
			index++
		}
	})
}

func DecimalTo62(num int) string {
	/*
	把十进制数转化为62进制数

	不断对62取模，直到余数为0
	*/
	initBase62()
	result := ""
	var mod int // 模数

	for {
		mod = num % 62
		result = base62[mod] + result

		num = num / 62
		if (num == 0) {
			break
		}
	}

	return result
}
