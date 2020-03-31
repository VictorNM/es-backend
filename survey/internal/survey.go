package internal

import (
	"time"
)

type SurveyRow struct {
	ID         int
	CourseID   int
	TemplateID int
	Answered   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TemplateRow struct {
	ID int
}

type QuestionRow struct {
	ID      int
	Content string
}

type TemplateQuestionRow struct {
	TemplateID int
	QuestionID int
}

type SurveyQuestionRow struct {
	SurveyID   int
	QuestionID int
	Answer     int
}
