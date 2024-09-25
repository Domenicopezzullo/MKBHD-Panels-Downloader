package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// Wallpaper represents the structure of the wallpaper data
type Wallpaper struct {
	DHD string `json:"dhd"`
	DSD string `json:"dsd"`
}

// Response represents the overall response structure
type Response struct {
	Version int                       `json:"version"`
	Data    map[string]Wallpaper `json:"data"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "<wallpaperid>",
		Short: "Downloads Wallpaper from MKBHD Overpriced App!",
		Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
		Run: func(cmd *cobra.Command, args []string) {
			wallpaperID := args[0]
			downloadWallpaper(wallpaperID)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}

func downloadWallpaper(wallpaperID string) {
	url := "https://storage.googleapis.com/panels-api/data/20240916/media-1a-i-p~s" // Example URL

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

	wallpaper, exists := apiResponse.Data[wallpaperID]
	if !exists {
		fmt.Printf("Wallpaper ID %s does not exist.\n", wallpaperID)
		return
	}

	// Save the wallpaper URLs to disk
	saveWallpaper(wallpaperID, wallpaper.DHD, "dhd")
	saveWallpaper(wallpaperID, wallpaper.DSD, "dsd")
}

func saveWallpaper(wallpaperID, url, quality string) {
	// Create the directory if it doesn't exist
	err := os.MkdirAll("wallpapers", os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Create a file to save the image
	fileName := fmt.Sprintf("wallpapers/%s_%s.jpg", wallpaperID, quality)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	// Fetch the image and write to the file
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

	fmt.Printf("Saved wallpaper ID %s (%s) to %s\n", wallpaperID, quality, fileName)
}

