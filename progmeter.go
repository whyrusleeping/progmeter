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
	Info   string
	Start  time.Time
}

type ProgMeter struct {
	lk    sync.Mutex
	Items []Item
	tick  *time.Ticker
	start time.Time

	total int
	done  int
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
		p.progdisp()
	}
}

func (p *ProgMeter) progdisp() {
	fmt.Printf("\r[%d / %d] %ds", p.done, p.total, int(time.Since(p.start).Seconds()))
}

func (p *ProgMeter) Stop() {
	p.tick.Stop()
}

func (p *ProgMeter) AddEntry(key, name, inf string) {
	p.lk.Lock()
	defer p.lk.Unlock()
	it := Item{
		Key:    key,
		Name:   name,
		Info:   inf,
		Active: true,
		State:  "actv",
		Start:  time.Now(),
	}
	p.Items = append(p.Items, it)
	it.Print(color(yellow, "get "))
	fmt.Println()
	p.progdisp()
}

func (it *Item) Print(state string) {
	fmt.Printf("\r[%s] %s%s", state, rightPad(it.Name, 35), it.Info)
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
	p.lk.Lock()
	defer p.lk.Unlock()
	p.done++
	for i := 1; i <= len(p.Items); i++ {
		it := p.Items[len(p.Items)-i]
		if it.Key == key {
			fmt.Printf(up, i)
			now := time.Now().Round(time.Millisecond)
			before := it.Start.Round(time.Millisecond)
			dur := now.Sub(before)
			it.Info += " " + dur.String()
			it.Print(color(green, "done"))
			fmt.Printf(down+"\r", i)
			break
		}
	}
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

func (p *ProgMeter) MarkDone() {
	p.lk.Lock()
	defer p.lk.Unlock()
	p.done++
}

func (p *ProgMeter) AddTodos(n int) {
	p.lk.Lock()
	defer p.lk.Unlock()
	p.total += n
}
