package main

import (
	"bufio"
	"fmt"
	//"os"
	//"io"
	"net"
)

func main() {
	fmt.Println("Hello World!");
	// fmt.Println("Starting server...");
	// server, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090));
	// if err != nil {
	// 	fmt.Println("Error starting server:", err.Error());
	// 	os.Exit(1);
	// }
	// fmt.Println("Server started.");
	// defer server.Close();
	// fmt.Println("Waiting for client...");
	// isAllEventsReceived := false;
	// for !isAllEventsReceived {
	// 	connection, err := server.Accept();
	// 	if err != nil {
	// 		fmt.Println("Error in client connection:", err.Error());
	// 		os.Exit(1);
	// 	}
	// 	fmt.Println("client connected");
	// 	// do go routine here
	// 	buffer := make([]byte, 1024);
	// 	messageLength, err := connection.Read(buffer);
	// 	if err != nil {
	// 		fmt.Println("Error reading:", err.Error());
	// 	}
	// 	message := string(buffer[:messageLength]);
	// 	fmt.Println("Message received: ", message);
	// 	if message == "EVENT END" {
	// 		isAllEventsReceived = true;
	// 		connection.Close();
	// 	}
	// }
	connection, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", 9090));
	if err != nil {
		fmt.Println("Error connecting to EVENT SOURCE:", err.Error());
	}
	defer connection.Close()
	reader := bufio.NewReader(connection);
	fmt.Println("Waiting for server...");
	isAllEventsReceived := false;
	for !isAllEventsReceived {
		//buffer := make([]byte, 1024);
		//messageLength, err := connection.Read(buffer);
		message, err := reader.ReadString('\n');
		if err != nil {
			fmt.Println("Error reading:", err.Error());
		}
		//message := string(buffer[:messageLength]);
		message = message[:len(message)-1]
		fmt.Println("Message received: ", message);
		if message == "EVENT END" {
			isAllEventsReceived = true;
			//connection.Close();
			//break;
		}
	}
	//connection.Close();
	fmt.Println("All messages received.");
	fmt.Println("Connection closed.");
	fmt.Println("Server shutdown.");
}