// Code generated by "stringer -type=AccessModel"; DO NOT EDIT.

package skydb

import "fmt"

const _AccessModel_name = "RoleBasedAccessRelationBasedAccess"

var _AccessModel_index = [...]uint8{0, 15, 34}

func (i AccessModel) String() string {
	i -= 1
	if i < 0 || i >= AccessModel(len(_AccessModel_index)-1) {
		return fmt.Sprintf("AccessModel(%d)", i+1)
	}
	return _AccessModel_name[_AccessModel_index[i]:_AccessModel_index[i+1]]
}
