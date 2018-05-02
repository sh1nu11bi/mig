// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Zack Mullaly zmullaly@mozilla.com [:zack]

package modules

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Module is implemented by types that contain parameters for a module
// supported by the MIG agent.  The `ToParameters` method is expected to
// validate the module configuration data and return a data type that
// can be set as the `Parameters` field in an `Action`'s
// [Operation.Parameters](https://github.com/mozilla/mig/blob/master/action.go#L84)
// field.
type Module interface {
	Name() string
	ToParameters() (interface{}, error)
}

// FromMap attempts to populate a `Module` with data from a `map` containing
// configuration data for a module specified by `moduleName`.
func FromMap(moduleName string, jsonMap map[string]interface{}) (Module, error) {
	pkg := Pkg{}

	moduleContainers := map[string]Module{
		pkg.Name(): pkg,
	}

	module, found := moduleContainers[moduleName]
	if !found {
		return InvalidModule{}, errors.New("Not a recognized module.")
	}

	encoded, encodeErr := json.Marshal(&jsonMap)
	if encodeErr != nil {
		return InvalidModule{}, encodeErr
	}

	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decodeErr := decoder.Decode(&module)
	if decodeErr != nil {
		return InvalidModule{}, decodeErr
	}

	return module, nil
}
