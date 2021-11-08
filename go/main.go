package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/thatisuday/commando"
)

type SafeBuffer struct {
	sync.Mutex
	buffer *bufio.Reader
}

func (b *SafeBuffer) Read() (string, error) {
	b.Lock()
	defer b.Unlock()

	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = b.buffer.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

func producer(buffer *SafeBuffer, lines chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		text, err := buffer.Read()

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
			break
		} else {
			lines <- text
		}
	}
}

func consumer(text_to_find string, lines <-chan string, counts chan<- int) {
	count := 0

	for line := range lines {
		if strings.Contains(line, text_to_find) {
			count = count + 1
		}
	}

	counts <- count
}

func execute(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	file_path := args["file"].Value
	text_to_find := args["text"].Value
	num_producers, _ := flags["producers"].GetInt()
	num_consumers, _ := flags["consumers"].GetInt()

	fmt.Printf("Procurando %v no arquivo %v\n", text_to_find, file_path)

	file, err := os.Open(file_path)

	if err != nil {
		log.Fatal(err)
	}

	buffer := SafeBuffer{
		buffer: bufio.NewReader(file),
	}

	var wg sync.WaitGroup
	line_channel := make(chan string, 100)
	counts_channel := make(chan int, num_consumers)

	for i := 0; i < num_producers; i++ {
		wg.Add(1)
		go producer(&buffer, line_channel, &wg)
	}

	for i := 0; i < num_consumers; i++ {
		go consumer(text_to_find, line_channel, counts_channel)
	}

	wg.Wait()
	close(line_channel)

	finished_routines := 0
	total := 0

	for count := range counts_channel {
		total = total + count

		finished_routines = finished_routines + 1

		if finished_routines == num_consumers {
			close(counts_channel)
		}
	}

	fmt.Printf("O texto %v foi encontrado %v vez(es) no arquivo %v", text_to_find, total, file_path)
}

func main() {
	commando.
		SetExecutableName("desafio-4").
		SetVersion("1.0.0").
		SetDescription("Buscador de texto em um arquivo")

	commando.
		Register(nil).
		AddArgument("file", "Arquivo de texto a ser lido", "").
		AddArgument("text", "Texto a ser procurado", "").
		AddFlag("producers,P", "Número de threads produtoras", commando.Int, runtime.NumCPU()/2).
		AddFlag("consumers,C", "Número de threads consumidoras", commando.Int, runtime.NumCPU()/2).
		SetAction(execute)

	commando.Parse(nil)
}
