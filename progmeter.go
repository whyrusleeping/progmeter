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

	dorun *sync.Once

	total   int
	done    int
	minimal bool
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

func (p *ProgMeter) color(c, s string) string {
	if p.minimal {
		return s
	}
	return fmt.Sprintf("%s%s%s", c, s, colorNorm)
}

func rightPad(s string, w int) string {
	if len(s) < w {
		return s + strings.Repeat(" ", w-len(s))
	}
	return s
}

func NewProgMeter(minimal bool) *ProgMeter {
	pm := &ProgMeter{
		tick:    time.NewTicker(time.Second),
		start:   time.Now(),
		minimal: minimal,
		dorun:   new(sync.Once),
	}
	return pm
}

func (p *ProgMeter) run() {
	if p.minimal {
		return
	}
	go func() {
		for range p.tick.C {
			p.progdisp()
		}
	}()
}

func (p *ProgMeter) progdisp() {
	fmt.Printf("\r[%d / %d] %ds", p.done, p.total, int(time.Since(p.start).Seconds()))
}

func (p *ProgMeter) Stop() {
	if p == nil {
		return
	}
	p.tick.Stop()
}

func (p *ProgMeter) AddEntry(key, name, inf string) {
	if p == nil {
		return
	}
	p.AddEntryWithState("get ", key, name, inf)
}

func (p *ProgMeter) AddEntryWithState(state, key, name, inf string) {
	if p == nil {
		return
	}
	p.dorun.Do(p.run)
	p.lk.Lock()
	defer p.lk.Unlock()
	it := Item{
		Key:    key,
		Name:   name,
		Info:   inf,
		Active: true,
		State:  state,
		Start:  time.Now(),
	}
	p.Items = append(p.Items, it)
	it.Print(p.color(yellow, state))
	fmt.Println()
	if !p.minimal {
		p.progdisp()
	}
}

func (it *Item) Print(state string) {
	fmt.Printf("\r[%s] %s%s", state, rightPad(it.Name, 40), it.Info)
}

func (p *ProgMeter) SetState(key, state string) {
	if p == nil {
		return
	}
	p.lk.Lock()
	defer p.lk.Unlock()
	for i := 1; i <= len(p.Items); i++ {
		it := p.Items[len(p.Items)-i]
		if it.Key == key {
			if p.minimal {
				it.Print(state)
			} else {
				fmt.Printf(up, i)
				fmt.Printf("\r[%s]", state)
				fmt.Printf(down+"\r", i)
			}
			break
		}
	}
}

func (p *ProgMeter) Finish(key string) {
	if p == nil {
		return
	}
	p.lk.Lock()
	defer p.lk.Unlock()
	p.done++
	for i := 1; i <= len(p.Items); i++ {
		it := p.Items[len(p.Items)-i]
		if it.Key == key {
			now := time.Now().Round(time.Millisecond)
			before := it.Start.Round(time.Millisecond)
			dur := now.Sub(before)
			it.Info += " " + dur.String()

			if !p.minimal {
				fmt.Printf(up, i)
			}

			it.Print(p.color(green, "done"))

			if !p.minimal {
				fmt.Printf(down+"\r", i)
			} else {
				fmt.Println()
			}
			break
		}
	}
}

func (p *ProgMeter) Error(key, err string) {
	if p == nil {
		return
	}
	p.lk.Lock()
	defer p.lk.Unlock()
	for i := 1; i <= len(p.Items); i++ {
		it := p.Items[len(p.Items)-i]
		if it.Key == key {
			it.Info += " " + err
			if !p.minimal {
				fmt.Printf(up, i)
			}

			it.Print(p.color(red, "err "))

			if !p.minimal {
				fmt.Printf(down+"\r", i)
			}
			break
		}
	}
}

func (p *ProgMeter) Working(key, state string) {
	if p == nil {
		return
	}
	p.SetState(key, p.color(magenta, state))
}

func (p *ProgMeter) MarkDone() {
	if p == nil {
		return
	}
	p.lk.Lock()
	defer p.lk.Unlock()
	p.done++
}

func (p *ProgMeter) AddTodos(n int) {
	if p == nil {
		return
	}
	p.lk.Lock()
	defer p.lk.Unlock()
	p.total += n
}
