package ctrl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	. "github.com/mickael-kerjean/webpty/common"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var resizeMessage = struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
	X    uint16
	Y    uint16
}{}

func HandleSocket(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("upgrade error"))
		return
	}
	defer conn.Close()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd")
	} else if _, err = exec.LookPath("/bin/bash"); err == nil {
		bashCommand := `bash --noprofile --init-file <(cat <<EOF
export TERM="xterm"
export PS1="\[\033[1;34m\]\w\[\033[0;37m\] # \[\033[0m\]"
export EDITOR="emacs"`
		bashCommand += strings.Join([]string{
			"",
			"export PATH=" + os.Getenv("PATH"),
			"export HOME=" + os.Getenv("HOME"),
			"",
		}, "\n")
		bashCommand += `
alias ls='ls --color'
alias ll='ls -lah'
EOF
)`
		cmd = exec.Command("/bin/bash", "-c", bashCommand)
	} else if _, err = exec.LookPath("/bin/sh"); err == nil {
		cmd = exec.Command("/bin/sh")
		cmd.Env = []string{
			"TERM=xterm",
			"PATH=" + os.Getenv("PATH"),
			"HOME=" + os.Getenv("HOME"),
		}
	} else {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("No terminal found"))
		Log.Error("No terminal found")
		return
	}

	tty, err := pty.Start(cmd)
	if err != nil {
		fmt.Printf("plugin::plg_handler_console pty.Start error '%s'\n", err)
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
		tty.Close()
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := tty.Read(buf)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				return
			}
			conn.WriteMessage(websocket.BinaryMessage, buf[:read])
		}
	}()

	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			return
		} else if messageType == websocket.TextMessage {
			conn.WriteMessage(websocket.TextMessage, []byte("Unexpected text message"))
			continue
		}

		dataTypeBuf := make([]byte, 1)
		read, err := reader.Read(dataTypeBuf)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Unable to read message type from reader"))
			return
		} else if read != 1 {
			return
		}

		switch dataTypeBuf[0] {
		case 0:
			b, err := io.ReadAll(reader)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error copying bytes: "+err.Error()))
				continue
			}
			_, err = tty.Write(b)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error writing bytes: "+err.Error()))
				continue
			}
		case 1:
			decoder := json.NewDecoder(reader)
			if err := decoder.Decode(&resizeMessage); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error decoding resize message: "+err.Error()))
				continue
			}
			if _, _, errno := syscall.Syscall(
				syscall.SYS_IOCTL,
				tty.Fd(),
				syscall.TIOCSWINSZ,
				uintptr(unsafe.Pointer(&resizeMessage)),
			); errno != 0 {
				conn.WriteMessage(websocket.TextMessage, []byte("Unable to resize terminal: "+err.Error()))
			}
		default:
			conn.WriteMessage(websocket.TextMessage, []byte("Unknown data type: "+err.Error()))
		}
	}
}
