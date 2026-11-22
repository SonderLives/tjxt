package examlogic

// 题目类型
const (
	QuestionTypeSingle     = 1 // 单选题
	QuestionTypeMultiple   = 2 // 多选题
	QuestionTypeUncertain  = 3 // 不定向选择题
	QuestionTypeJudge      = 4 // 判断题
	QuestionTypeSubjective = 5 // 主观题
)

// 难易度
const (
	QuestionDifficultyEasy   = 1 // 简单
	QuestionDifficultyMedium = 2 // 中等
	QuestionDifficultyHard   = 3 // 困难
)
