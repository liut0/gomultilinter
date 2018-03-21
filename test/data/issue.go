package data

func testIssue() {
	// issue
}

func TestIssue2() {
}

// nolint: testlinter
func testIssue3() {
	// issue
}

// nolint
func testIssue4() {
}

func testIssue5() int {
	for i := 0; i < 50; i++ {
		switch i {
		case 1:
		case 2:
		case 3:
		case 4:
		case 5:
		case 6:
		}
	}
	return 0
}

func testIssue6() {
	// ignored
}
