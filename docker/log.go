package docker

import (
	"bytes"
	"errors"
	"strings"
)

type DockerLog struct {
	Type      string            `json:"Type"`
	Labels    map[string]string `json:"Labels"`
	Timestamp string            `json:"Timestamp"`
	IP        string            `json:"IP"`
	Message   string            `json:"Message"`
}

// ParseLog of docker log
func ParseLog(b []byte) (*DockerLog, error) {
	// It is encoded on the first 8 bytes like this:
	//
	// header := [8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}
	//
	// `STREAM_TYPE` can be:
	//
	// -   0: stdin (will be written on stdout)
	// -   1: stdout
	// -   2: stderr
	//
	// `SIZE1, SIZE2, SIZE3, SIZE4` are the 4 bytes of
	// the uint32 size encoded as big endian.

	h := make([]byte, 8)
	buf := bytes.NewBuffer(b)
	buf.Read(h)

	var t string

	switch h[0] {
	case 0:
		t = "stdin"
	case 1:
		t = "stdout"
	case 2:
		t = "stderr"
	}

	//xlog.Debugf("Docker Log: %s", string(buf.Bytes()))

	msg := buf.Bytes()

	if len(msg) == 0 {
		return nil, errors.New("Log line empty")
	}

	logmsg := bytes.SplitN(msg, []byte(" "), 3)

	//part := bytes.SplitN(logmsg[0], []byte(" "), 2)

	labels := make(map[string]string)

	l := string(logmsg[1])

	if l != "" {
		items := strings.Split(l, ",")

		for _, val := range items {
			pair := strings.SplitN(val, "=", 2)

			labels[pair[0]] = pair[1]
		}
	}

	// 2017-01-20T05:50:42.047552194Z  10.0.1.3 - - [20/Jan/2017:05:50:41 +0000] "GET /logo.png HTTP/1.1" 200 13133 "http://localhost:8012/" "Mozilla
	// 2017-01-20T05:50:42.047490838Z com.docker.swarm.node.id=8dsrkozezn9j2lbvmciwtn2wf,com.docker.swarm.service.id=3hwqpkilg4saotx96jeiwgggg,com.docker.swarm.task.id=z037zkyr99u649qimsumnho0w 10.0.1.3 - - [20/Jan/2017:05:50:41 +0000] "GET / HTTP/1.1" 200 485 "-" "Mozilla

	return &DockerLog{
		Type:      t,
		Labels:    labels,
		Timestamp: string(logmsg[0]),
		//IP:        string(part[2]),
		Message: string(logmsg[2]),
	}, nil
}
