package nav

func Resume(resumeInfo *NewResumerInfo) (*TraverseResult, error) {

	resumer, err := NewResumer(resumeInfo)

	if err != nil {
		return nil, err
	}
	_ = resumer
	// result := resumer.Continue()

	return &TraverseResult{}, nil
}
