package audio

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/gordonklaus/portaudio"
)

type AudioCapturer interface {
	StartRecording(f *os.File, stopCh <-chan bool, samplesCh chan<- int) error
	StopRecording(stopCh chan<- bool, samplesCh <-chan int, missingSamplesFiller func(int)) error
	ConvertToText()
}

type AudioCapture struct {
	stream *portaudio.Stream
}

func NewAudioCapture() *AudioCapture {
	portaudio.Initialize()

	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		log.Fatalf("failed to open default stream: %v", err)
	}

	return &AudioCapture{
		stream: stream,
	}
}

func (a *AudioCapture) StartRecording(
	f *os.File,
	stopCh <-chan bool,
	samplesCh chan<- int,
) error {
	fmt.Println("Starting recording...")

	in := make([]int32, 64)
	nSamples := 0

	err := a.stream.Start()
	if err != nil {
		return fmt.Errorf("failed to start stream: %v", err)
	}

	go func() {
		defer fmt.Println("go routine ended")

		for {
			err = a.stream.Read()
			if err != nil {
				log.Fatalf("error while reading from stream: %v", err)
			}

			err = binary.Write(f, binary.BigEndian, in)
			if err != nil {
				log.Fatalf("error while reading from stream: %v", err)
			}

			nSamples += len(in)
			select {
			case <-stopCh:
				samplesCh <- nSamples
				return
			default:
			}
		}
	}()

	return nil
}

func (a *AudioCapture) StopRecording(
	stopCh chan<- bool,
	samplesCh <-chan int,
	missingSamplesFiller func(int),
) error {
	fmt.Println("Stopping recording...")

	defer portaudio.Terminate()
	defer a.stream.Close()

	go func(samplesCh <-chan int) {
		select {
		case nSamples := <-samplesCh:
			missingSamplesFiller(nSamples)
			return
		default:
		}
	}(samplesCh)

	stopCh <- true

	err := a.stream.Stop()
	if err != nil {
		return fmt.Errorf("error while stopping stream: %v", err)
	}

	fmt.Println("Stopping recording ended...")

	return nil
}

func (a *AudioCapture) ConvertToText() {

}
