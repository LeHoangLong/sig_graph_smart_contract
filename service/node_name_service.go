// Copyright (C) 2022 Le Hoang Long
// This file is part of SigGraph smart contract <https://github.com/LeHoangLong/sig_graph_smart_contract>.
//
// SigGraph is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SigGraph is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SigGraph.  If not, see <http://www.gnu.org/licenses/>.
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

func (s *nodeNameService) GenerateFullId(id string) (string, utility.Error) {
	if s.IsFullId(id) {
		return id, nil
	}

	if !s.IsIdValid(id) {
		return "", utility.NewError(utility.ErrInvalidArgument).AddMessage("invalid id")
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
