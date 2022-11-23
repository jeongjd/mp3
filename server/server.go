package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Message struct {
	messageContent string
	node           string
	round          int
}

var (
	// Map - key: (client) username, value: connection
	clientConnections = make(map[string]net.Conn)

	// Read/Write mutex to synchronize the clientConnections hashmap between the threads (instead of a channel)
	clientConnectionsMutex = sync.RWMutex{}
)

func main() {
	// var port string
	// fmt.Print("Enter a port number: ")
	// fmt.Scanln(&port)
	// port = ":" + port
	totalNodes := 10
	portNum := 1111
	for i := 0; i < totalNodes; i++ {
		portNum += i
		port := strconv.Itoa(portNum)
		fmt.Println("Launching a TCP Chatroom Server...")
		l, err := net.Listen("tcp4", ":"+port)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer l.Close()
		messages := make(chan Message)
		go createTCPClient(port, messages)
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func createTCPClient(port string, messages chan Message) {
	fmt.Println("creating a TCP client...")
	hostAddress := "127.0.0.1"
	c, err := net.Dial("tcp", hostAddress+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	go read(c, messages)
	write(c)
}

// Handle client connections - invoke other functions depending on the messages received
func handleConnection(c net.Conn) {
	var count = 0
	for {
		var text string
		// Reads and decodes data from connection
		dec := gob.NewDecoder(c)
		err := dec.Decode(&text)

		// If a connection is closed delete the username from map (clientConnections)
		if err != nil {
			log.Fatal(err)
		}

		// var m Message
		// broadcastMessage()

		// unicastSend(c, m.messageContent, delay)
		count++
	}
}

// Split a string into a string array
func parseLine(line string) []string {
	return strings.Split(line, " ")
}

// Send private message to a specific client using gob
func broadcastMessage(m Message) {
	// Prevents other go routines from reading the clientConnections hashmap in order to synchronize the routines
	clientConnectionsMutex.RLock()
	defer clientConnectionsMutex.RUnlock()
	for item := range clientConnections {
		enc := gob.NewEncoder(clientConnections[item])
		enc.Encode(m.messageContent)
	}
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

func delayTime() int {
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	delay := strings.Fields(scanner.Text())
	minDelay, _ := strconv.Atoi(delay[0])
	maxDelay, _ := strconv.Atoi(delay[1])
	delayed := rand.Intn(maxDelay-minDelay+1) + minDelay
	return delayed
}

func unicastSend(c net.Conn, m Message, delayed int) {
	go func() {
		// delayTime()
		// implement delay
		c.Write([]byte(m.messageContent))
	}()
}

func approximateConsensus() {
	// round := 1
	// sum := 0
	// messageNum := 0
}

func read(c net.Conn, messages chan Message) {
	for {
		dec := gob.NewDecoder(c)
		msg := new(Message)
		err := dec.Decode(&msg)
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		messages <- *msg
		fmt.Println(msg)
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

		enc := gob.NewEncoder(c)
		if err := enc.Encode(message); err != nil {
			log.Fatal(err)
		}
	}
}
