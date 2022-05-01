// Copyright 2022 cedar12, cedar12.zxd@qq.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"gopkg.in/ini.v1"
)

type IniParser struct {
	fileName   string
	confReader *ini.File // config reader
}

type IniParserError struct {
	errorInfo string
}

func (e *IniParserError) Error() string { return e.errorInfo }

func (p *IniParser) Load(configFileName string) error {

	conf, err := ini.Load(configFileName)
	if err != nil {
		p.confReader = nil
		return err
	}
	p.fileName = configFileName
	p.confReader = conf
	return nil
}

func (p *IniParser) GetString(section string, key string) string {
	if p.confReader == nil {
		return ""
	}

	s := p.confReader.Section(section)
	if s == nil {
		return ""
	}

	return s.Key(key).String()
}

func (p *IniParser) GetInt32(section string, key string) int32 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int()

	return int32(valueInt)
}

func (p *IniParser) GetUint32(section string, key string) uint32 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint()

	return uint32(valueInt)
}

func (p *IniParser) GetInt64(section string, key string) int64 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Int64()
	return valueInt
}

func (p *IniParser) GetUint64(section string, key string) uint64 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueInt, _ := s.Key(key).Uint64()
	return valueInt
}

func (p *IniParser) GetFloat32(section string, key string) float32 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return float32(valueFloat)
}

func (p *IniParser) GetFloat64(section string, key string) float64 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	valueFloat, _ := s.Key(key).Float64()
	return valueFloat
}

func (p *IniParser) Reload() error {
	return p.confReader.Reload()
}

func (p *IniParser) FileName() string {
	return p.fileName
}
