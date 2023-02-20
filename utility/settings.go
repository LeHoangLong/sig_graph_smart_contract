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
package utility

//go:generate mockgen -source=$GOFILE -destination ../testutils/utility.go -package mock
type SettingsI interface {
	GraphName() string
	MaxTimeDifference_ms() uint64
}

type settings struct {
	graphName string
}

func NewSettings(
	graphName string,
) SettingsI {
	return &settings{
		graphName: graphName,
	}
}

func (s *settings) GraphName() string {
	return s.graphName
}

func (s *settings) MaxTimeDifference_ms() uint64 {
	/// maximum 5 minutes by default
	return 300000
}
