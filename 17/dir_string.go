// Code generated by "stringer -type=Dir"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Up-0]
	_ = x[Right-1]
	_ = x[Down-2]
	_ = x[Left-3]
}

const _Dir_name = "UpRightDownLeft"

var _Dir_index = [...]uint8{0, 2, 7, 11, 15}

func (i Dir) String() string {
	if i < 0 || i >= Dir(len(_Dir_index)-1) {
		return "Dir(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Dir_name[_Dir_index[i]:_Dir_index[i+1]]
}
