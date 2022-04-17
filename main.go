package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	showIntro()
	for {
		showOptions()
		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("here are the logs:")
			showLogs()
		case 0:
			fmt.Println("bye...")
			os.Exit(0)
		default:
			fmt.Println("invalid option!")
			os.Exit(-1)
		}
	}

}

func showOptions() {
	fmt.Println("1- Start website monitoring")
	fmt.Println("2- Show logs ")
	fmt.Println("0- Exit")
}

func showIntro() {
	name := "Marincor"
	github := "@" + name
	version := "1.0"
	fmt.Println("Hello,", "for more access the github", github)
	fmt.Println("This software version is:", version)
}

func readCommand() (command int) {
	fmt.Scan(&command)
	return command
}

func startMonitoring() {
	const delay = 5 * time.Minute
	const monitorings = 5
	fmt.Println("monitoring...")
	sites := readSitesTxt()

	for i := 0; i < monitorings; i++ {
		for _, site := range sites {
			siteTesting(site)
		}
		time.Sleep(delay)
	}

}

func siteTesting(site string) {
	res, err := http.Get(site)
	if err != nil {
		fmt.Println("something is wrong: ", err)
		log.Panic(err)
	}

	switch res.StatusCode {
	case 200:
		fmt.Println("Site:", site, "loading with success!")
		registerLog(site, true)
	default:
		fmt.Println("Site", site, "with error, status code:", res.StatusCode)
		registerLog(site, false)
	}
}

func readSitesTxt() []string {
	var sites []string
	file, err := os.Open("./sites.txt")
	if err != nil {
		fmt.Println("something is wrong: ", err)
		log.Panic(err)
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		sites = append(sites, line)
		if err == io.EOF {
			break
		}
	}
	file.Close()
	return sites
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("./sites_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("something is wrong, error: ", err)
	}
	timeFormat := "02/01/2006 15:04:05" /* dd/mm/yyyy hh:mm:ss */
	timenow := time.Now().Format(timeFormat)
	breakLine := "\n"

	file.WriteString(timenow + " - " + "The site: " + site + "is online:" + strconv.FormatBool(status) + breakLine)

	file.Close()
}

func showLogs() {
	file, err := ioutil.ReadFile("./sites_log.txt")
	if err != nil {
		fmt.Println("shomething is wrong: ", err)
	}
	fmt.Println(string(file))

}
