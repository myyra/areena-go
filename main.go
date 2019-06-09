package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {

	appKey := "8930d72170e48303cf5f3867780d549b"

	baseURL, err := url.Parse("https://player.api.yle.fi")
	if err != nil {
		fmt.Println(err)
	}

	areenaID := os.Args[1]

	preview := getPreview(baseURL, appKey, areenaID)
	playlist := getPlaylist(preview.Data.OngoingOndemand.ManifestURL)
	getVideoFiles(playlist)
}

func getVideoFiles(playlistURL string) {

	videoFileRegex := regexp.MustCompile(`.*.ts`)

	fmt.Println("Getting list of files")
	resp, err := http.Get(playlistURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	playlist := string(body)
	fileNames := videoFileRegex.FindAllString(playlist, -1)
	for _, fileName := range fileNames {
		index := strings.Split(fileName, "_")[1]
		url := strings.ReplaceAll(playlistURL, ".m3u8", "_"+index)
		fmt.Println("Downloading file", fileName, "out of", len(fileNames))
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println(err)
			panic(err)
		}
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		out, err := os.Create(dir + "/" + fileName)
		if err != nil {
			panic(err)
		}
		io.Copy(out, resp.Body)
		err = out.Close()
		if err != nil {
			panic(err)
		}
	}
}

func getPlaylist(manifestURL string) string {

	fmt.Println("Getting playlist")

	manifestRegex := regexp.MustCompile(`.*index\d.*`)

	resp, err := http.Get(manifestURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	content := string(body)
	matches := manifestRegex.FindAllString(content, -1)
	bestResolutionPlaylistURL := matches[len(matches)-1]
	return bestResolutionPlaylistURL
}

type Preview struct {
	Data struct {
		OngoingOndemand struct {
			Description struct {
				Fin string `json:"fin,omitempty"`
			} `json:"description,omitempty"`

			Image struct {
				Id      string `json:"id,omitempty"`
				Version int    `json:"version,omitempty"`
			} `json:"image,omitempty"`

			ManifestURL string `json:"manifest_url,omitempty"`

			Title struct {
				Fin string `json:"fin,omitempty"`
			} `json:"title,omitempty"`
		} `json:"ongoing_ondemand,omitempty"`
	} `json:"data,omitempty"`
}

func getPreview(baseURL *url.URL, appKey, id string) Preview {

	fmt.Println("Getting stream info")

	base := *baseURL
	url := &base
	url.Path = fmt.Sprintf("/v1/preview/%s.json", id)
	q := url.Query()
	q.Set("app_key", appKey)
	q.Set("language", "fin")
	q.Set("app_id", "player_static_prod")
	url.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Origin", "https://areena.yle.fi")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var response Preview
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
	}

	return response
}
