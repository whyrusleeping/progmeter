package progmeter

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Key    string
	Name   string
	State  string
	Active bool
}

type ProgMeter struct {
	lk    sync.Mutex
	Items []Item
	tick  *time.Ticker
	start time.Time
}

const up = "\033[%dA"
const down = "\033[%dB"
const right = "\033[%dC"

const (
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	magenta   = "\033[35m"
	colorNorm = "\033[0m"
)

func color(c, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, colorNorm)
}

func rightPad(s string, w int) string {
	if len(s) < w {
		return s + strings.Repeat(" ", w-len(s))
	}
	return s
}

func NewProgMeter() *ProgMeter {
	pm := &ProgMeter{
		tick:  time.NewTicker(time.Second),
		start: time.Now(),
	}
	go pm.run()
	return pm
}

func (p *ProgMeter) run() {
	for range p.tick.C {
		fmt.Printf("\r%s", time.Since(p.start))
	}
}

func (p *ProgMeter) Stop() {

}

func (p *ProgMeter) AddEntry(key, name, inf string) {
	p.lk.Lock()
	defer p.lk.Unlock()
	p.Items = append(p.Items, Item{Key: key, Name: name, Active: true, State: "actv"})
	fmt.Printf("\r[%s] %s%s", color(yellow, "get "), rightPad(name, 40), inf)
	fmt.Println()
}

func (p *ProgMeter) SetState(key, state string) {
	p.lk.Lock()
	defer p.lk.Unlock()
	for i := 1; i <= len(p.Items); i++ {
		if p.Items[len(p.Items)-i].Key == key {
			fmt.Printf(up, i)
			fmt.Printf("\r[%s]", state)
			fmt.Printf(down+"\r", i)
			break
		}
	}
}

func (p *ProgMeter) Finish(key string) {
	p.SetState(key, color(green, "done"))
}

func (p *ProgMeter) Error(key, err string) {
	p.lk.Lock()
	defer p.lk.Unlock()
	for i := 1; i <= len(p.Items); i++ {
		it := p.Items[len(p.Items)-i]
		if it.Key == key {
			fmt.Printf(up, i)
			fmt.Printf("\r[%s]", color(red, "err "))
			fmt.Printf(" %s (%s)", it.Name, err)
			fmt.Printf(down+"\r", i)
			break
		}
	}
}

func (p *ProgMeter) Working(key, state string) {
	p.SetState(key, color(magenta, state))
}
