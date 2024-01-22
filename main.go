package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CEPResponse struct {
	Provider string
	CEP      string
}

func getCep(url string, ch chan<- string, c *http.Client) {
	resp, err := c.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response, _ := io.ReadAll(resp.Body)
	ch <- string(response)
}

func getResult(via_cep_channel chan string, via_cep_cdn chan string) (*CEPResponse, error) {
	select {
	case msgViaCep := <-via_cep_channel:
		return &CEPResponse{Provider: "VIACEP", CEP: msgViaCep}, nil
	case msgCDN := <-via_cep_cdn:
		return &CEPResponse{Provider: "CDNCEP", CEP: msgCDN}, nil
	case <-time.After(time.Second):
		return nil, errors.New("endpoint request timeout")
	}
}

func main() {
	c := http.Client{}
	via_cep_channel := make(chan string)
	via_cep_cdn := make(chan string)

	go getCep("https://viacep.com.br/ws/63111-020/json/", via_cep_channel, &c)
	go getCep("https://cdn.apicep.com/file/apicep/63111-020.json", via_cep_cdn, &c)

	result, err := getResult(via_cep_channel, via_cep_cdn)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}
