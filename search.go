package main
import (
	"fmt"
	"os"
	"strings"
	"bufio"
)

func listdir(domain string) []string {
	var fileList []string
	dir, _ := os.Getwd()
	file_path := dir+"/"+domain+"/"
	fmt.Println("Searching in:",file_path,"\n")
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
	webpath__ := strings.Replace(webpath_, "!!!", ":",-1)
	webpath_final := strings.Replace(webpath__, ".txt","",1)
	fmt.Println(webpath_final)
}

func main() {
	var domain string
	fmt.Print("Specify the target domain (only lowercase): ")
	fmt.Scanln(&domain)
	fmt.Print("Type your keyword below (case insensitive): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.ToLower(scanner.Text())
        fmt.Println("Searching for:",input)
	for _, file_ := range listdir(domain) {
		file, _ := os.Open(file_)
		defer file.Close()
		scanner = bufio.NewScanner(file)
		for scanner.Scan() {
		    line := strings.ToLower(scanner.Text())
		    if strings.Contains(line, input) {
		        printUrl(file_)
			break
		     }
		  }
	 }
	 fmt.Println("Search finished, press enter to close window")
	 fmt.Scanln()
}
