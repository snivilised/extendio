package nav

import (
	"io/fs"
	"os"

	"github.com/samber/lo"
	xi18n "github.com/snivilised/extendio/i18n"
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
			return xi18n.NewPathNotFoundError("Traverse Item", params.path)
		},
		func() error {
			return xi18n.NewThirdPartyErr(params.err)
		},
	)

	callbackErr := params.agent.proxy(&agentProxyParams{
		item: &TraverseItem{
			Path:     params.path,
			Info:     params.info,
			Error:    err,
			Children: []fs.DirEntry{},
		},
		frame: params.frame,
	})

	return callbackErr
}
