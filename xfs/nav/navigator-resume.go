package nav

func ResumeLegacy(resumeInfo *ResumerInfo) (*TraverseResult, error) {
	resumer, err := resumerFactory{}.new(resumeInfo)

	if err != nil {
		return &TraverseResult{
			err: err,
		}, err
	}

	return resumer.Continue()
}

func Resume(info *ResumerInfo) (*TraverseResult, error) {
	return nil, nil
}

func Walk(path string) (*TraverseResult, error) {

	return nil, nil
}
