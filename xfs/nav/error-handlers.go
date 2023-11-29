package nav

import (
	"io/fs"
	"os"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/i18n"
)

type fileSystemErrorParams struct {
	err   error
	path  string
	info  fs.FileInfo
	agent *navigationAgent
	frame *navigationFrame
}

type fileSystemErrorHandler interface {
	accept(params *fileSystemErrorParams) error
}

type notifyCallbackErrorHandler struct {
}

func (h *notifyCallbackErrorHandler) accept(params *fileSystemErrorParams) error {
	err := lo.TernaryF(os.IsNotExist(params.err),
		func() error {
			return i18n.NewPathNotFoundError("Traverse Item", params.path)
		},
		func() error {
			return i18n.NewThirdPartyErr(params.err)
		},
	)

	callbackErr := params.frame.proxy(&TraverseItem{
		Path:     params.path,
		Info:     params.info,
		Error:    err,
		Children: []fs.DirEntry{},
	}, nil)

	return callbackErr
}
