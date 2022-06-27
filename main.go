package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type ItemUpdate struct {
	Id      string `json:"id"`
	Collumn string `json:"collumn"`
	Value   string `json:"value"`
}

type ItemDelete struct {
	Id string `json:"id"`
}

func wait(max int) {
	rand.Seed(time.Now().UnixNano())
	sleepTime := rand.Intn(max)
	time.Sleep(time.Duration(sleepTime) * time.Second)
}

func readData(fileName string) ([][]string, error) {

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	csvFile := csv.NewReader(file)

	if _, err := csvFile.Read(); err != nil {
		log.Fatal(err)
	}

	lines, err := csvFile.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	return lines, nil
}

func getListFromFile(path string) []string {
	var terms []string
	lines, err := readData(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, line := range lines {
		terms = append(terms, line[0])
	}
	return terms

}

func deleteProduct(id string) {
	log.Printf("Removendo item: %s", id)
	item := ItemDelete{id}
	jsonReq, err := json.Marshal(item)
	request := fmt.Sprintf("https://tccmarina.sj.ifsc.edu.br/app/delete/%s", id)
	req, err := http.NewRequest(http.MethodDelete, request, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
}

func updateProduct(id string) {
	log.Printf("Atualizando item: %s", id)
	item := ItemUpdate{id, "title", "Titulo atualizado"}
	jsonReq, err := json.Marshal(item)
	req, err := http.NewRequest(http.MethodPut, "https://tccmarina.sj.ifsc.edu.br/app/update", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
}

func searchTerm(s string) {
	start := time.Now()
	requestURL := fmt.Sprintf("https://tccmarina.sj.ifsc.edu.br/finder/search/%s", s)
	res, err := http.Get(requestURL)
	elapsed := time.Since(start).Seconds()
	if err != nil {
		fmt.Printf("Erro: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Bucando item: %s, status code: %d, tempo de resposta: %f ", s, res.StatusCode, elapsed)
}

func searchItens(queries []string) {
	for _, query := range queries {
		wait(1)
		searchTerm(query)
	}
}

func updateItens(ids []string) {
	for _, id := range ids {
		wait(5)
		updateProduct(id)
	}
}

func removeItens(ids []string) {
	for _, id := range ids {
		wait(20)
		deleteProduct(id)
	}
}

func main() {

	var queries = getListFromFile("./frequent-terms")
	var updatedIds = getListFromFile("./updated-ids")
	var removedIds = getListFromFile("./deleted-ids")
	go updateItens(updatedIds)
	go removeItens(removedIds)
	searchItens(queries)
}
