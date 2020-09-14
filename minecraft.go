// Check the minecraft server download site to see if there's a new version
// then start the minecraft server with default recommended args

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const remoteDownloadPage = "https://www.minecraft.net/en-us/download/server"
const pageParseRegex = `<a href="(.*server.jar)" aria-label="mincraft version">(minecraft_server[0-9\.]*jar)</a>`
const archiveLocation = "versionArchives"

var javaArgs1 = []string{
	"-Xmx1024M",
	"-Xms1024M",
	"-jar",
}
var javaArgs2 = []string{
	"nogui",
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("encountered fatal error=%v\n", err)
		} else {
			fmt.Printf("exiting normally")
		}

		fmt.Println()
		fmt.Println("press any key to exit...")
		fmt.Scanf("h")
	}()

	minecraft()
}

func minecraft() {
	fmt.Println("getting local version...")
	localVersion, err := getLocalVersion()
	if err != nil {
		panic("could not get local version")
	}
	fmt.Printf("localVersion=%v\n", localVersion)

	fmt.Println("getting remote version...")
	remoteVersion, remoteHREF, err := getRemoteVersion()
	if err != nil {
		panic("could not get remote version")
	}
	fmt.Printf("remoteVersion=%v remoteHREF=%v\n", remoteVersion, remoteHREF)

	if localVersion != remoteVersion {
		fmt.Println("updating local version...")
		err = updateLocalVersion(remoteVersion, remoteHREF)
		if err != nil {
			panic("could not update local version")
		}
		fmt.Println("finished updating local version")

		if localVersion != "" {
			fmt.Println("archiving old version...")
			err = archiveOldVersion(localVersion)
			if err != nil {
				fmt.Printf("could not archive old version err=%v\n", err)
			}
			fmt.Println("finished archiving old version")
		}

		localVersion = remoteVersion
	} else {
		fmt.Println("local version is up to date.")
	}

	err = agreeToEULA()
	if err != nil {
		panic("could not agree to eula")
	}

	fmt.Println("starting server...")
	javaArgs := append(javaArgs1, localVersion)
	javaArgs = append(javaArgs, javaArgs2...)
	javaCmd := exec.Command("java", javaArgs...)
	javaCmd.Stdin = os.Stdin
	javaCmd.Stdout = os.Stdout
	javaCmd.Stderr = os.Stderr
	err = javaCmd.Run()
	if err != nil {
		panic(err)
	}
}

func getLocalVersion() (jarFile string, err error) {
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, archiveLocation) {
			return nil
		}

		if filepath.Ext(path) == ".jar" {
			jarFile = info.Name()
		}

		return nil
	})
	return
}

func getRemoteVersion() (jarFile, HREF string, err error) {
	resp, err := http.Get(remoteDownloadPage)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	re := regexp.MustCompile(pageParseRegex)

	submatches := re.FindSubmatch(body)
	// 0 is the full match
	HREF = string(submatches[1])
	jarFile = string(submatches[2])

	return
}

func updateLocalVersion(remoteVersion, remoteURL string) (err error) {
	resp, err := http.Get(remoteURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(remoteVersion)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return
}

func archiveOldVersion(oldJarFile string) (err error) {
	err = os.Rename(oldJarFile, fmt.Sprintf("%v/%v", archiveLocation, oldJarFile))
	return
}

func agreeToEULA() (err error) {
	out, err := os.Create("eula.txt")
	if err != nil {
		return
	}

	_, err = out.WriteString("eula=true")
	return
}
