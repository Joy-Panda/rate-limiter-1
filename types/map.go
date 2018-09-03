package types

import (
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
)

// NewMap 创建一个非并发安全的Map实例
func NewMap() *PB_Map {
	return &PB_Map{
		Value: make(map[string]*any.Any),
	}
}

// Set 向map中设置KV
func (m *PB_Map) Set(key string, v proto.Message) error {
	any, err := NewAny(v)
	if err != nil {
		return err
	}
	m.Value[key] = any
	return nil
}

// Get 从map中获取值, 如果没有找到将返回nil
func (m *PB_Map) Get(key string) *any.Any {
	if v, exist := m.Value[key]; exist {
		return v
	}
	return nil
}
