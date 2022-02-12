package student

import (
	"encoding/json"
	"errors"
	"os"
	"unicode"
)

type Students [26][26][]Student

type Student struct {
	Initials       [2]rune `json:"initials"`
	Grade          int     `json:"grade"`
	TeacherInitial rune    `json:"teacher_initial"`
	StudentNumber  int     `json:"student_number"`
}

var (
	ErrInvalidInitials      = errors.New("invalid initials")
	ErrInvalidInitialLength = errors.New("empty initial")
)

func MakeStudents(filename string) Students {
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return MakeStudentsFromBytes(contents)
}

func MakeStudentsFromBytes(b []byte) Students {
	ss := [26][26][]Student{}

	err := json.Unmarshal(b, &ss)
	if err != nil {
		panic(err)
	}

	return ss
}

func (s Students) Verify(student Student) bool {
	if unicode.IsLetter(student.Initials[0]) == false || unicode.IsLetter(student.Initials[1]) == false || unicode.IsLetter(student.TeacherInitial) == false {
		return false
	}

	if student.Grade < 1 || student.Grade > 12 {
		return false
	}

	for _, s := range s[student.Initials[0]-'A'][student.Initials[1]-'A'] {
		if s == student {
			return true
		}
	}

	return false
}

func (s *Student) UnmarshalJSON(b []byte) error {
	type Alias Student
	aux := &struct {
		*Alias
		TeacherInitial string    `json:"teacher_initial"`
		Initials       [2]string `json:"initials"`
	}{
		Alias: (*Alias)(s),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	if len(aux.Initials) != 2 {
		return ErrInvalidInitials
	}

	if len(aux.Initials[0]) < 1 || len(aux.Initials[1]) < 1 || len(aux.TeacherInitial) < 1 {
		return ErrInvalidInitialLength
	}

	if unicode.IsLetter(rune(aux.Initials[0][0])) == false || unicode.IsLetter(rune(aux.Initials[1][0])) == false || unicode.IsLetter(rune(aux.TeacherInitial[0])) == false {
		return ErrInvalidInitials
	}

	s.Initials = [2]rune{
		unicode.ToUpper(rune(aux.Initials[0][0])),
		unicode.ToUpper(rune(aux.Initials[1][0])),
	}

	s.TeacherInitial = unicode.ToUpper(rune(aux.TeacherInitial[0]))

	return nil
}
