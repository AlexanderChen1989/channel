package server

import "github.com/garyburd/redigo/redis"

const (
	FLUSHDB    = "FLUSHDB"
	MULTI      = "MULTI"
	EXEC       = "EXEC"
	LPUSH      = "LPUSH"
	RPOP       = "RPOP"
	BRPOPLPUSH = "BRPOPLPUSH"
	HMSET      = "HMSET"
	HGETALL    = "HGETALL"
)

type Command struct {
	CmdName string
	Args    []interface{}
}

func Cmd(cmd string, args ...interface{}) Command {
	return Command{CmdName: cmd, Args: args}
}

func Do(pool *redis.Pool, cmd Command) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	return conn.Do(cmd.CmdName, cmd.Args...)
}

func Pipe(pool *redis.Pool, cmd1 Command, cmd2 Command, cmds ...Command) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send(MULTI)
	conn.Send(cmd1.CmdName, cmd1.Args...)
	conn.Send(cmd2.CmdName, cmd2.Args...)

	for _, cmd := range cmds {
		conn.Send(cmd.CmdName, cmd.Args...)
	}

	return conn.Do(EXEC)
}
