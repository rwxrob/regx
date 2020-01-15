package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

const usage = `
usage: regx REGX STRING
       ... | regx REGX
       regx (looks for regx.txt)
`

const yes = `
_____.___.             ._.
\__  |   | ____   _____| |
 /   |   |/ __ \ /  ___/ |
 \____   \  ___/ \___ \ \|
 / ______|\___  >____  >__
 \/           \/     \/ \/
`

const no = `
  _________                           
 /   _____/ __________________ ___.__.
 \_____  \ /  _ \_  __ \_  __ <   |  |
 /        (  <_> )  | \/|  | \/\___  |
/_______  /\____/|__|   |__|   / ____|
        \/                     \/     
`

func display(match bool) {
	ws := getwinsize()
	ycws := ws.Col - 26
	ncws := ws.Col - 38
	rws := ws.Row - 6
	if ycws > 0 {
		ycws /= 2
	}
	if ncws > 0 {
		ncws /= 2
	}
	if rws > 0 {
		rws /= 2
	}
	ypad := strings.Repeat(" ", int(ycws))
	npad := strings.Repeat(" ", int(ncws))
	rpad := strings.Repeat("\n", int(rws))

	fmt.Print(Clear + rpad)
	if match {
		buf := strings.Replace(yes, "\n", "\n"+ypad, -1)
		fmt.Print(Y + buf + X)
	} else {
		buf := strings.Replace(no, "\n", "\n"+npad, -1)
		fmt.Print(R + buf + X)
	}
}

type winsize struct {
	Row, Col       uint16
	Xpixel, Ypixel uint16
}

func main() {
	var regx, buf string
	argc := len(os.Args)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			fmt.Print(Clear)
			os.Exit(0)
		}
	}()
	if argc > 3 || argc < 2 {
		last := false
		display(last)
		for {
			file, err := ioutil.ReadFile("regx.txt")
			if err != nil {
				fmt.Print(Clear)
				fmt.Println(usage)
				os.Exit(1)
			}
			i := strings.Index(string(file), "\n")
			if i < 0 {
				display(false)
				last = false
			} else {
				regx = string(file[:i])
				buf = string(file[len(regx)+1:])
				matches, _ := regexp.MatchString(regx, buf)
				if matches != last {
					display(matches)
				}
				last = matches
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	if argc == 3 {
		regx = os.Args[1]
		buf = os.Args[2]
	} else if argc == 2 {
		regx = os.Args[1]
		b, _ := ioutil.ReadAll(os.Stdin)
		buf = string(b)
	}
	matches, _ := regexp.MatchString(regx, buf)
	fmt.Println(matches)
}
