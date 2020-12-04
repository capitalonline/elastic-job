package icontext

import (
	"context"
	"time"
)

type DefaultContext struct {
	clusterId string
	ctx       context.Context
	Keys      map[string]string
	err       error
}

func Background() *DefaultContext {
	c := &DefaultContext{
		ctx:  context.Background(),
		Keys: make(map[string]string),
	}
	return c
}

//Implement context interface
func (c *DefaultContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *DefaultContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *DefaultContext) Err() error {
	if c.err != nil {
		return c.err
	}
	return c.ctx.Err()
}

func (c *DefaultContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// Cluster metas
func (c *DefaultContext) GetContext() context.Context {
	return c.ctx
}

func (c *DefaultContext) GetClusterId() string {
	return c.clusterId
}

func (c *DefaultContext) SetClusterId(clusterId string) {
	c.clusterId = clusterId
}

func (c *DefaultContext) SetError(err error) {
	c.err = err
}

// Implement Cache
func (c *DefaultContext) Set(key string, value string) {
	if c.Keys == nil {
		c.Keys = make(map[string]string)
	}
	c.Keys[key] = value
}

func (c *DefaultContext) Get(key string) (value string, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *DefaultContext) WithCancel() (DefaultContext, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.ctx)
	return DefaultContext{
		clusterId: c.clusterId,
		Keys:      c.Keys,
		ctx:       ctx,
	}, cancel
}

func (c *DefaultContext) WithTimeout(timeout int) (DefaultContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Duration(timeout)*time.Second)
	return DefaultContext{
		clusterId: c.clusterId,
		Keys:      c.Keys,
		ctx:       ctx,
	}, cancel
}
