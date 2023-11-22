package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	path2 "path"
	"path/filepath"
	"strings"
)

type Resp struct {
	Arguments struct {
		Torrents []struct {
			ID          int    `json:"id"`
			DownloadDir string `json:"downloadDir"`
			Files       []struct {
				BytesCompleted int64  `json:"bytesCompleted"`
				Length         int64  `json:"length"`
				Name           string `json:"name"`
			} `json:"files"`
		} `json:"torrents"`
	} `json:"arguments"`
	Result string `json:"result"`
}

var pPath = flag.String("path", "", "Torrent path")
var pSession = flag.String("session", "", "Transmission session id")

func main() {

	flag.Parse()
	var path = *pPath
	var session = *pSession

	var files = make(map[string]int64)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files[path] = info.Size()
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	url := "http://192.168.0.10:8400/transmission/rpc"
	method := "POST"

	payload := strings.NewReader(`{"method":"torrent-get","arguments":{"fields":["id","files","downloadDir"]}}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "json")
	req.Header.Add("Cookie", "_session=MTY5OTAwMzEyNnxEdi1CQkFFQ180SUFBUkFCRUFBQVR2LUNBQUVHYzNSeWFXNW5EQkFBRG5KbFkyVnVkR3g1UldScGRHVmtCbk4wY21sdVp3d29BQ1pUWVhSMWNtNXBibVZUZEdGbmZIeDhVbVZzWVhSbFpGSm9hVzV2WTJWeWIzTjhmSHhzY3c9PXyqEEPCOpCDN018l_D6odiXlXPDXvjVKSUxCkUDjVWpcQ==; agh_session=b8aed7a449883a4c1bde830d46ce0fe1; pma_lang=zh_CN; pmaUser-1=algZ1zUjCJqxZ%2B7jTIZwkSPWrpOcVoLLLZxoEJwbPoHs94imHfkO1KJPds4%3D; firstVisit=1699436185528; compact_display_state=false; immich_access_token=A8PUlZw6dBNP6st75ri7p5XIt52gLCPiDTTm9Jsx4YQ; immich_auth_type=password; filter=all")
	req.Header.Add("Origin", "http://192.168.0.10:8400")
	req.Header.Add("Referer", "http://192.168.0.10:8400/transmission/web/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-Transmission-Session-Id", session)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resp Resp
	json.Unmarshal(body, &resp)

	var torrents []string

	for _, torrent := range resp.Arguments.Torrents {
		for _, file := range torrent.Files {
			torrents = append(torrents, path2.Join(torrent.DownloadDir, file.Name))
		}
	}
	for _, torrent := range torrents {
		delete(files, torrent)
	}
	for v, s := range files {
		fmt.Println(v, s)
	}
}
