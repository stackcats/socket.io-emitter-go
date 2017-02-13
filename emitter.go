package emitter

import (
	"bytes"
	"fmt"
	"gopkg.in/redis.v5"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const (
	// https://github.com/socketio/socket.io-parser/blob/master/index.js
	gEvent       = 2
	gBinaryEvent = 5
	uid          = "emitter"
)

// Options ...
type Options struct {
	// host to connect to redis on (localhost)
	Host string
	// port to connect to redis on (6379)
	Port int
	// the name of the key to pub/sub events on as prefix (socket.io)
	Key string
	// unix domain socket to connect to redis on ("/tmp/redis.sock")
	Socket string
}

// Emitter Socket.IO redis base emitter
type Emitter struct {
	redis  *redis.Client
	prefix string
	rooms  []string
	flags  map[string]interface{}
}

// NewEmitter Emitter constructor
func NewEmitter(opts *Options) *Emitter {
	emitter := &Emitter{}

	host := "127.0.0.1"
	if opts.Host != "" {
		host = opts.Host
	}

	port := 6379
	if opts.Port > 0 && opts.Port < 65536 {
		port = opts.Port
	}

	redisURI := fmt.Sprintf("%s:%d", host, port)
	emitter.redis = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})

	emitter.prefix = "socket.io"
	if opts.Key != "" {
		emitter.prefix = opts.Key
	}

	emitter.rooms = make([]string, 0, 0)

	emitter.flags = make(map[string]interface{})

	return emitter
}

// Close release redis client
func (e *Emitter) Close() {
	if e.redis != nil {
		e.redis.Close()
	}
}

// Emit Send the packet
func (e *Emitter) Emit(data ...interface{}) (*Emitter, error) {
	packet := make(map[string]interface{})
	packet["type"] = gEvent
	if hasBin(data...) {
		packet["type"] = gBinaryEvent
	}

	packet["data"] = data

	packet["nsp"] = "/"
	if nsp, ok := e.flags["nsp"]; ok {
		packet["nsp"] = nsp
		delete(e.flags, "nsp")
	}

	opts := map[string]interface{}{
		"rooms": e.rooms,
		"flags": e.flags,
	}

	chn := fmt.Sprintf("%s#%s#", e.prefix, packet["nsp"])

	buf, err := msgpack.Marshal([]interface{}{uid, packet, opts})
	if err != nil {
		return nil, err
	}

	fmt.Println(string(buf))
	if len(e.rooms) > 0 {
		for _, room := range e.rooms {
			chnRoom := fmt.Sprintf("%s%s#", chn, room)
			e.redis.Publish(chnRoom, string(buf))
		}
	} else {
		e.redis.Publish(chn, string(buf))
	}

	e.rooms = make([]string, 0, 0)
	e.flags = make(map[string]interface{})
	return e, nil
}

// In Limit emission to a certain `room`
func (e *Emitter) In(room string) *Emitter {
	for _, r := range e.rooms {
		if r == room {
			return e
		}
	}
	e.rooms = append(e.rooms, room)
	return e
}

// To Limit emission to a certain `room`
func (e *Emitter) To(room string) *Emitter {
	return e.In(room)
}

// Of Limit emission to certain `namespace`
func (e *Emitter) Of(namespace string) *Emitter {
	e.flags["nsp"] = namespace
	return e
}

// JSON flag
func (e *Emitter) JSON() *Emitter {
	e.flags["json"] = true
	return e
}

// Volatile flag
func (e *Emitter) Volatile() *Emitter {
	e.flags["volatile"] = true
	return e
}

// Broadcast flag
func (e *Emitter) Broadcast() *Emitter {
	e.flags["broadcast"] = true
	return e
}

func hasBin(data ...interface{}) bool {
	if data == nil {
		return false
	}

	for _, d := range data {
		switch res := d.(type) {
		case []byte:
			return true
		case bytes.Buffer:
			return true
		case []interface{}:
			for _, each := range res {
				if hasBin(each) {
					return true
				}
			}
		case map[string]interface{}:
			for _, val := range res {
				if hasBin(val) {
					return true
				}
			}
		default:
			return false
		}
	}

	return false
}
