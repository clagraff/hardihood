package service

import "github.com/clagraff/hardihood/checkup"

type Service interface {
	Name() string
	Checkups() []checkup.Checkup
}

type service struct {
	name   string
	checks []checkup.Checkup
}

func (s service) Name() string { return s.name }
func (s service) Checkups() []checkup.Checkup {
	checks := []checkup.Checkup{}
	for _, c := range s.checks {
		checks = append(checks, c)
	}
	return checks
}

func Make(name string, checks []checkup.Checkup) Service {
	return service{
		name:   name,
		checks: checks,
	}
}
