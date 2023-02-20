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

import (
	"context"
	"encoding/json"
)

type ClonerI interface {
	Clone(ctx context.Context, src any, dst any) error
}

type cloner struct {
}

func NewCloner() *cloner {
	return &cloner{}
}

func (c *cloner) Clone(ctx context.Context, src any, dst any) error {
	srcStr, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(srcStr, dst)
	if err != nil {
		return err
	}
	return nil
}
