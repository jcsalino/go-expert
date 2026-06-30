package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	c1 := make(chan string)
	c2 := make(chan string)

	for _, cep := range os.Args[1:] {

		go func() {
			req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
			if err != nil {
				c1 <- "ViaCEP: Error ao fazer a requisicao: " + err.Error()
			}
			res, err := io.ReadAll(req.Body)
			if err != nil {
				c1 <- "ViaCEP: Error ao ler a resposta: " + err.Error()
			}
			c1 <- "ViaCEP: " + string(res)

		}()
		go func() {
			// time.Sleep(2 * time.Second)
			req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
			if err != nil {
				c2 <- "BrasilAPI: Error ao fazer a requisicao: " + err.Error()
			}
			res, err := io.ReadAll(req.Body)
			if err != nil {
				c2 <- "BrasilAPI: Error ao ler a resposta: " + err.Error()
			}
			c2 <- "BrasilAPI: " + string(res)

		}()

	}
	for {
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg)
		case <-time.After(3 * time.Second):
			fmt.Println("timeout")
		}
	}

}
