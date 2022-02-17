package arraysHelper

func IndexOf(array []int, element int) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	return -1
}
