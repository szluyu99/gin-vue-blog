package handle

import "testing"

func TestCheckQueryPage(t *testing.T) {
	{
		pageNum, pageSize := 0, 0
		num, size := checkQueryPage(pageNum, pageSize)
		if num != 1 || size != 10 {
			t.Errorf("checkQueryPage(%d, %d) = %d, %d; want 1, 10", pageNum, pageSize, num, size)
		}
	}
	{
		pageNum, pageSize := 1, 5
		num, size := checkQueryPage(pageNum, pageSize)
		if num != 1 || size != 5 {
			t.Errorf("checkQueryPage(%d, %d) = %d, %d; want 1, 5", pageNum, pageSize, num, size)
		}
	}
}
