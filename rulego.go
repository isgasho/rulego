/*
 * Copyright 2023 The RuleGo Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rulego

import (
	"sync"
)

var DefaultRuleGo = &RuleGo{}

//RuleGo 规则引擎实例池
type RuleGo struct {
	ruleEngines sync.Map
}

// New creates a new RuleEngine and stores it in the RuleGo.
func (g *RuleGo) New(id string, rootRuleChainSrc []byte, opts ...RuleEngineOption) (*RuleEngine, error) {
	if v, ok := g.ruleEngines.Load(id); ok {
		return v.(*RuleEngine), nil
	} else {
		if ruleEngine, err := newRuleEngine(id, rootRuleChainSrc, opts...); err != nil {
			return nil, err
		} else {
			// Store the new RuleEngine in the ruleEngines map with the Id as the key.
			g.ruleEngines.Store(ruleEngine.Id, ruleEngine)
			return ruleEngine, err
		}

	}
}

//Get 获取指定ID规则引擎实例
func (g *RuleGo) Get(id string) (*RuleEngine, bool) {
	v, ok := g.ruleEngines.Load(id)
	if ok {
		return v.(*RuleEngine), ok
	} else {
		return nil, false
	}

}

//Del 删除指定ID规则引擎实例
func (g *RuleGo) Del(id string) {
	v, ok := g.ruleEngines.Load(id)
	if ok {
		v.(*RuleEngine).Stop()
		g.ruleEngines.Delete(id)
	}

}

//Stop 释放所有规则引擎实例
func (g *RuleGo) Stop() {
	g.ruleEngines.Range(func(key, value any) bool {
		if item, ok := value.(*RuleEngine); ok {
			item.Stop()
		}
		g.ruleEngines.Delete(key)
		return true
	})
}

// New creates a new RuleEngine and stores it in the RuleGo.
func New(id string, rootRuleChainSrc []byte, opts ...RuleEngineOption) (*RuleEngine, error) {
	return DefaultRuleGo.New(id, rootRuleChainSrc, opts...)
}

//Get 获取指定ID规则引擎实例
func Get(id string) (*RuleEngine, bool) {
	return DefaultRuleGo.Get(id)
}

//Del 删除指定ID规则引擎实例
func Del(id string) {
	DefaultRuleGo.Del(id)
}

//Stop 释放所有规则引擎实例
func Stop() {
	DefaultRuleGo.Stop()
}
