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
	file_path := dir+"/"+domain+"/"
	fmt.Println("Searching in:",file_path)		
	files, _ := os.ReadDir(file_path)
	for _, file := range files {
		fileName := file.Name()
		fileList = append(fileList, file_path+fileName)
	}
	return fileList	
}

func printUrl(filename string) {
	array := strings.Split(filename, "/")
	webpath := strings.Replace(array[len(array)-1],"$",":",-1)
	webpath_ := strings.Replace(webpath, "£", "/",-1)
	webpath_final := strings.Replace(webpath_, ".txt","",1)
	fmt.Println(webpath_final)
}
	
func main() {
	var domain string
	fmt.Print("Specify the target domain (only lowercase): ")
	fmt.Scanln(&domain)
	fmt.Println("Type your keyword below (case insensitive):")
	inputReader := bufio.NewReader(os.Stdin)
        input_, _ := inputReader.ReadString('\n')
        input := strings.Replace(input_, "\n", "", -1)
        fmt.Println("Searching for:",input)
	for _, file := range listdir(domain) {
		content, _:= os.ReadFile(file)
		string_content := string(content)
		keyword, _ := regexp.MatchString("(?i)"+input, string_content)
		if keyword {
			printUrl(file)
		}
	}
	fmt.Println("Search finished")
	time.Sleep(time.Second * 100000)
}
