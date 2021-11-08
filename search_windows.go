package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"bufio"
	"time"
)

func listdir(domain string) []string {
	var fileList []string
	dir, _ := os.Getwd()
	file_path := dir+"\\"+domain+"\\"
	fmt.Println("Ricerca in:",file_path)		
	files, _ := os.ReadDir(file_path)
	for _, file := range files {
		fileName := file.Name()
		fileList = append(fileList, file_path+fileName)
	}
	return fileList	
}

func printUrl(filename string) {
	array := strings.Split(filename, "\\")
	webpath := strings.Replace(array[len(array)-1],"$",":",-1)
	webpath_ := strings.Replace(webpath, "£", "/",-1)
	webpath_final := strings.Replace(webpath_, ".txt","",1)
	fmt.Println(webpath_final)
}
	
func main() {
	var domain string
	fmt.Print("Specificare il dominio (solo caratteri minuscoli): ")
	fmt.Scanln(&domain)
	fmt.Println("Specificare la stringa di ricerca (case insensitive):")
	inputReader := bufio.NewReader(os.Stdin)
        input_, _ := inputReader.ReadString('\n')
        input := strings.Replace(input_, "\r\n", "", -1)
        fmt.Println("Ricerca di:",input)
	for _, file := range listdir(domain) {
		content, _:= os.ReadFile(file)
		string_content := string(content)
		keyword, _ := regexp.MatchString("(?i)"+input, string_content)
		if keyword {
			printUrl(file)
		}
	}
	fmt.Println("Ricerca terminata")
	time.Sleep(time.Second * 100000)
}
