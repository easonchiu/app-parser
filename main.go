/*
 * @Author: easonchiu
 * @Date: 2023-07-03 10:53:15
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-03 19:17:05
 * @Description:
 */
package main

import (
	"fmt"

	"github.com/easonchiu/app-parser/parser"
)

func main() {
	data, _ := parser.ParseAPPData("1142110895")
	fmt.Println(data)
}
