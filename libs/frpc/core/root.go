package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/featuregate"
	"github.com/fatedier/frp/pkg/util/log"
)

type instance struct {
	Service    *client.Service
	CancelFunc context.CancelFunc
}

var (
	mutex         sync.Mutex
	isInitialized bool
	instances     map[string]*instance
	callback      func(action string, path string)
)

func Initialize(
	to string,
	level string,
	maxDays int,
	cb func(action string, path string),
) {
	mutex.Lock()
	defer mutex.Unlock()
	if isInitialized {
		return
	}

	log.InitLogger(to, level, maxDays, true)
	instances = make(map[string]*instance)
	callback = cb
	isInitialized = true
}

func Verify(path string) error {
	common, proxy, visitor, _, err := config.LoadClientConfig(path, true)
	if err != nil {
		return err
	}

	warning, err := validation.ValidateAllClientConfig(common, proxy, visitor)
	if warning != nil {
		return warning
	}
	if err != nil {
		return err
	}

	return nil
}

func Start(path string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if !isInitialized {
		return fmt.Errorf("please initialize first")
	}

	if _, exists := instances[path]; exists {
		log.Errorf("config [%s] is already running", path)
		return fmt.Errorf("config [%s] is already running", path)
	}

	common, proxy, visitor, _, err := config.LoadClientConfig(path, true)
	if err != nil {
		return err
	}

	if len(common.FeatureGates) > 0 {
		if err := featuregate.SetFromMap(common.FeatureGates); err != nil {
			return err
		}
	}

	svr, err := client.NewService(client.ServiceOptions{
		Common:         common,
		ProxyCfgs:      proxy,
		VisitorCfgs:    visitor,
		ConfigFilePath: path,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	instances[path] = &instance{
		Service:    svr,
		CancelFunc: cancel,
	}
	callback("started", path)

	go func() {
		if err := svr.Run(ctx); err != nil {
			log.Errorf("client for config file [%s] stopped with error: %v", path, err)
		}

		mutex.Lock()
		defer mutex.Unlock()
		delete(instances, path)
		callback("stopped", path)
	}()

	return nil
}

func Stop(path string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if !isInitialized {
		return fmt.Errorf("please initialize first")
	}

	instance, exisits := instances[path]

	if exisits {
		instance.CancelFunc()
	}

	return nil
}
