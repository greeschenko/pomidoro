package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

type Workoverlord struct {
	TIMER           int
	LVL             int
	EXP             int
	NEXP            int
	DP              int
	AP              int
	T               int
	G               int
	NAP             int
	NT              int
	NG              int
	Userfile        *os.File `json:"-"`
	Userlog         *os.File `json:"-"`
	Userfilepath    string   `json:"-"`
	Userlogfilepath string   `json:"-"`
}

func (w *Workoverlord) Init() {
	var err error
	usr, _ := user.Current()
	homedir := usr.HomeDir

	mainpath := fmt.Sprintf("%s/.config/workoverlord", homedir)
	if _, err := os.Stat(mainpath); os.IsNotExist(err) {
		os.Mkdir(mainpath, os.ModePerm)
	}

	w.Userfilepath = fmt.Sprintf("%s/user.json", mainpath)
	w.Userfile, err = os.OpenFile(w.Userfilepath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	w.Userlogfilepath = fmt.Sprintf("%s/user.log", mainpath)
	w.Userlog, err = os.OpenFile(w.Userfilepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := w.Userfile.Stat()
	if err != nil {
		fmt.Println("error4:", err)
	}

	if fi.Size() < 5 {
		w.Write()
	}
}

func (w *Workoverlord) Write() {
	j, err := json.Marshal(w)
	if err != nil {
		fmt.Println("error1:", err)
	}

	if _, err := w.Userfile.WriteAt(j, 0); err != nil {
		fmt.Println("error2:", err)
	}

	os.Stdout.Write(j)
}

func (w *Workoverlord) Read() {
	file, _ := ioutil.ReadFile(w.Userfilepath)
	err := json.Unmarshal(file, &w)
	if err != nil {
		fmt.Println("error3:", err)
	}
}

func main() {
	App := Workoverlord{}
	App.Init()
	//App.Write()
	App.Read()

	fmt.Println(App.AP)
	fmt.Println(App.DP)

	App.DP = 2
	App.Write()

	App.Userfile.Close()
	App.Userlog.Close()
}
