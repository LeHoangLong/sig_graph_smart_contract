package service

import (
	"fmt"
	"sig_graph/controller"
	"sig_graph/utility"
	"strings"
)

type nodeNameService struct {
	settings utility.SettingsI
}

func NewNodeNameService(
	settings utility.SettingsI,
) controller.NodeNameServiceI {
	return &nodeNameService{
		settings: settings,
	}
}

func (s *nodeNameService) GenerateFullId(id string) (string, error) {
	if s.IsFullId(id) {
		return id, nil
	}

	if !s.IsIdValid(id) {
		return "", fmt.Errorf("%w: invalid id", utility.ErrInvalidArgument)
	}

	return fmt.Sprintf("%s:%s", s.settings.GraphName(), id), nil
}

func (s *nodeNameService) IsIdValid(id string) bool {
	if strings.Contains(id, ":") {
		return false
	}

	return true
}

func (s *nodeNameService) IsFullId(id string) bool {
	graphName := fmt.Sprintf("%s:", s.settings.GraphName())
	if !strings.Contains(id, graphName) {
		return false
	}

	splittedStr := strings.Split(id, graphName)

	if len(splittedStr) != 2 || len(splittedStr[0]) != 0 {
		return false
	}

	if !s.IsIdValid(splittedStr[1]) {
		return false
	}

	return true
}
