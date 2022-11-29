package nav

func Resume(resumeInfo *NewResumerInfo) (*TraverseResult, error) {

	resumer, err := newResumer(resumeInfo)

	if err != nil {
		return nil, err
	}
	_ = resumer
	result := resumer.Continue()

	return result, result.Error
}
