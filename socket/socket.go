package socket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JustBugLord/investigo/bypass"
	"github.com/gorilla/websocket"
)

type InvestigoSocket struct {
	rng        *bypass.RandomGenerator
	dialer     websocket.Dialer
	rwm        sync.RWMutex
	connection *websocket.Conn
	ticker     *time.Ticker
	handlers   map[EventType]func(*WSResponse) error
	ctx        context.Context
	cancel     context.CancelFunc
	errHandler func(error)
	health     bool
}

func NewInvestigoSocketDefault() *InvestigoSocket {
	return &InvestigoSocket{
		rng: bypass.NewRandomGenerator(),
		dialer: websocket.Dialer{
			ReadBufferSize:    8192,
			WriteBufferSize:   8192,
			EnableCompression: true,
			Proxy:             http.ProxyFromEnvironment,
			HandshakeTimeout:  45 * time.Second,
		},
		handlers: make(map[EventType]func(*WSResponse) error),
		errHandler: func(err error) {
			panic(err)
		},
	}
}

func (is *InvestigoSocket) Connect() error {
	id, err := is.rng.NumberString(1e3)
	if err != nil {
		return errors.New("could not generate random number: " + err.Error())
	}
	sign, err := is.rng.String(8)
	if err != nil {
		return errors.New("could not generate random signature: " + err.Error())
	}
	conn, _, err := is.dialer.Dial(fmt.Sprintf("wss://streaming.forexpros.com/echo/%s/%s/websocket", id, sign), nil)
	if err != nil {
		return errors.New("fail open public socket: " + err.Error())
	}
	is.connection = conn
	ctx, cancel := context.WithCancel(context.Background())
	is.ctx = ctx
	is.cancel = cancel
	is.health = true
	is.ping()
	is.channel()
	return nil
}

func (is *InvestigoSocket) ping() {
	go func() {
		ticker := time.NewTicker(40 * time.Second)
		is.ticker = ticker
		for {
			select {
			case <-is.ctx.Done():
				return
			case <-ticker.C:
				if err := is.Write(websocket.TextMessage, []byte("[\"{ \"_event\": \"heartbeat\", \"data\": \"h\"}\"]")); err != nil {
					if is.errHandler != nil {
						is.errHandler(errors.New("fail write ping to connection: " + err.Error()))
					}
				}
			}
		}
	}()
}

func (is *InvestigoSocket) channel() {
	go func() {
		for {
			select {
			case <-is.ctx.Done():
				return
			default:
				_, data, err := is.Read()
				if err != nil {
					is.errHandler(err)
				}
				if len(data) == 0 {
					continue
				}
				if data[0] == 'a' {
					response := WSResponseFromRaw(data)
					if response != nil && response.Event != Heartbeat {
						if value, ok := is.handlers[response.Event]; ok && value != nil {
							if err := value(response); err != nil {
								is.errHandler(err)
							}
						}
					}
				}
			}
		}
	}()
}

func (is *InvestigoSocket) Subscribe(req *WSRequest) error {
	if err := is.SendRequest(websocket.TextMessage, req); err != nil {
		return err
	}
	return nil
}

func (is *InvestigoSocket) AddHandler(event EventType, handler func(*WSResponse) error) {
	if is.handlers == nil {
		is.handlers = make(map[EventType]func(*WSResponse) error)
	}
	is.rwm.Lock()
	defer is.rwm.Unlock()
	is.handlers[event] = handler
}

func (is *InvestigoSocket) SetErrHandler(handler func(err error)) {
	if handler == nil {
		return
	}
	is.errHandler = func(err error) {
		is.health = false
		handler(err)
	}
}

func (is *InvestigoSocket) Health() bool {
	is.rwm.RLock()
	defer is.rwm.RUnlock()
	return is.health
}

func (is *InvestigoSocket) SendRequest(msgType int, request *WSRequest) error {
	if err := is.Write(msgType, []byte(request.String())); err != nil {
		return errors.New("fail send message: " + err.Error())
	}
	return nil
}

func (is *InvestigoSocket) ReadToStruct(to any) error {
	_, bytes, err := is.Read()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, to); err != nil {
		return errors.New("fail to unmarshal json: " + err.Error())
	}
	return nil
}

func (is *InvestigoSocket) Write(msgType int, bytes []byte) error {
	is.rwm.Lock()
	defer is.rwm.Unlock()
	if err := is.connection.WriteMessage(msgType, bytes); err != nil {
		return errors.New("fail write to public socket: " + err.Error())
	}
	return nil
}

func (is *InvestigoSocket) Read() (int, []byte, error) {
	msgType, bytes, err := is.connection.ReadMessage()
	if err != nil {
		return 0, nil, errors.New("fail read from public socket: " + err.Error())
	}
	return msgType, bytes, nil
}

func (is *InvestigoSocket) Close() {
	if is.connection != nil {
		is.connection.Close()
	}
	if is.ticker != nil {
		is.ticker.Stop()
	}
	if is.cancel != nil {
		is.cancel()
	}
	is.health = false
}
