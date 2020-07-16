// Copyright 2019 morgine.com. All rights reserved.

package cfg

import (
	"github.com/BurntSushi/toml"
)

func Decode(data []byte) (Env, error) {
	env := make(Env)
	err := toml.Unmarshal(data, &env)
	return env, err
}
