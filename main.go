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

var domains = [2]string{"191.36.8.53:4040", "191.36.8.54:4040"}

//var domainsDatabase = [3]string{"191.36.8.53:6060", "191.36.8.54:6060", "tccmarina.sj.ifsc.edu.br/app"}

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
	requestURL := fmt.Sprintf("https://tccmarina.sj.ifsc.edu.br/app/update")
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

}

func openFile(filename string) *os.File {
	file, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	return file
}

func writeFile(file *os.File, str string) {
	if _, err := file.WriteString(str); err != nil {
		log.Println(err)
	}

}

func searchTerm(term string, file *os.File) {
	for _, domain := range domains {
		start := time.Now()
		requestURL := fmt.Sprintf("http://%s/search/%s", domain, term)
		res, err := http.Get(requestURL)
		elapsed := time.Since(start).Seconds()
		if err != nil {
			logs := fmt.Sprintf("%s, %s, %s, %d, %f, %s \n", start.Format(time.UnixDate), domain, term, 516, elapsed, err)
			writeFile(file, logs)
			continue
		}
		logs := fmt.Sprintf("%s, %s, %s, %d, %f \n", start.Format(time.UnixDate), domain, term, res.StatusCode, elapsed)
		writeFile(file, logs)
		log.Printf("Bucando item: %s, status code: %d, tempo de resposta: %f ", term, res.StatusCode, elapsed)

	}
}

func searchItens(queries []string, file *os.File) {
	writeFile(file, "Time, EndPoint, Termo, StatusCode, ResponseTime, Details \n")
	for _, query := range queries {
		searchTerm(query, file)
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
	file := openFile("logs.txt")
	var queries = getListFromFile("./frequent-terms")
	var updatedIds = getListFromFile("./updated-ids")
	var removedIds = getListFromFile("./deleted-ids")
	go updateItens(updatedIds)
	go removeItens(removedIds)
	searchItens(queries, file)
}
