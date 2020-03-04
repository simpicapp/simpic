package processing

type action int

const (
	actionNoop action = 1 << iota
	actionSaveRaw
	actionGenerateSamples
	actionExtractExif
	actionCalculateTimestamp
)

// initialActions is the set of actions to perform for newly uploaded images
const initialActions = actionSaveRaw | actionGenerateSamples | actionExtractExif | actionCalculateTimestamp

// allActions is a list of all possible actions, in the order they will be executed if multiple are required
var allActions = []action{actionSaveRaw, actionGenerateSamples, actionExtractExif, actionCalculateTimestamp}

// migrations defines which actions to perform to migrate images that were uploaded using previous versions of
// Simpic. Whenever Simpic starts, any migrations added to the list since it last processed photos are performed
// (after de-duping, and in the order specified in allActions).
var migrations = []action{
	actionNoop,               // Used to be save raw, must have happened already
	actionNoop,               // Used to be extract photo type, now done as sampling
	actionGenerateSamples,    // Newly added
	actionExtractExif,        // Newly added
	actionExtractExif,        // Switch to using the Simpic fork of goexif with support for ORFs
	actionGenerateSamples,    // Switch how image formats are stored in the DB, add more details
	actionCalculateTimestamp, // Newly added
	actionCalculateTimestamp, // More file name formats added
}
