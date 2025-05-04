package fileops

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
)

func NewAudioWriter(fileName string) (*os.File, func(int)) {
	if !strings.HasSuffix(fileName, ".aiff") {
		fileName += ".aiff"
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("error while creating audio file: %v", err)
	}

	// form chunk
	_, err = f.WriteString("FORM")
	if err != nil {
		log.Fatalf("error while writing FORM chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int32(0))
	if err != nil {
		log.Fatalf("error while writing size/placeholder for FORM chunk: %v", err)
	}

	_, err = f.WriteString("AIFF")
	if err != nil {
		log.Fatalf("error while writing AIFF chunk: %v", err)
	}

	// common chunk
	_, err = f.WriteString("COMM")
	if err != nil {
		log.Fatalf("error while writing COMM chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int32(18))
	if err != nil {
		log.Fatalf("error while writing size/placeholder for COMM chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int16(1))
	if err != nil {
		log.Fatalf("error while writing channels for COMM chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int32(0))
	if err != nil {
		log.Fatalf("error while writing number of samples for COMM chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int16(32))
	if err != nil {
		log.Fatalf("error while writing bits per sample for COMM chunk: %v", err)
	}

	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0})
	if err != nil {
		log.Fatalf("error while writing 80-bit sample rate for COMM chunk: %v", err)
	}

	// sound chunk
	_, err = f.WriteString("SSND")
	if err != nil {
		log.Fatalf("error while writing SSND chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int32(0))
	if err != nil {
		log.Fatalf("error while writing size for SSND chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int32(0))
	if err != nil {
		log.Fatalf("error while writing offset for SSND chunk: %v", err)
	}

	err = binary.Write(f, binary.BigEndian, int32(0))
	if err != nil {
		log.Fatalf("error while writing block for SSND chunk: %v", err)
	}

	missingSamplesFiller := func(nSamples int) {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = f.Seek(4, 0)
		if err != nil {
			log.Fatalf("error while seeking - 1: %v", err)
		}

		err = binary.Write(f, binary.BigEndian, int32(totalBytes))
		if err != nil {
			log.Fatalf("error while filling samples - 1: %v", err)
		}

		_, err = f.Seek(22, 0)
		if err != nil {
			log.Fatalf("error while seeking - 2: %v", err)
		}

		err = binary.Write(f, binary.BigEndian, int32(nSamples))
		if err != nil {
			log.Fatalf("error while filling samples - 2: %v", err)
		}

		_, err = f.Seek(42, 0)
		if err != nil {
			log.Fatalf("error while seeking - 3: %v", err)
		}

		err = binary.Write(f, binary.BigEndian, int32(4*nSamples+8))
		if err != nil {
			log.Fatalf("error while filling samples - 3: %v", err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("error while closing file: %v", err)
		}
	}

	return f, missingSamplesFiller
}
