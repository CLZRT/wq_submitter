package viewer

type Idea struct {
	ID                int64
	IdeaAlphaTemplate string
	IdeaTitle         string
	IdeaDesc          string
	StartIdx          int64
	EndIdx            int64
	NextIdx           int64
	SuccessNum        int64
	FailNum           int64
	ConcurrencyNum    int64
	IsFinished        int64
}
