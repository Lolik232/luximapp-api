package domain

import "github.com/Lolik232/luximapp-api/internal/context_keys"

type Log struct {
	ID     string
	Msg    string
	Data   context_keys.InputData
	Device string 
}
