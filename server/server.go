package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Message struct {
	receiverID     string
	senderID       string
	messageContent string
}

var (
	// Map - key: (client) username, value: connection
	clientConnections = make(map[string]net.Conn)

	// For switch/cases - printing error messages
	option = 0

	// Read/Write mutex to synchronize the clientConnections hashmap between the threads (instead of a channel)
	clientConnectionsMutex = sync.RWMutex{}
)

func main() {
	var port string
	fmt.Print("Enter a port number: ")
	fmt.Scanln(&port)
	port = ":" + port

	fmt.Println("Launching a TCP Chatroom Server...")

	l, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
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
		if item == m.receiverID {
			enc := gob.NewEncoder(clientConnections[item])
			enc.Encode(m.messageContent)
		}
	}
}

func readConfig() {
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// store data into a hashmap (host address, port etc)
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
