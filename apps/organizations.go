package apps

import (
	"encoding/json"
	"github.com/innatical/octii-cli/utils"
	"github.com/urfave/cli/v2"
	"net/http"
	"sync"
)

func Organizations(c *cli.Context) error {
	token, err := utils.Authorization()
	if err != nil {
		return &errorString{errorStyle.Render("You aren't authenticated! To login please use 'octii account login'")}
	}

	claims, err := utils.ParseAuthorization(token)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't parse token! Please reauthenticate using 'octii account login'")}
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.octii.chat/v1/users/"+claims.Subject+"/organizations", nil)

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

	var res []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return &errorString{errorStyle.Render("Couldn't decode gateway response!")}
	}

	info := make(map[string]string)

	for _, v := range res {
		info[keyStyle.Render(v.ID)] = dataStyle.Render(v.Name)
	}

	println(utils.List(info, 75))
	return nil
}

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getProduct(ch chan interface{}, wg *sync.WaitGroup, id string, token string) {
	defer wg.Done()
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.octii.chat/v1/products/"+id, nil)

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

	var res Product
	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		ch <- &errorString{"Couldn't decode gateway response"}
		return
	}

	ch <- res
}

func Products(c *cli.Context) error {
	id := c.Args().First()
	token, err := utils.Authorization()
	if err != nil {
		return &errorString{errorStyle.Render("You aren't authenticated! To login please use 'octii account login'")}
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.octii.chat/v1/communities/"+id+"/products", nil)

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
	var results []Product
	ch := make(chan interface{})
	collectorCh := make(chan error)

	for _, productID := range res {
		wg.Add(1)
		go getProduct(ch, &wg, productID, token)
	}

	go func() {
		for res := range ch {
			switch v := res.(type) {
			case error:
				collectorCh <- v
				return
			case Product:
				results = append(results, v)
			}
		}
		collectorCh <- nil
	}()

	wg.Wait()
	close(ch)
	err = <-collectorCh

	if err != nil {
		return &errorString{errorStyle.Render(err.Error())}
	}

	info := make(map[string]string)

	for _, v := range results {
		info[keyStyle.Render(v.ID)] = dataStyle.Render(v.Name)
	}

	println(utils.List(info, 75))

	return nil
}
