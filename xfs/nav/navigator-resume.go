package nav

func ResumeLegacy(resumeInfo *ResumerInfo) (*TraverseResult, error) {
	// TODO: should only return a result with error embedded as member
	//
	resumer, err := resumerFactory{}.construct(resumeInfo)

	if err != nil {
		return &TraverseResult{
			Error: err,
		}, err
	}
	result := resumer.Continue()

	return result, result.Error
}

func Resume(info *ResumerInfo) (*TraverseResult, error) {
	return nil, nil
}

func Walk(path string) (*TraverseResult, error) {

	return nil, nil
}
