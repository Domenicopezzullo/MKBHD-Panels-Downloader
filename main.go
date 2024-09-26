package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Wallpaper struct {
	DHD string `json:"dhd,omitempty"`
}

type Response struct {
	Version int                `json:"version"`
	Data    map[string]Wallpaper `json:"data"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "<wallpapername>",
		Short: "Downloads Wallpaper from MKBHD Overpriced App!",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			wallpaperName := args[0]
			downloadWallpaper(wallpaperName)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}

func downloadWallpaper(wallpaperName string) {
	url := "https://storage.googleapis.com/panels-api/data/20240916/media-1a-i-p~s"

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching wallpaper data: %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch wallpaper data: %s\n", response.Status)
		return
	}

	var apiResponse Response
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	wallpaperNameLower := strings.ToLower(wallpaperName)

	for wallpaperID, wallpaper := range apiResponse.Data {
		if strings.Contains(strings.ToLower(wallpaper.DHD), wallpaperNameLower) {
			saveWallpaper(wallpaperID, wallpaper.DHD)
			return
		}
	}

	fmt.Printf("No wallpaper found with name %s in the DHD URL.\n", wallpaperName)
}

func saveWallpaper(wallpaperID, url string) {
	err := os.MkdirAll("wallpapers", os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	fileName := fmt.Sprintf("wallpapers/%s_dhd.jpg", wallpaperID)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	imageResponse, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching image for wallpaper ID %s: %v\n", wallpaperID, err)
		return
	}
	defer imageResponse.Body.Close()

	if _, err := io.Copy(file, imageResponse.Body); err != nil {
		fmt.Printf("Error saving image for wallpaper ID %s: %v\n", wallpaperID, err)
		return
	}

	fmt.Printf("Saved wallpaper ID %s (dhd) to %s\n", wallpaperID, fileName)
}

