package grules

import (
	jsonencoding "encoding/json"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
)

type Engine struct {
	Rules []Rule
}

type Result struct {
	Ok     bool
	Rule   string
	Action string
}

func Init() Engine {
	return Engine{}
}

func (e *Engine) SetRules(json string) error {
	var r []Rule
	err := jsonencoding.NewDecoder(strings.NewReader(json)).Decode(&r)
	if err != nil {
		return err
	}
	e.Rules = r

	return nil
}

func (e *Engine) EvaluateRules(json string) []Result {
	results := make(chan Result, len(e.Rules))
	var wg sync.WaitGroup
	wg.Add(len(e.Rules))
	object := gjson.Parse(json)
	for _, rule := range e.Rules {
		go func(rule Rule, object gjson.Result, wg *sync.WaitGroup) {
			defer wg.Done()
			ok := evaluateObject(object, rule)
			results <- Result{
				Ok:     ok,
				Rule:   rule.Name,
				Action: rule.Action,
			}
		}(rule, object, &wg)
	}
	wg.Wait()
	close(results)
	res := make([]Result, 0)
	for result := range results {
		res = append(res, result)
	}
	return res
}
