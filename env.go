// Copyright 2019 morgine.com. All rights reserved.

package cfg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strconv"
)

var OsEnvGetter = func(namespace, key string) string {
	if namespace != "" {
		key = namespace + "." + key
	}
	return os.Getenv(key)
}

type NamespaceError struct {
	Namespace string
	Err       string
}

func (n NamespaceError) Error() string {
	return fmt.Sprintf("config namespace [%s]: %s", n.Namespace, n.Err)
}

type Env map[string]interface{}

// Unmarshal for decoding Env data to schema
func (e Env) Unmarshal(schema interface{}) error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(e)
	if err != nil {
		return err
	}
	return toml.Unmarshal(buf.Bytes(), schema)
}

// UnmarshalSub is shorthand for GetSub and Unmarshal
func (e Env) UnmarshalSub(namespace string, schema interface{}) error {
	envs, err := e.GetSub(namespace)
	if err != nil {
		return err
	}
	return envs.Unmarshal(schema)
}

// Get sub environments and load os environment values
func (e Env) GetSub(namespace string) (env Env, err error) {
	subs, ok := e[namespace]
	if !ok {
		return nil, NamespaceError{
			Namespace: namespace,
			Err:       "not exist",
		}
	}

	mp, o := subs.(map[string]interface{})
	if !o {
		return nil, NamespaceError{
			Namespace: namespace,
			Err:       "config data only support 'key = value' format",
		}
	} else {
		env = mp
		err = env.LoadOSEnv(namespace)
		if err != nil {
			return nil, err
		}
		return env, nil
	}
}

func (e Env) Len() int {
	return len(e)
}

func (e Env) GetStr(name string) string {
	return e[name].(string)
}

func (e Env) GetInt(name string) int {
	return e[name].(int)
}

func (e Env) GetFloat(name string) float64 {
	return e[name].(float64)
}

func (e Env) GetBool(name string) bool {
	return e[name].(bool)
}

func (e Env) GetSliceStr(name string) []string {
	return e[name].([]string)
}

func (e Env) GetSliceInt(name string) []int {
	return e[name].([]int)
}

func (e Env) GetSliceFloat(name string) []float64 {
	return e[name].([]float64)
}

func (e Env) GetSliceBool(name string) []bool {
	return e[name].([]bool)
}

// LoadOSEnv for loading the OS environment values, if namespace
func (e Env) LoadOSEnv(namespace string) (err error) {
	for key, value := range e {
		ov := OsEnvGetter(namespace, key)
		if ov != "" {
			switch value.(type) {
			case string:
				e[key] = ov
			case int:
				e[key], err = strconv.Atoi(ov)
				if err != nil {
					return fmt.Errorf("cfg.LoadOSEnv: key: %s]: %s", key, err)
				}
			case bool:
				switch ov {
				case "y", "Y", "yes", "YES", "Yes", "1", "t", "T", "true", "TRUE", "True":
					e[key] = true
				case "n", "N", "no", "NO", "No", "0", "f", "F", "false", "FALSE", "False":
					e[key] = false
				default:
					return fmt.Errorf("OS env value [%s]: need boolean", key)
				}
			case float64:
				e[key], err = strconv.ParseFloat(ov, 64)
				if err != nil {
					return fmt.Errorf("OS env value [%s]: %s", key, err)
				}
			default:
				return errors.New("os config value only support 'string', 'int', 'float64' or 'bool'")
			}
		}
	}
	return nil
}
