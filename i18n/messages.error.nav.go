package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// see https://dave.cheney.net/2014/12/24/inspecting-errors
//
// The requirement for localisation runs counter to that
// explained in this article, in particular the definition
// of typed errors increasing the api surface of a package
// and therefore makes the api more brittle. This issue is
// due-ly noted, but if translations are important, then we
// have to live with this problem unless another approach
// is available. Its not really recommended to provide foreign
// translations for external packages as this creates as
// undesirable coupling, but the option is there just in case.
// To ameliorate api surface area issue, limit error definitions
// to those errors that are intended to be displayed to
// the end user. Internal errors that can be handled, should not
// have translations templates defined for them as the user
// won't see them.
//
// As is presented in the article, clients are better off
// asserting errors for behaviour, not type, but this aspect
// should not be at cross purposes with the requirement for
// localisation.
//
//  In summary then, for ...
//
// * package authors: provide predicate interface definitions
// for errors that can be handled, eg "Timeout() bool". Also,
// use errors.Wrap to add context to another error.
// * package users: don't check an error's type, query for the
// declared interface, and invoke the provided predicates
// to identify an actual error.
//
// An alternative to providing foreign translations is just
// to handle the 3rd party error and Wrapping it up with a
// local error in the desired language. Sure, the inner error
// will be defined in the library's default language, but that
// can be wrapped (errors.Wrap), providing content in the
// required but library un-supported language.
//

type ExtendioTemplData struct{}

func (td ExtendioTemplData) SourceId() string {
	return EXTENDIO_SOURCE_ID
}

// ❌ FailedToReadDirectoryContents

// FailedToReadDirectoryContentsTemplData failed to resume using file
type FailedToReadDirectoryContentsTemplData struct {
	ExtendioTemplData
	Path   string
	Reason error
}

func (td FailedToReadDirectoryContentsTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-read-directory-contents.extendio.nav",
		Description: "Failed to read directory contents from the path specified",
		Other:       "failed to read directory contents '{{.Path}}' (reason: {{.REASON}})",
	}
}

// FailedToReadDirectoryContentsErrorBehaviourQuery used to query if an error is:
// "Failed to read directory contents from the path specified"
type FailedToReadDirectoryContentsErrorBehaviourQuery interface {
	FailedToReadDirectoryContents() bool
}

type FailedToReadDirectoryContentsError struct {
	LocalisableError
}

// FailedToReadDirectoryContents enables the client to check if error is FailedToReadDirectoryContentsError
// via QueryFailedToReadDirectoryContentsError
func (e FailedToReadDirectoryContentsError) FailedToReadDirectoryContents() bool {
	return true
}

// NewFailedToReadDirectoryContentsError creates a FailedToReadDirectoryContentsError
func NewFailedToReadDirectoryContentsError(path string, reason error) FailedToReadDirectoryContentsError {
	return FailedToReadDirectoryContentsError{
		LocalisableError: LocalisableError{
			Data: FailedToReadDirectoryContentsTemplData{
				Path:   path,
				Reason: reason,
			},
		},
	}
}

// QueryFailedToReadDirectoryContentsError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryFailedToReadDirectoryContentsError(target error) bool {
	return QueryGeneric[FailedToReadDirectoryContentsErrorBehaviourQuery]("FailedToReadDirectoryContents", target)
}

// ❌ FailedToResumeFromFile

// FailedToResumeFromFileTemplData failed to resume using file
type FailedToResumeFromFileTemplData struct {
	ExtendioTemplData
	Path   string
	Reason error
}

func (td FailedToResumeFromFileTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-resume-from-file.extendio.nav",
		Description: "Failed to resume traverse operation from the resume file specified",
		Other:       "failed to resume from file '{{.Path}}' (reason: {{.REASON}})",
	}
}

// FailedToResumeFromFileErrorBehaviourQuery used to query if an error is:
// "Failed to resume traverse operation from the resume file specified"
type FailedToResumeFromFileErrorBehaviourQuery interface {
	FailedToResumeFromFile() bool
}

type FailedToResumeFromFileError struct {
	LocalisableError
}

// FailedToResumeFromFile enables the client to check if error is FailedToResumeFromFileError
// via QueryFailedToResumeFromFileError
func (e FailedToResumeFromFileError) FailedToResumeFromFile() bool {
	return true
}

// NewFailedToResumeFromFileError creates a FailedToResumeFromFileError
func NewFailedToResumeFromFileError(path string, reason error) FailedToResumeFromFileError {
	return FailedToResumeFromFileError{
		LocalisableError: LocalisableError{
			Data: FailedToResumeFromFileTemplData{
				Path:   path,
				Reason: reason,
			},
		},
	}
}

// QueryFailedToResumeFromFileError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryFailedToResumeFromFileError(target error) bool {
	return QueryGeneric[FailedToResumeFromFileErrorBehaviourQuery]("FailedToResumeFromFile", target)
}

// ❌ InvalidConfigEntry

// InvalidConfigEntryTemplData failed to resume using file
type InvalidConfigEntryTemplData struct {
	ExtendioTemplData
	Value string
	At    string
}

func (td InvalidConfigEntryTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-config.entry.extendio.nav",
		Description: "Invalid entry specified in config at the location specified",
		Other:       "invalid entry '{{.Value}}' specified in config at {{.At}}",
	}
}

// InvalidConfigEntryErrorBehaviourQuery used to query if an error is:
// "Failed to resume traverse operation from the resume file specified"
type InvalidConfigEntryErrorBehaviourQuery interface {
	InvalidConfigEntry() bool
}

type InvalidConfigEntryError struct {
	LocalisableError
}

// InvalidConfigEntry enables the client to check if error is InvalidConfigEntryError
// via QueryInvalidConfigEntryError
func (e InvalidConfigEntryError) InvalidConfigEntry() bool {
	return true
}

// NewInvalidConfigEntryError creates a InvalidConfigEntryError
func NewInvalidConfigEntryError(value, at string) InvalidConfigEntryError {
	return InvalidConfigEntryError{
		LocalisableError: LocalisableError{
			Data: InvalidConfigEntryTemplData{
				Value: value,
				At:    at,
			},
		},
	}
}

// QueryInvalidConfigEntryError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryInvalidConfigEntryError(target error) bool {
	return QueryGeneric[InvalidConfigEntryErrorBehaviourQuery]("InvalidConfigEntry", target)
}

// ❌ InvalidResumeStrategy

// InvalidResumeStrategyTemplData failed to resume using file
type InvalidResumeStrategyTemplData struct {
	ExtendioTemplData
	Value string
}

func (td InvalidResumeStrategyTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-resume-strategy.internal.extendio.nav",
		Description: "Invalid resume strategy specified",
		Other:       "invalid resume strategy '{{.Value}}' specified",
	}
}

// InvalidResumeStrategyErrorBehaviourQuery used to query if an error is:
// "Failed to resume traverse operation from the resume file specified"
type InvalidResumeStrategyErrorBehaviourQuery interface {
	InvalidResumeStrategy() bool
}

type InvalidResumeStrategyError struct {
	LocalisableError
}

// InvalidResumeStrategy enables the client to check if error is InvalidResumeStrategyError
// via QueryInvalidResumeStrategyError
func (e InvalidResumeStrategyError) InvalidResumeStrategy() bool {
	return true
}

// NewInvalidResumeStrategyError creates a InvalidResumeStrategyError
func NewInvalidResumeStrategyError(value string) InvalidResumeStrategyError {
	return InvalidResumeStrategyError{
		LocalisableError: LocalisableError{
			Data: InvalidResumeStrategyTemplData{
				Value: value,
			},
		},
	}
}

// QueryInvalidResumeStrategyError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryInvalidResumeStrategyError(target error) bool {
	return QueryGeneric[InvalidResumeStrategyErrorBehaviourQuery]("InvalidResumeStrategy", target)
}

// ❌ MissingCallback

// missing callback (internal)
type MissingCallbackTemplData struct {
	ExtendioTemplData
}

func (td MissingCallbackTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "missing-callback.internal.extendio",
		Description: "Missing callback (internal error)",
		Other:       "missing callback (internal error)",
	}
}

// MissingCallbackBehaviourQuery used to query if an error is:
// "Missing callback (internal error)"
type MissingCallbackBehaviourQuery interface {
	MissingCallback() bool
}

// MissingCallbackError, this is a coding error where client has not provided
// a callback required by the api.
type MissingCallbackError struct {
	LocalisableError
}

// MissingCallback enables the client to check if error is MissingCallbackError
// via QueryMissingCallbackError
func (e MissingCallbackError) MissingCallback() bool {
	return true
}

// NewMissingCallbackError creates a MissingCallbackError
func NewMissingCallbackError() MissingCallbackError {
	return MissingCallbackError{
		LocalisableError: LocalisableError{
			Data: MissingCallbackTemplData{},
		},
	}
}

// QueryMissingCallbackError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryMissingCallbackError(target error) bool {
	return QueryGeneric[MissingCallbackBehaviourQuery]("MissingCallback", target)
}

// ❌ MissingCustomFilterDefinition

// Missing custom filter definition (config)
type MissingCustomFilterDefinitionTemplData struct {
	ExtendioTemplData
	At string
}

func (td MissingCustomFilterDefinitionTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "missing-custom-filter-definition.config.extendio",
		Description: "Missing custom filter definition (config error)",
		Other:       "missing custom filter definition at {{.At}} (config error)",
	}
}

// MissingCustomFilterDefinitionBehaviourQuery used to query if an error is:
// "Missing callback (internal error)"
type MissingCustomFilterDefinitionBehaviourQuery interface {
	MissingCustomFilterDefinition() bool
}

// MissingCustomFilterDefinitionError, this is a config error where client has not provided
// the definition of a custom filter having set the filter type to custom and a pattern
type MissingCustomFilterDefinitionError struct {
	LocalisableError
}

// MissingCustomFilterDefinition enables the client to check if error is
// MissingCustomFilterDefinitionError via QueryMissingCustomFilterDefinitionError
func (e MissingCustomFilterDefinitionError) MissingCustomFilterDefinition() bool {
	return true
}

// NewMissingCustomFilterDefinitionError creates a MissingCustomFilterDefinitionError
func NewMissingCustomFilterDefinitionError(at string) MissingCustomFilterDefinitionError {
	return MissingCustomFilterDefinitionError{
		LocalisableError: LocalisableError{
			Data: MissingCustomFilterDefinitionTemplData{
				At: at,
			},
		},
	}
}

// QueryMissingCustomFilterDefinitionError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryMissingCustomFilterDefinitionError(target error) bool {
	return QueryGeneric[MissingCustomFilterDefinitionBehaviourQuery]("MissingCustomFilterDefinition", target)
}

// ❌ NotADirectory

// NotADirectoryTemplData path is not a directory
type NotADirectoryTemplData struct {
	ExtendioTemplData
	Path string
}

func (td NotADirectoryTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "not-a-directory.extendio.nav",
		Description: "File system path is not a directory",
		Other:       "file system path '{{.Path}}', is not a directory",
	}
}

// NotADirectoryErrorBehaviourQuery used to query if an error is:
// "File system path is not a directory"
type NotADirectoryErrorBehaviourQuery interface {
	NotADirectory() bool
}

type NotADirectoryError struct {
	LocalisableError
}

// NotADirectory enables the client to check if error is NotADirectoryError
// via QueryNotADirectoryError
func (e NotADirectoryError) NotADirectory() bool {
	return true
}

// NewNotADirectoryError creates a NotADirectoryError
func NewNotADirectoryError(path string) NotADirectoryError {
	return NotADirectoryError{
		LocalisableError: LocalisableError{
			Data: NotADirectoryTemplData{
				Path: path,
			},
		},
	}
}

// QueryNotADirectoryError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryNotADirectoryError(target error) bool {
	return QueryGeneric[NotADirectoryErrorBehaviourQuery]("NotADirectory", target)
}

// ❌ SortFnFailed

// sort function failed (internal)
type SortFnFailedTemplData struct {
	ExtendioTemplData
}

func (td SortFnFailedTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "sort-fn-failed.internal.extendio.nav",
		Description: "Sort function failed (internal error)",
		Other:       "sort function failed (internal error)",
	}
}

// SortFnFailedBehaviourQuery used to query if an error is:
// "Sort function failed (internal error)"
type SortFnFailedBehaviourQuery interface {
	SortFnFailed() bool
}

type SortFnFailedError struct {
	LocalisableError
}

// SortFnFailed enables the client to check if error is SortFnFailedError
// via QuerySortFnFailedError
func (e SortFnFailedError) SortFnFailed() bool {
	return true
}

// NewSortFnFailedError creates a SortFnFailedError
func NewSortFnFailedError() SortFnFailedError {
	return SortFnFailedError{
		LocalisableError: LocalisableError{
			Data: SortFnFailedTemplData{},
		},
	}
}

// QuerySortFnFailedError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QuerySortFnFailedError(target error) bool {
	return QueryGeneric[SortFnFailedBehaviourQuery]("SortFnFailed", target)
}

// ❌ TerminateTraverse

// terminate traverse
type TerminateTraverseTemplData struct {
	ExtendioTemplData
	Reason string
}

func (td TerminateTraverseTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "terminate-traverse.extendio.nav",
		Description: "Traversal terminated",
		Other:       "terminate traversal: '{{.Reason}}'",
	}
}

// TerminateTraverseBehaviourQuery used to query if an error is:
// "Traversal terminated"
type TerminateTraverseBehaviourQuery interface {
	TraverseTerminated() bool
}

// TerminateTraverseError indicates that traversal has been terminated early
type TerminateTraverseError struct {
	LocalisableError
}

// TerminateTraverse enables the client to check if error is SortFnFailedError
// via QueryTerminateTraverseError
func (e SortFnFailedError) TerminateTraverse() bool {
	return true
}

// NewTerminateTraverseError creates a TerminateTraverseError
func NewTerminateTraverseError() TerminateTraverseError {
	return TerminateTraverseError{
		LocalisableError: LocalisableError{
			Data: SortFnFailedTemplData{},
		},
	}
}

// QueryTerminateTraverseError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryTerminateTraverseError(target error) bool {
	return QueryGeneric[SortFnFailedBehaviourQuery]("TerminateTraverse", target)
}

// ❌ ThirdPartyError

// ====================================================================

// This file is automatically processed to create the active.en.json
// by running:
// $ goi18n extract -format json

// ThirdPartyErrorTemplData third party un-translated error
type ThirdPartyErrorTemplData struct {
	ExtendioTemplData

	Error error
}

func (td ThirdPartyErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "third-party-error.extendio",
		Description: "These errors are generated by dependencies that don't support localisation",
		Other:       "third party error: '{{.Error}}'",
	}
}

// ThirdPartyError represents an error received by a dependency that does
// not support i18n.
type ThirdPartyError struct {
	LocalisableError
}

// NewThirdPartyErr creates a ThirdPartyErr
func NewThirdPartyErr(err error) ThirdPartyError {

	return ThirdPartyError{
		LocalisableError: LocalisableError{
			Data: ThirdPartyErrorTemplData{
				Error: err,
			},
		},
	}
}

// ❌ UnknownMarshalFormat

// UnknownMarshalFormatTemplData unknown marshall format specified in config by user
type UnknownMarshalFormatTemplData struct {
	ExtendioTemplData
	Format string
	At     string
}

func (td UnknownMarshalFormatTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "unknown-marshal-format.config.extendio.nav",
		Description: "Unknown marshal format specified in config",
		Other:       "unknown marshal format {{.Format}} specified at {{.At}}",
	}
}

// UnknownMarshalFormatErrorBehaviourQuery used to query if an error is:
// "Unknown marshal format specified in config"
type UnknownMarshalFormatErrorBehaviourQuery interface {
	UnknownMarshalFormat() bool
}

type UnknownMarshalFormatError struct {
	LocalisableError
}

// UnknownMarshalFormat enables the client to check if error is UnknownMarshalFormatError
// via QueryUnknownMarshalFormatError
func (e UnknownMarshalFormatError) UnknownMarshalFormat() bool {
	return true
}

// NewUnknownMarshalFormatError creates a UnknownMarshalFormatError
func NewUnknownMarshalFormatError(format string, at string) UnknownMarshalFormatError {
	return UnknownMarshalFormatError{
		LocalisableError: LocalisableError{
			Data: UnknownMarshalFormatTemplData{
				Format: format,
				At:     at,
			},
		},
	}
}

// QueryUnknownMarshalFormatError helper function to enable identification of
// an error via its behaviour, rather than by its type.
func QueryUnknownMarshalFormatError(target error) bool {
	return QueryGeneric[UnknownMarshalFormatErrorBehaviourQuery]("UnknownMarshalFormat", target)
}
