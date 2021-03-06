// Code generated by "stringer -type=ItemType -output=gen_ItemType.go"; DO NOT EDIT.

package itemtype

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ItemError-0]
	_ = x[ItemEOF-1]
	_ = x[ItemText-2]
	_ = x[ItemCodeFence-3]
	_ = x[ItemCode-4]
	_ = x[ItemTmplBlockStart-5]
	_ = x[ItemJsonBlockStart-6]
	_ = x[ItemBlockEnd-7]
	_ = x[ItemCommEnd-8]
	_ = x[ItemArg-9]
	_ = x[ItemQuoteArg-10]
	_ = x[ItemArgComment-11]
	_ = x[ItemOption-12]
}

const _ItemType_name = "ItemErrorItemEOFItemTextItemCodeFenceItemCodeItemTmplBlockStartItemJsonBlockStartItemBlockEndItemCommEndItemArgItemQuoteArgItemArgCommentItemOption"

var _ItemType_index = [...]uint8{0, 9, 16, 24, 37, 45, 63, 81, 93, 104, 111, 123, 137, 147}

func (i ItemType) String() string {
	if i < 0 || i >= ItemType(len(_ItemType_index)-1) {
		return "ItemType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ItemType_name[_ItemType_index[i]:_ItemType_index[i+1]]
}
