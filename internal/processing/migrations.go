package processing

type action int

const (
	actionNoop action = 1 << iota
	actionSaveRaw
	actionGenerateSamples
	actionExtractExif

	initialActions = actionSaveRaw | actionGenerateSamples | actionExtractExif
)

// allActions is a list of all possible actions
var allActions = []action{actionSaveRaw, actionGenerateSamples, actionExtractExif}

// migrations defines which actions to perform to migrate images
var migrations = []action{
	actionNoop,            // Used to be save raw, must have happened already
	actionNoop,            // Used to be extract photo type, now done as sampling
	actionGenerateSamples, // First time generate of samples
	actionExtractExif,     // First time extraction of EXIF
	actionExtractExif,     // Switch to using the Simpic fork of goexif with support for ORFs
	actionGenerateSamples, // Switch how image formats are stored in the DB, add more details
}
