package main

import (
	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
	"fmt"
	"os"
)

func main() {
	///home/dabzelos/Downloads
	file, err := os.Open("/home/dabzelos/Downloads/nginx_logs.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	stat := domain.NewStatHolder()
	stat.DataProcessor(file)
	fmt.Printf("%+v\n", stat)
}

/*
	singleLog := "54.86.157.236 - - [17/May/2015:10:05:48 +0000] \"GET /downloads/product_1 HTTP/1.1\" 404 336 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.20.1)\""
	logsFormat := regexp.MustCompile("^(\\S+) - (\\S*) \\[(.*?)] \"(\\S+) (\\S+) (\\S+)\" (\\d{3}) (\\d+) \"(.*?)\" \"(.*?)\"$")

	matches := logsFormat.FindStringSubmatch(singleLog)
	fmt.Println(matches[1])
	fmt.Println(matches[2])
	fmt.Println(matches[3])
	fmt.Println(matches[4])
	fmt.Println(matches[5])
	fmt.Println(matches[6])
	fmt.Println(matches[7])
	fmt.Println(matches[8])
	fmt.Println(matches[9])
	fmt.Println(matches[10])

*/
/*file, err := os.Create("./log-analyzer.adoc")
defer func(file *os.File) {
	err := file.Close()
	if err != nil {

	}
}(file)
file.Write([]byte("### wefefwefwefew\n----------------\n____________________"))
if err != nil {
panic(err)
}
err = file.Close()
if err != nil {
return
}*/
/*package main

import "os"

func main() {
	file, err := os.Create("./log-analyzer.md")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		return
	}

}
*/
/*	// a line of log
	logsExample := `[2021-08-27T07:39:54.173Z] "GET /healthz HTTP/1.1" 200 - 0 61 225 - "111.114.195.106,10.0.0.11" "okhttp/3.12.1" "0557b0bd-4c1c-4c7a-ab7f-2120d67bee2f" "example.com" "172.16.0.1:8080"`

	// your defined log format
	logsFormat := `\[$time_stamp\] \"$http_method $request_path $_\" $response_code - $_ $_ $_ - \"$ips\" \"$_\" \"$_\" \"$_\" \"$_\"`

	// transform all the defined variable into a regex-readable named format
	regexFormat := regexp.MustCompile(`\$([\w_]*)`).ReplaceAllString(logsFormat, `(?P<$1>.*)`)

	// compile the result
	re := regexp.MustCompile(regexFormat)

	// find all the matched data from the logsExample
	matches := re.FindStringSubmatch(logsExample)

	for i, k := range re.SubexpNames() {
		// ignore the first and the $_
		if i == 0 || k == "_" {
			continue
		}

		// print the defined variable
		fmt.Printf("%-15s => %s\n", k, matches[i])
	}*/
