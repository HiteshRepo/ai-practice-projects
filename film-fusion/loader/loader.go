package loader

import (
	"fmt"
	"sync"
	"time"
)

type Loader struct {
	message string
	frames  []string
	delay   time.Duration
	running bool
	mu      sync.Mutex
	done    chan bool
}

func NewLoader(message string) *Loader {
	return &Loader{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:   100 * time.Millisecond,
		done:    make(chan bool),
	}
}

// []string{"◐", "◓", "◑", "◒"}
func NewLoaderWithFrames(message string, frames []string, delay time.Duration) *Loader {
	return &Loader{
		message: message,
		frames:  frames,
		delay:   delay,
		done:    make(chan bool),
	}
}

func (l *Loader) Start() {
	l.mu.Lock()
	l.running = true
	l.mu.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-l.done:
				return
			default:
				l.mu.Lock()
				if !l.running {
					l.mu.Unlock()
					return
				}

				fmt.Printf("\r%s %s", l.frames[i%len(l.frames)], l.message)
				i++
				l.mu.Unlock()

				time.Sleep(l.delay)
			}
		}
	}()
}

func (l *Loader) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.running {
		l.running = false
		l.done <- true
		fmt.Print("\r")
	}
}

func (l *Loader) UpdateMessage(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.message = message
}

func (l *Loader) Complete(message string) {
	l.Stop()
	fmt.Printf("\n✓ %s\n", message)
}

func (l *Loader) Fail(message string) {
	l.Stop()
	fmt.Printf("\n✗ %s\n", message)
}
