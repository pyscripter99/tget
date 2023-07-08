/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/minio/selfupdate"
	"github.com/pyscripter99/tget/internals"
	"github.com/spf13/cobra"
)

var version string = "UNSET"

type GithubReleaseAsset struct {
	Name                 string `json:"name"`
	Size                 int    `json:"size"`
	Browser_download_url string `json:"browser_download_url"`
}

type GithubRelease struct {
	Html_url     string               `json:"html_url"`
	Tag_name     string               `json:"tag_name"`
	Published_at string               `json:"published_at"`
	Assets       []GithubReleaseAsset `json:"assets"`
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Downloads latest version of tget from GitHub",
	Long:  `Downloads latest version of tget from GitHub`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("https://api.github.com/repos/pyscripter99/tget/releases/latest")
		if err != nil {
			internals.Fatal(err)
		}

		respText, err := io.ReadAll(resp.Body)
		if err != nil {
			internals.Fatal(err)
		}

		var release GithubRelease
		if err := json.Unmarshal(respText, &release); err != nil {
			internals.Fatal(err)
		}

		fmt.Println("Latest: " + release.Tag_name)
		fmt.Println("Current: v" + version)
		if release.Tag_name == "v"+version {
			internals.Fatal(fmt.Errorf("already up to date"))
		}
		published, err := time.Parse(time.RFC3339, release.Published_at)
		if err != nil {
			internals.Fatal(err)
		}
		fmt.Println("Published: " + published.Format(time.RFC822))

		for _, artifact := range release.Assets {
			artifact.Name = strings.ToLower(artifact.Name)
			goarch := runtime.GOARCH
			goarch = strings.ReplaceAll(goarch, "amd64", "x86")
			goarch = strings.ReplaceAll(goarch, "386", "i386")
			if strings.HasPrefix(artifact.Name, "tget_"+runtime.GOOS+"_"+goarch) {
				fmt.Println("Downloading: " + artifact.Name)
				resp, err := http.Get(artifact.Browser_download_url)
				if err != nil {
					internals.Fatal(err)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					internals.Fatal(err)
				}

				if strings.HasSuffix(artifact.Name, ".zip") {

					zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
					if err != nil {
						internals.Fatal(err)
					}

					for _, file := range zipReader.File {
						if file.Name == "tget" || file.Name == "tget.exe" {
							f, err := file.Open()
							if err != nil {
								internals.Fatal(err)
							}
							defer f.Close()
							if err := selfupdate.Apply(resp.Body, selfupdate.Options{}); err != nil {
								internals.Fatal(err)
							}
							break
						}
					}
				} else if strings.HasSuffix(artifact.Name, ".tar.gz") {
					uncompressed, err := gzip.NewReader(bytes.NewReader(body))
					if err != nil {
						internals.Fatal(err)
					}

					tarReader := tar.NewReader(uncompressed)

					for {
						header, err := tarReader.Next()
						if err == io.EOF {
							break
						}

						if err != nil {
							internals.Fatal(err)
						}

						if header.Name == "tget" || header.Name == "tget.exe" {
							if err := selfupdate.Apply(tarReader, selfupdate.Options{}); err != nil {
								internals.Fatal(err)
							}
							break
						}
					}
				} else {
					internals.Fatal(fmt.Errorf("unknown archive format: '%s'", strings.Join(strings.Split(artifact.Name, ".")[1:], ".")))
				}
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
