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

package progress

import "fmt"

type Progress struct {
	percent int64
	current int64
	total   int64
	rate    string
	graph   string
}

func NewProgress(start, total int64) *Progress {
	p := new(Progress)
	p.current = start
	p.total = total
	p.graph = "="
	p.percent = p.GetPercent()
	return p
}

func (p *Progress) GetPercent() int64 {
	return int64(float32(p.current) / float32(p.total) * 100)
}

func (p *Progress) Add(i int64) {

	p.current += i

	if p.current > p.total {
		return
	}

	last := p.percent
	p.percent = p.GetPercent()

	if p.percent != last && p.percent%2 == 0 {
		p.rate += p.graph
	}

	fmt.Printf("\r[%-50s]%3d%% %8d/%d", p.rate, p.percent, p.current, p.total)

	if p.current == p.total {
		p.Done()
	}
}

func (p *Progress) Done() {
	fmt.Println()
}
