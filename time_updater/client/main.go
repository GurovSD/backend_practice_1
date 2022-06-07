package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1) // создаем буфер
	for {
		_, err = conn.Read(buf)
		if err == io.EOF {
			break
		}
		// io.WriteString(os.Stdout, fmt.Sprintf("Custom output! %s", string(buf)))
		// buf = buf[:0]

		//меняем использование буфера на копирование содержимого соединения, чтобы избежать "хвостов" от предыдущего сообщения.
		// использование буфера на 1 символ оставляем как проверку открытого соединения
		io.Copy(os.Stdout, conn)
		// выводим измененное сообщение сервера в консоль
	}
}
