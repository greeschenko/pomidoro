package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

const Period = 2 // work or rest period in minuts
const Break = 5  // after period break time in minuts

var Needwrite int = 0

type Workoverlord struct {
	TIMER           int //start timer timestamp
	PAUSE           int //turn pause timestamp
	TIMERMODE       int //current timer mode 1-work mode 2-restmode
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
	DAYSTART        int
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
		Needwrite = 1
	}
}

func (w *Workoverlord) Write() {
	j, err := json.Marshal(w)
	if err != nil {
		fmt.Println("error1:", err)
	}

	w.Userfile.Truncate(0)

	if _, err := w.Userfile.Write(j); err != nil {
		fmt.Println("error2:", err)
	}
}

func (w *Workoverlord) Read() {
	file, _ := ioutil.ReadFile(w.Userfilepath)
	err := json.Unmarshal(file, &w)
	if err != nil {
		fmt.Println("error3:", err)
	}
}

func (w *Workoverlord) ShowHelp() {
	fmt.Println("")
	fmt.Println("Workoverlord - time and effectivity tracker app with gamification")
	fmt.Println("")
	fmt.Println("use:")
	fmt.Println("  workoverlord [command]")
	fmt.Println("------------")
	fmt.Println("commands:")
	fmt.Println("  status - return inline status")
	fmt.Println("  work   - start timer for work process (reduce duty points)")
	fmt.Println("  rest   - start timer for rest activity (require action points)")
}

func (w *Workoverlord) GetStatus() {
	t := ""

	if w.DAYSTART == 0 || int64(w.DAYSTART+60*60*24) < time.Now().Unix() {
		w.DAYSTART = int(time.Now().Unix())
		w.DP = w.DP - 6
		Needwrite = 1
	}

	if w.TIMERMODE == 1 {
		t = "W|"
	} else if w.TIMERMODE == 2 {
		t = "R|"
	}

	res := fmt.Sprintf("%d %s%s %d", w.DP, t, w.GetTime(), w.AP)
	os.Stdout.Write([]byte(res))
}

func (w *Workoverlord) GetTime() string {
	var res string
	now := time.Now().Unix()
	periodseconds := Period * 60
	periodend := w.TIMER + periodseconds
	delta := int(now) - periodend
	if delta < 0 {
		res = secondsToMinutes(delta)
	} else {
		w.DonePeriod()
		res = "DONE"
	}
	return res
}

func (w *Workoverlord) StartWorkTimer() {
	now := time.Now().Unix()
	w.TIMER = int(now)
	w.TIMERMODE = 1
	Needwrite = 1
}

func (w *Workoverlord) StartRestTimer() {
	now := time.Now().Unix()
	w.TIMER = int(now)
	w.TIMERMODE = 2
	Needwrite = 1
}

func (w *Workoverlord) DonePeriod() {
	if w.TIMERMODE == 1 {
		w.DP++
		if w.DP > 0 {
			w.AP++
		}
		w.TIMERMODE = 0
		Needwrite = 1
	} else if w.TIMERMODE == 2 {
		w.AP--
		w.TIMERMODE = 0
		Needwrite = 1
	}
}

func secondsToMinutes(inSeconds int) string {
	minutes := -1 * inSeconds / 60
	seconds := -1 * inSeconds % 60
	str := fmt.Sprintf("%d:%d", minutes, seconds)
	return str
}

func main() {
	App := Workoverlord{}
	App.Init()
	App.Read()
	defer App.Userfile.Close()
	defer App.Userlog.Close()

	if len(os.Args) > 1 {
		if os.Args[1] == "status" {
			App.GetStatus()
		} else if os.Args[1] == "work" {
			App.StartWorkTimer()
		} else if os.Args[1] == "rest" {
			App.StartRestTimer()
		}
	} else {
		App.ShowHelp()
	}

	if Needwrite == 1 {
		App.Write()
	}
}
