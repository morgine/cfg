// Copyright 2020 morgine.com. All rights reserved.

package cfg_test

import (
	"fmt"
	"github.com/morgine/cfg"
	"github.com/morgine/service"
	"os"
)

func ExampleEnv_Unmarshal() {
	var data = `
[mysql]
host = "127.0.0.1"
port= "3306"
`

	// Config structure
	type mysql struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	}

	// Decode data to cfg.Env
	envs, err := cfg.Decode([]byte(data))
	if err != nil {
		panic(err)
	}

	// Set os environment. OPTIONAL
	err = os.Setenv("mysql.host", "localhost")
	if err != nil {
		panic(err)
	}

	// Get sub environment values, it use OS environment value if the value exist
	mysqlEnvs, err := envs.GetSub("mysql")
	if err != nil {
		panic(err)
	}
	fmt.Println(mysqlEnvs["host"])
	fmt.Println(mysqlEnvs["port"])

	config := mysql{}
	// Marshal mysqlEnvs to TOML data, and then Unmarshal data into config
	err = mysqlEnvs.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Host)
	fmt.Println(config.Port)

	config = mysql{}
	_ = os.Setenv("mysql.host", "127.0.0.1")
	// UnmarshalSub is shorthand for GetSub and Unmarshal
	err = envs.UnmarshalSub("mysql", &config)
	fmt.Println(config.Host)

	// Output:
	// localhost
	// 3306
	// localhost
	// 3306
	// 127.0.0.1
}

// NewService provide a configuration service, usually many other services depend on
// configuration service, this example shows how to inject configuration service
// to other services
func ExampleNewService() {
	var data = `
[client]
addr = "127.0.0.1:3306"
`

	// configuration service
	configService := cfg.NewService(cfg.NewMemoryStorageService(data))

	// client service, client service depend on configuration service
	cltService := newClientService("client", configService)

	// dependency injection container
	container := service.NewContainer()

	// get client from dependency injection container
	c, err := cltService.Get(container)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.addr)
	// Output:
	// 127.0.0.1:3306
}

type clientService struct {
	provider service.Provider
}

type client struct {
	addr string
}

type config struct {
	Addr string `toml:"addr"`
}

// Get singleton client
func (r *clientService) Get(ctn *service.Container) (*client, error) {
	c, err := ctn.Get(&r.provider)
	if err != nil {
		return nil, err
	} else {
		return c.(*client), nil
	}
}

func newClientService(namespace string, cfgService *cfg.Service) *clientService {
	// client provider
	var provider service.ProviderFunc = func(ctn *service.Container) (interface{}, error) {
		// Get the singleton envs
		envs, err := cfgService.Get(ctn)
		if err != nil {
			return nil, nil
		}
		c := config{}
		err = envs.UnmarshalSub(namespace, &c)
		if err != nil {
			return nil, err
		}
		return &client{addr: c.Addr}, nil
	}
	return &clientService{
		provider: provider,
	}
}
