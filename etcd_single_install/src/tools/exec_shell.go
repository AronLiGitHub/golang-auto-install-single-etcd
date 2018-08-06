package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"io/ioutil"
)

/**
阻塞方式(需要执行结果)
主要用于执行shell命令，并且返回shell的标准输出
适用于执行普通非阻塞shell命令，且需要shell标准输出的
*/

//阻塞式的执行外部shell命令的函数,等待执行完毕并打印标准输出
func ExecShell(shell string) string {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", shell)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()       //接收错误
	result := out.String() //接收执行结果

	//打印结果
	fmt.Println("执行命令：" + shell)
	fmt.Println("执行结果：" + result)
	checkErr(err) //打印错误

	return result
}

//错误处理函数
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// 测试print
func MyPrint(msg string) {
	fmt.Print(msg)
}

/**
定义写入文件的方法
 */
//使用ioutil.WriteFile方式写入文件,是将[]byte内容写入文件,如果content字符串中没有换行符的话，默认就不会有换行符
func WriteWithIoutil(name,content string) {
	data :=  []byte(content)
	if ioutil.WriteFile(name,data,0644) == nil {
		fmt.Println("写入文件成功:",content)
	}
}
