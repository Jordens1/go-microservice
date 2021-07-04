package grades

func init() {
	students = []Student{
		{
			ID:        1,
			FirstName: "dalang",
			LastName:  "wu",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 90,
				},
				{
					Title: "Quiz 2",
					Type:  GradeTest,
					Score: 100,
				},
			},
		},
		{
			ID:        2,
			FirstName: "jinlian",
			LastName:  "pan",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 67,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 56,
				},
				{
					Title: "Quiz 2",
					Type:  GradeTest,
					Score: 87,
				},
			},
		},
		{
			ID:        3,
			FirstName: "qin",
			LastName:  "ximen",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 44,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 66,
				},
				{
					Title: "Quiz 2",
					Type:  GradeTest,
					Score: 88,
				},
			},
		},
	}

}
