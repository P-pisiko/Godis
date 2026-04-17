package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"DEL":     del,
	"HDEL":    hdel,
	"EXITS":   exists,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk
	SETsMu.RLock()
	value, err := SETs[key]
	SETsMu.RUnlock()

	if !err {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	fields, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}
	for k, v := range fields {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}

func del(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}

	var deleted int

	SETsMu.Lock()
	for _, arg := range args {
		key := arg.bulk
		if _, exists := SETs[key]; exists {
			delete(SETs, key)
			deleted++
		}
	}
	SETsMu.Unlock()

	return Value{typ: "integer", num: deleted}
}

func hdel(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hdel' command"}
	}

	hash := args[0].bulk
	var deleted int

	HSETsMu.Lock()
	if fields, exists := HSETs[hash]; exists {
		for _, arg := range args[1:] {
			key := arg.bulk
			if _, fieldExists := fields[key]; fieldExists {
				delete(fields, key)
				deleted++
			}
		}
		if len(HSETs[hash]) == 0 {
			delete(HSETs, hash)
		}
	}
	HSETsMu.Unlock()

	return Value{typ: "integer", num: deleted}
}

func exists(args []Value) Value {
	if len(args) <= 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'exits' command"}
	}
	var count int

	SETsMu.RLock()
	for _, arg := range args {
		if _, ok := SETs[arg.bulk]; ok {
			count++
		}
	}
	SETsMu.RUnlock()
	return Value{typ: "integer", num: count}
}
