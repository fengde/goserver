package serviceSentinel

import "github.com/alibaba/sentinel-golang/core/flow"

// 这里集中维护限流、熔断规则
func LoadRules() error {
	var rules []*flow.Rule
	rules = append(rules, &flow.Rule{
		Resource:               "test",
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
		StatIntervalInMs:       1000,
		Threshold:              10,
	})
	_, err := flow.LoadRules(rules)
	return err
}
