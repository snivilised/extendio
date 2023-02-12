package nav

func Resume(resumeInfo *NewResumerInfo) (*TraverseResult, error) {
	// TODO: should only return a result with error embedded as member
	//
	resumer, err := resumerFactory{}.create(resumeInfo)

	if err != nil {
		return &TraverseResult{
			Error: err,
		}, err
	}
	result := resumer.Continue()

	return result, result.Error
}
