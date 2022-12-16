package nav

func Resume(resumeInfo *NewResumerInfo) (*TraverseResult, error) {

	resumer, err := (&resumerFactory{}).create(resumeInfo)

	if err != nil {
		return nil, err
	}
	_ = resumer
	result := resumer.Continue()

	return result, result.Error
}
