package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type client chan<- string

var (
	entering    = make(chan client)
	leaving     = make(chan client)
	messages    = make(chan string)
	tasks       = make(chan string)
	clients     = make(map[client]bool)
	rightAnswer string
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	go taskManager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
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

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	fmt.Fprintln(conn, "Your address is "+who+".\nEnter your nickname and press enter")
	input := bufio.NewScanner(conn)
	for input.Scan() {
		who = input.Text()
		break
	}

	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	for input.Scan() {
		if input.Text() == rightAnswer {
			messages <- who + " give the right answer! It is " + rightAnswer
			<-tasks
		}
	}
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func taskManager() {
	for {
		if len(clients) > 0 {
			expression := generateTaskExpression()
			messages <- "New task: " + expression
			tasks <- expression
		}
	}
}

func generateTaskExpression() string {
	operators := []string{"+", "-", "*"}

	rand.Seed(time.Now().UnixNano())
	// for {
	selectedOperator := operators[rand.Intn(3)]
	operand_1 := rand.Intn(101)
	operand_2 := rand.Intn(101)
	expression := strconv.Itoa(operand_1) + " " + selectedOperator + " " + strconv.Itoa(operand_2)

	var answer int
	switch selectedOperator {
	case "+":
		answer = operand_1 + operand_2
	case "-":
		answer = operand_1 - operand_2
	case "*":
		answer = operand_1 * operand_2
	}

	rightAnswer = strconv.Itoa(answer)
	fmt.Println(expression)

	return expression
	// }

}
