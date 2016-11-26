package cli

type Socket struct {
	endpoint  string
	onOpen    []func()
	onClose   []func()
	onError   []func()
	onMessage []func(string)

	hb struct {
		sync.Mutex
		done chan struct{}
		started bool
		timeout func() time.Duration
	}

	rc struct {
		sync.Mutex
		done chan struct{}
		started bool
		timeout func() time.Duration
	}
}

func (sock *Socket) startReconnect() {
	sock.rc.Lock()
	defer sock.rc.Unlock()

	if (sock.rc.started) {return}
	sock.rc.started = true
	sock.rc.done = make(chan struct{})
	go func() {
		done := sock.rc.done
		for {
			select {
			case <-done:
				return
			default:
				sock.Connect()	
			}
			time.Sleep(sock.rc.timeout())
		}
	}()
}

func (sock *Socket) stopReconnect() {
	sock.rc.Lock()
	defer sock.rc.Unlock()
	if (!sock.rc.started) {return}
	close(sock.rc.done)
	sock.rc.started = false
}

func (sock *Socket) Connect() {

}

func (sock *Socket) startHeartbeat() {
	sock.hb.Lock()
	defer sock.hb.Unlock()

	if (sock.hb.started) {return}

	sock.hb.started = true
	sock.hb.done = make(chan struct{})

	go func() {
		done := sock.hb.done
		for {
			select {
			case <-done:
				return
			default:
				sock.send("")
			}
			time.Sleep(sock.hb.timeout())
		}
	}()
}

func (sock *Socket) stopHeartbeat() {
	sock.hb.Lock()
	defer sock.hb.Unlock()

	if (!sock.hb.started) {return}
	sock.hb.started = false
	close(sock.hb.done)
}


func (sock *Socket) send(m string) {}

func (sock *Socket) flushSendBuffer() {	
}

func (sock *Socket) onConnOpen() {
	sock.flushSendBuffer()
	sock.stopReconnect()
	sock.startHeartbeat()

	for _, callback := range sock.onOpen {
		callback()
	}
}

func (sock *Socket) OnOpen(callback func()) {
	sock.onOpen = append(sock.onOpen, callback)
}

func (sock *Socket) OnClose(callback func()) {
	sock.onClose = append(sock.onClose, callback)
}

func (sock *Socket) OnError(callback func()) {
	sock.onError = append(sock.onError, callback)
}

func (sock *Socket) OnMessage(callback func()) {
	sock.onMessage = append(sock.onMessage, callback)
}




