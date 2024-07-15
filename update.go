package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/ulikunitz/xz"
)

type Version struct {
	Version     string `json:"version"`
	DownloadUrl string `json:"download_url"`
}

func getLocalVersion() (string, error) {
	versionFile := path.Join("puppilot", "version.json")
	// read version file
	data, err := os.ReadFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			// brand new installation
			return "", nil
		}
		return "", err
	}
	var version Version
	if err := json.Unmarshal(data, &version); err != nil {
		return "", err
	}
	return version.Version, nil
}

func getLatestVersion() (string, string, error) {
	res, err := http.Get(versionUrl)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	var version Version
	if err := json.Unmarshal(data, &version); err != nil {
		return "", "", err
	}
	return version.Version, version.DownloadUrl, nil
}

func checkForUpdate() error {
	fmt.Println("Checking for updates")
	// get local puppilot version
	localVersion, err := getLocalVersion()
	if err != nil {
		return err
	}
	// get latest puppilot version
	latestVersion, downloadUrl, err := getLatestVersion()
	if err != nil {
		fmt.Println("Failed to get latest version")
		fmt.Println(err.Error())
		if localVersion != "" {
			fmt.Println("Continue with current version")
			return nil
		} else {
			return err
		}
	}
	// compare versions
	if localVersion == latestVersion {
		fmt.Println("Puppilot is up to date")
		return nil
	}
	fmt.Println("Downloading update")
	if err := downloadUpdate(downloadUrl); err != nil {
		return err
	}
	fmt.Println("Update downloaded")
	// write latest version to version file
	versionFile := path.Join("puppilot", "version.json")
	version := Version{Version: latestVersion}
	data, err := json.Marshal(version)
	if err != nil {
		return err
	}
	if err := os.WriteFile(versionFile, data, 0644); err != nil {
		return err
	}
	return nil
}

func downloadUpdate(downloadUrl string) error {
	// download update
	res, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	uncompressed, err := xz.NewReader(res.Body)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressed)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(path.Join("puppilot", header.Name), 0755); err != nil {
				return err
			}
			continue
		case tar.TypeReg:
			file, err := os.Create(path.Join("puppilot", header.Name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
			file.Close()
			continue
		}
	}
	return nil
}
