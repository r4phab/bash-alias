package main

import (
    "os"
    "os/exec"
    "fmt"
    "strings"
    "path/filepath"
    "bytes"
    "runtime"
	"io/ioutil"
	"log"
	"io"
	"net/http"
)

func main() {
	fmt.Println("Start")

	fmt.Println("Install ZSH ...")
	directoryDelete(userHomeDir() + "/.zsh/zsh-autosuggestions")
	exe_cmd([]string{
    	"sh -c \"$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)\"",
    	"git clone https://github.com/zsh-users/zsh-autosuggestions "+userHomeDir()+"/.zsh/zsh-autosuggestions",
    })

    fmt.Println("Done ...")
    fmt.Println("Configure ZSH ...")

	os.MkdirAll(userHomeDir() + "/.oh-my-zsh/custom/themes/", os.ModePerm);
	fileDelete(userHomeDir() + "/.oh-my-zsh/custom/themes/r4phab.zsh-theme")
	downloadFile("https://gist.githubusercontent.com/r4phab/c3d6cda9861c3595c316f7de97b6c69b/raw/",userHomeDir() + "/.oh-my-zsh/custom/themes/r4phab.zsh-theme")

	replaceInFile(userHomeDir() + "/.zshrc", "robbyrussell", "r4phab")
    safeWriteToFile(userHomeDir() + "/.zshrc", "source ~/.zsh/zsh-autosuggestions/zsh-autosuggestions.zsh")

    fmt.Println("Done ...")
    fmt.Println("Install Git-Extras ...")

    exe_cmd([]string{
   	    "curl -sSL http://git.io/git-extras-setup | sudo bash /dev/stdin",
    })

    fmt.Println("Done ...")
    fmt.Println("Install Aliasses ...")

    safeWriteToFile(userHomeDir() + "/.zshrc", "function gcommit() { \n git add . \n git commit -a -m \"$1\" \n git push \n }")

    fmt.Println("Done ...")
    fmt.Println("Don't forget to source .zshrc !")
}

func directoryDelete(path string){
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func fileDelete(path string){
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var (
    output_path = filepath.Join("./output")
    bash_script = filepath.Join( "_script.sh" )
)

func checkError( e error){
    if e != nil {
        panic(e)
    }
}

func exe_cmd(cmds []string) {
    os.RemoveAll(output_path)
    err := os.MkdirAll( output_path, os.ModePerm|os.ModeDir )
    checkError(err)
    file, err := os.Create( filepath.Join(output_path, bash_script))
    checkError(err)
    defer file.Close()
    file.WriteString("#!/bin/sh\n")
    file.WriteString( strings.Join(cmds, "\n"))
    err = os.Chdir(output_path)
    checkError(err)
    out, err := exec.Command("sh", bash_script).Output()
    checkError(err)
    fmt.Println(string(out))
}

func safeWriteToFile(path string, text string) {
    replaceInFile(path, text, "")
    writeToFile(path, text)
}

func writeToFile(path string, text string) error {
      f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
      if err != nil {
              return err
      }
      defer f.Close()

      _, err = f.WriteString("\n"+text)
      if err != nil {
              return err
      }
      return nil
}

func userHomeDir() string {
    if runtime.GOOS == "windows" {
        home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
        if home == "" {
            home = os.Getenv("USERPROFILE")
        }
        return home
    }
    return os.Getenv("HOME")
}

func createFile(path string){
	newFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
}

func downloadFile(url string, path string){
	// Create output file
	newFile, err := os.Create(path)
	if err != nil {
	  log.Fatal(err)
	}
	defer newFile.Close()

	// HTTP GET request devdungeon.com 
	response, err := http.Get(url)
	defer response.Body.Close()

	// Write bytes from HTTP response to file.
	// response.Body satisfies the reader interface.
	// newFile satisfies the writer interface.
	// That allows us to use io.Copy which accepts
	// any type that implements reader and writer interface
	io.Copy(newFile, response.Body)
	log.Printf("Downloaded file : " + url + " -> " + path)
}

func replaceInFile(path string, replaced string, withreplace string) error {
		input, err := ioutil.ReadFile(path)
		if err != nil {
		     fmt.Println(err)
		     os.Exit(1)
		}

		output := bytes.Replace(input, []byte(replaced), []byte(withreplace), -1)

		if err = ioutil.WriteFile(path, output, 0666); err != nil {
		     fmt.Println(err)
		     os.Exit(1)
		}
		return nil
}
