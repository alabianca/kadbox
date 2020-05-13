package kadprotocol

import (
	"bufio"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/alabianca/kadbox/log"
	"github.com/libp2p/go-libp2p-core/network"
	"io"
	"os"
	"path"
)

const (
	Message_WantType int = 1
	Message_AckType  int = 2
	Message_NackType int = 3
)

func New() core.KadProtocolService {
	return &kadprotocolService{}
}

type kadprotocolService struct {
}

func (k *kadprotocolService) HandleStream(s network.Stream) core.KadProtocol {
	fmt.Println("Handling the stream")
	kp := kadprotocol{
		want: make(chan []byte),
		ack:  make(chan string),
		nack: make(chan string),
		errc: make(chan error),
	}
	outr, outw := io.Pipe()

	kp.outr = outr
	kp.outw = outw

	go kp.read(s)
	go kp.write(s)

	return &kp
}

type kadprotocol struct {
	want chan []byte
	ack  chan string
	nack chan string
	errc chan error
	outr io.ReadCloser
	outw io.WriteCloser
}

func (k *kadprotocol) read(s network.Stream) {
	reader := bufio.NewReader(s)
	defer k.outw.Close()


	for {
		// 1. read the first byte which indicates to us which operation has to be performed
		b, err := reader.ReadByte()
		if err != nil {
			return
		}

		// 2. read the body based on the operation
		log.Debugf("Handling Message Type: %d\n", int(b))
		switch int(b) {
		case Message_WantType:
			// read the 64 byte key
			buf := make([]byte, 64)
			_, err := reader.Read(buf)
			if err != nil {
				return
			}
			// must find the file they are looking for
			k.ack <- string(buf)

		case Message_AckType:
			// we must read the content now which can have a variable length.
			// we can read till EOF
			io.Copy(k.outw, reader)
			return

		default:
			log.Debug("Unknown operation")

		}
	}
}

func (k *kadprotocol) write(s network.Stream) {
	defer s.Close()

	writer := bufio.NewWriter(s)

	for {
		select {
		case key := <-k.want:
			// send a 'want' message type
			log.Debug("Sending Message Type Want")
			writer.Write(prependOpKey(Message_WantType, key))
			writer.Flush()
		case fileHash := <-k.ack:
			// send a 'ack' message type
			log.Debug("Sending Message Type ACK")
			p, err := core.GetClosestKadboxRepoRelativeToWd()
			if err != nil {
				log.Debugf("Error %s\n", err)
				return
			}

			file, err := os.Open(path.Join(p, "store", fileHash))
			if err != nil {
				log.Debugf("Error %s\n", err)
				return
			}

			writer.Write([]byte{byte(Message_AckType)})
			if _, err := io.Copy(writer, file); err != nil {
				log.Debugf("Error %s\n", err)
				return
			}

			writer.Flush()
			return
		}
	}
}

func (k *kadprotocol) Want(key string) (io.Reader, chan error) {

	go func() {
		key := []byte(key)
		msg := append([]byte{}, key...)

		k.want <- msg
	}()

	return k.outr, nil
}

func prependOpKey(op int, msg []byte) []byte {
	return append([]byte{byte(op)}, msg...)
}
