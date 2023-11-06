package net

type ErrorMessage struct {
	Text string
}

type ChallengeMessage struct {
	Algorithm  string
	Challenge  string
	Difficulty int
}

type SolutionMessage struct {
	Algorithm string
	Challenge string
	Solution  string
}

type GrantMessage struct {
	Text string
}
