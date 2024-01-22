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
		cmd = exec.CommandContext(req.Context(), "cmd")
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
		cmd = exec.CommandContext(req.Context(), "/bin/bash", "-c", bashCommand)
	} else if _, err = exec.LookPath("/bin/sh"); err == nil {
		cmd = exec.CommandContext(req.Context(), "/bin/sh")
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

	Log.Info("connected client")
	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				Log.Error("socket.go::nextReader unexpected close error %s", err.Error())
				return
			}
			Log.Debug("socket.go::disconnection_event - browser is disconnected")
			return
		} else if messageType == websocket.TextMessage {
			Log.Debug("socket.go::expectation_failed - Unexpected text message")
			return
		}

		dataTypeBuf := make([]byte, 1)
		read, err := reader.Read(dataTypeBuf)
		if err != nil {
			Log.Error("socket.go::error - Unable to read message type from reader")
			return
		} else if read != 1 {
			Log.Error("socket.go::expectation_failed - Unexpected message size")
			return
		}
		switch dataTypeBuf[0] {
		case 0:
			Log.Info("disconnected")
			return
		case 1:
			b, err := io.ReadAll(reader)
			if err != nil {
				Log.Error("socket.go::error - copying bytes: %s", err.Error())
				return
			}
			_, err = tty.Write(b)
			if err != nil {
				Log.Error("socket.go::error - writing bytes: %s", err.Error())
				return
			}
		case 2:
			decoder := json.NewDecoder(reader)
			if err := decoder.Decode(&resizeMessage); err != nil {
				Log.Error("socket.go::error - decoding resize message: %s", err.Error())
				return
			}
			if _, _, errno := syscall.Syscall(
				syscall.SYS_IOCTL,
				tty.Fd(),
				syscall.TIOCSWINSZ,
				uintptr(unsafe.Pointer(&resizeMessage)),
			); errno != 0 {
				Log.Error("socket.go::expectation_failed - errno[%+v]", errno)
				return
			}
		default:
			Log.Error("socket.go::expectation_failed - unknown socket data type: %+v", dataTypeBuf)
			return
		}
	}
}
