package apps

import (
	"encoding/json"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"octiiCli/utils"
	"os"
	"sync"
)

var (
	types = map[int]string{
		1: "Theme",
		2: "Client Integration",
		3: "Server Integration",
	}
)

type Resource struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type int `json:"type"`
}

func getResource(ch chan interface{}, wg *sync.WaitGroup, productID string, id string, token string) {
	defer wg.Done()
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://gateway.octii.chat/products/" + productID + "/resources/" + id, nil)

	if err != nil {
		ch <- &errorString{"Couldn't reach gateway!"}
		return
	}

	request.Header.Set("Authorization", token)

	response, err := client.Do(request)

	if err != nil {
		ch <- &errorString{"Couldn't reach gateway!"}
		return
	}

	if response.StatusCode != 200 {
		ch <- &errorString{"Request to gateway failed!"}
		return
	}

	var res Resource
	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		ch <- &errorString{"Couldn't decode gateway response"}
		return
	}

	ch <- res
}

func Resources(c *cli.Context) error {
	id := c.Args().First()
	token, err := utils.Authorization()
	if err != nil {
		println(errorStyle.Render("You aren't authenticated! To login please use 'octii account login'"))
		return nil
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://gateway.octii.chat/products/" + id + "/resources", nil)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't reach gateway!")}
	}

	request.Header.Set("Authorization", token)

	response, err := client.Do(request)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't reach gateway!")}

	}

	if response.StatusCode != 200 {
		return &errorString{errorStyle.Render("Request to gateway failed!")}
	}

	var res []string

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return &errorString{errorStyle.Render("Couldn't decode gateway response!")}
	}

	var wg sync.WaitGroup
	var results []Resource
	ch := make(chan interface{})
	collectorCh := make(chan error)

	for _, resourceID := range res {
		wg.Add(1)
		go getResource(ch, &wg, id, resourceID, token)
	}

	go func() {
		for res := range ch {
			switch v := res.(type) {
			case error:
				collectorCh <- v
				return
			case Resource:
				results = append(results, v)
			}
		}
		collectorCh <- nil
	}()

	wg.Wait()
	close(ch)
	err = <- collectorCh

	if err != nil {
		return &errorString{errorStyle.Render(err.Error())}

	}

	info := make(map[string]string)

	for _, v := range results {
		info[keyStyle.Render(v.ID) + " " + validStyle.Render(types[v.Type])] = dataStyle.Render(v.Name)
	}

	println(utils.List(info, 75))

	return nil
}

func GetResource(c *cli.Context) error {
	productId := c.Args().Get(0)
	resourceId := c.Args().Get(1)
	outputFile := c.Args().Get(2)

	token, err := utils.Authorization()
	if err != nil {
		return &errorString{errorStyle.Render("You aren't authenticated! To login please use 'octii account login'")}
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://gateway.octii.chat/products/" + productId + "/resources/" + resourceId + "/payload", nil)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't reach gateway!")}
	}

	request.Header.Set("Authorization", token)

	response, err := client.Do(request)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't reach gateway!")}
	}

	if response.StatusCode != 200 {
		return &errorString{errorStyle.Render("Request to gateway failed!")}
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return &errorString{errorStyle.Render("Couldn't open file!")}
	}

	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return &errorString{errorStyle.Render("Couldn't write file!")}
	}

	return nil
}

func PutResource(c *cli.Context) error {
	productId := c.Args().Get(0)
	resourceId := c.Args().Get(1)
	inputFile := c.Args().Get(2)

	token, err := utils.Authorization()
	if err != nil {
		return &errorString{errorStyle.Render("You aren't authenticated! To login please use 'octii account login'")}
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return &errorString{errorStyle.Render("Couldn't open file!")}
	}

	defer file.Close()

	client := &http.Client{}
	request, err := http.NewRequest("PUT", "https://gateway.octii.chat/products/" + productId + "/resources/" + resourceId + "/payload", file)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't reach gateway!")}
	}

	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't reach gateway!")}
	}

	if response.StatusCode != 200 {
		return &errorString{errorStyle.Render("Request to gateway failed!")}
	}

	return nil
}