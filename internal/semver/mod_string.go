// Code generated by "stringer -type=Mod -linecomment"; DO NOT EDIT.

package semver

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ModPatch-0]
	_ = x[ModMinor-1]
	_ = x[ModMajor-2]
	_ = x[ModDate-3]
}

const _Mod_name = "patchminormajordate"

var _Mod_index = [...]uint8{0, 5, 10, 15, 19}

func (i Mod) String() string {
	if i >= Mod(len(_Mod_index)-1) {
		return "Mod(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Mod_name[_Mod_index[i]:_Mod_index[i+1]]
}