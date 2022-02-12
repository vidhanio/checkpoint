package student

import (
	"encoding/json"
	"testing"
)

func TestMakeStudentsFromBytes(t *testing.T) {
	b := []byte(`
		[
			[
				[
					{"initials": ["A", "A"], "grade": 1, "teacher_initial": "A", "student_number": 1}
				]
			]
		]
	`)

	ss := MakeStudentsFromBytes(b)

	expected := Students(
		[26][26][]Student{
			{
				{
					{
						Initials:       [2]rune{'A', 'A'},
						Grade:          1,
						TeacherInitial: 'A',
						StudentNumber:  1,
					},
				},
			},
		},
	)

	if ss[0][0][0] != expected[0][0][0] {
		t.Errorf("expected %v, got %v", expected, ss)
	}

	b = []byte(`
	[
		[
			[
				{"initials": ["A", "A"], "grade": 1, "teacher_initial": "A", "student_number": 1}
			]
		]
	`)

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic, got nil")
		}
	}()

	ss = MakeStudentsFromBytes(b)
}

func TestVerify(t *testing.T) {
	students := []byte(`
		[
			[
				[
					{"initials": ["A", "A"], "grade": 1, "teacher_initial": "A", "student_number": 1}
				]
			]
		]
	`)

	ss := MakeStudentsFromBytes(students)

	expected := Student{
		Initials:       [2]rune{'A', 'A'},
		Grade:          1,
		TeacherInitial: 'A',
		StudentNumber:  1,
	}

	if !ss.Verify(expected) {
		t.Errorf("expected %v to be verified", expected)
	}

	expected = Student{
		Initials:       [2]rune{'A', 'A'},
		Grade:          1,
		TeacherInitial: 'A',
		StudentNumber:  2,
	}

	if ss.Verify(expected) {
		t.Errorf("expected %v to not be verified", expected)
	}
}

func TestStudent_UnmarshalJSON(t *testing.T) {
	student := []byte(`{"initials": ["A", "A"], "grade": 1, "teacher_initial": "A", "student_number": 1}`)

	s := &Student{}
	err := json.Unmarshal(student, s)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := Student{
		Initials:       [2]rune{'A', 'A'},
		Grade:          1,
		TeacherInitial: 'A',
		StudentNumber:  1,
	}

	if *s != expected {
		t.Errorf("expected %v, got %v", expected, *s)
	}

	student = []byte(`{"initials": ["A"], "grade": 1, "teacher_initial": "A", "student_number": 1}`)

	s = &Student{}
	err = json.Unmarshal(student, s)
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
