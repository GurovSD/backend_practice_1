package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	// добавляем броадкастер для рассылки сообщений с сервера по всем подключенным клиентам
	go broadcaster()
	// добавляем функцию, отправляющую в канал сообщений текущее время
	go timeUpdater()
	// добавляем функцию, обрабатывающую входящие соединения
	go connDispatcher(listener)

	// запускаем чтение из консоли с последующей отправкой в канал сообщений для рассылки на всех клиентов
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		messages <- text
		fmt.Println(text)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func connDispatcher(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)

	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	entering <- ch

	for msg := range ch {
		_, err := io.WriteString(conn, msg)
		if err != nil {
			return
		}

	}
	leaving <- ch

	conn.Close()
}

func timeUpdater() {
	for {
		messages <- time.Now().Format("15:04:05\n\r")
		time.Sleep(2 * time.Second)
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(conn, msg)
		if err != nil {
			return
		}

	}
}
