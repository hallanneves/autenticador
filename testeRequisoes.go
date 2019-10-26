package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type tarefa struct {
	numero  int
	usuario string
	senha   string
}

var fila = make(chan tarefa, 1000)

//Payload é a struct no formato do json da requisição
type Payload struct {
	Senha   string `json:"senha"`
	Usuario string `json:"usuario"`
}

func requisitor() {
	fmt.Println("requisitor inicado")
	for {
		tarefa := <-fila
		//tarefaString := strconv.Itoa(tarefa.numero)
		//fmt.Println("Requisitor " + tarefaString + " start")

		data := Payload{Senha: tarefa.usuario, Usuario: tarefa.senha}
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", "http://127.0.0.1:8080/v1/auth", body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//fmt.Println("Requisitor " + tarefaString + ", resposta " + resp.Status)

		defer resp.Body.Close()

	}
}

func main() {

	fmt.Println("Iniciando a fila")
	for contador := 0; contador < 199; contador++ {
		fila <- tarefa{numero: contador, usuario: "senha1", senha: "senha1"}
	}

	for contador := 200; contador < 399; contador++ {
		fila <- tarefa{numero: contador, usuario: "senha7", senha: "senha1"}
	}

	for contador := 400; contador < 599; contador++ {
		fila <- tarefa{numero: contador, usuario: "senha3", senha: "senha4"}
	}

	for contador := 600; contador < 799; contador++ {
		fila <- tarefa{numero: contador, usuario: "senha" + strconv.Itoa(contador), senha: "senha1" + strconv.Itoa(contador+1)}
	}

	fmt.Println("Iniciando requisitores")
	go requisitor()
	go requisitor()
	go requisitor()
	go requisitor()

	for len(fila) > 0 {
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Fim de teste")

}
