package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	hostAddress  = " "
	port         = " "
	nodeNums     = 0
	numFailures  = 0
	initialValue = 0
	addresses    = make(map[int][]string)
)

func main() {
	addresses = readConfig()
	fmt.Print("How many nodes? ")
	fmt.Scanln(&nodeNums)
	fmt.Print("What is the upper bound on the number of failures? ")
	fmt.Scanln(&numFailures)
	// Must be moved to somewhere else
	fmt.Print("What is the initial value ")
	fmt.Scanln(&initialValue)

	createTCPClients(nodeNums)
}
func createTCPClients(nodeNums int) {
	for i := 1; i < nodeNums; i++ {
		go createTCPClient(string(addresses[i][0]+addresses[i][1]), nodeNums)
	}
}

func parseLine(line string) []string {
	return strings.Split(line, " ")
}

func readConfig() map[int][]string {
	// open "config.txt" file
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Delay closing of the file until other functions return
	defer file.Close()

	configData := make(map[int][]string)
	scanner := bufio.NewScanner(file)
	currentLineNum := 0
	configLine := ""
	for scanner.Scan() {
		configLine = (scanner.Text())
		configLineParsed := parseLine(configLine)
		configData[currentLineNum] = configLineParsed
		currentLineNum++
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return configData
}

func createTCPClient(address string, nodeCount int) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	//Creating the connections and storing them

	go read(c)
	write(c)
}

func read(c net.Conn) {
	for {
		var message string
		dec := gob.NewDecoder(c)
		err := dec.Decode(&message)
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		// Close client connection if the server has shut down
		if err == io.EOF {
			fmt.Println("Server has shut down. Client is now closing... ")
			c.Close()
			os.Exit(0)
		}
		fmt.Println(message)
		fmt.Print(">> ")
	}
}

func write(c net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if strings.Contains(message, "EXIT") {
			fmt.Println("Exiting the client...")
			return
		}
		enc := gob.NewEncoder(c)
		if err := enc.Encode(message); err != nil {
			log.Fatal(err)
		}
	}
}
