package main

type FeatureType string

const (
	NONE                      FeatureType = "None"
	GOROUTINE                 FeatureType = "Goroutine"
	RECEIVE                   FeatureType = "Receive"
	SEND                      FeatureType = "Send"
	MAKE_CHAN                 FeatureType = "Synchronous chan"
	GO_IN_FOR                 FeatureType = "Go in for"
	RANGE_OVER_CHAN           FeatureType = "Range over chan"
	GO_IN_CONSTANT_FOR        FeatureType = "Goroutine in for with constant (constant)"
	KNOWN_CHAN_DEPTH          FeatureType = "Known chan length (length)"
	UNKNOWN_CHAN_DEPTH        FeatureType = "Unknown chan length"
	MAKE_CHAN_IN_FOR          FeatureType = "Make chan in for"
	MAKE_CHAN_IN_CONSTANT_FOR FeatureType = "Make chan in constant for"
	ARRAY_OF_CHANNELS         FeatureType = "Array of chans"
	CONSTANT_CHAN_ARRAY       FeatureType = "Constant array of chans (length)"
	CHAN_SLICE                FeatureType = "Slice array of chans"
	CHAN_MAP                  FeatureType = "Map of chans"
	CLOSE_CHAN                FeatureType = "Close chan"
	WAITGROUP                 FeatureType = "Waitgroup"
	KNOWN_ADD                 FeatureType = "Waitgroup Add(const)"
	UNKNOWN_ADD               FeatureType = "Waitgroup Add(var)"
	DONE                      FeatureType = "Waitgroup Done()"
	MUTEX                     FeatureType = "Mutex"
	UNLOCK                    FeatureType = "Mutex Unlock()"
	LOCK                      FeatureType = "Mutex Lock()"
	SELECT                    FeatureType = "Select (number of branch)"
	DEFAULT_SELECT            FeatureType = "Select with default (number of branch)"
	ASSIGN_CHAN_IN_FOR        FeatureType = "Assign chan in for"
	CHAN_OF_CHANS             FeatureType = "Channel of channels"
	RECEIVE_CHAN              FeatureType = "Receive only chan (<-chan)"
	SEND_CHAN                 FeatureType = "Send only chan (chan<-)"
	PARAM_CHAN                FeatureType = "chan used as a param"

	GOROUTINE_COUNT                 int = 1
	RECEIVE_COUNT                   int = 2
	SEND_COUNT                      int = 3
	MAKE_CHAN_COUNT                 int = 4
	GO_IN_FOR_COUNT                 int = 5
	RANGE_OVER_CHAN_COUNT           int = 6
	GO_IN_CONSTANT_FOR_COUNT        int = 7
	KNOWN_CHAN_DEPTH_COUNT          int = 8
	UNKNOWN_CHAN_DEPTH_COUNT        int = 9
	MAKE_CHAN_IN_FOR_COUNT          int = 10
	MAKE_CHAN_IN_CONSTANT_FOR_COUNT int = 23
	ARRAY_OF_CHANNELS_COUNT         int = 11
	CONSTANT_CHAN_ARRAY_COUNT       int = 12
	CHAN_SLICE_COUNT                int = 13
	CHAN_MAP_COUNT                  int = 14
	CLOSE_CHAN_COUNT                int = 15
	SELECT_COUNT                    int = 16
	DEFAULT_SELECT_COUNT            int = 17
	ASSIGN_CHAN_IN_FOR_COUNT        int = 18
	CHAN_OF_CHANS_COUNT             int = 19
	RECEIVE_CHAN_COUNT              int = 20
	SEND_CHAN_COUNT                 int = 21
	PARAM_CHAN_COUNT                int = 22
	WAITGROUP_COUNT                 int = 23
	KNOWN_ADD_COUNT                 int = 24
	UNKNOWN_ADD_COUNT               int = 25
	DONE_COUNT                      int = 26
	MUTEX_COUNT                     int = 27
	UNLOCK_COUNT                    int = 28
	LOCK_COUNT                      int = 29
)

type Feature struct {
	F_type         FeatureType
	F_type_num     int
	F_filename     string
	F_package_name string
	F_line_num     int
	F_number       string // A number used to report additional info about a feature
	F_commit       string // the commit of the project at the time the feature was found
	F_project_name string // the project name of the feature
}

// takes a list of feature and sets their feature number according to their types
func setFeaturesNumber(counter *Counter) {

	features := counter.Features
	counter.Features = []*Feature{}

	for _, feature := range features {
		switch feature.F_type {
		case GOROUTINE:
			feature.F_type_num = GOROUTINE_COUNT
		case RECEIVE:
			feature.F_type_num = RECEIVE_COUNT
		case SEND:
			feature.F_type_num = SEND_COUNT
		case MAKE_CHAN:
			feature.F_type_num = MAKE_CHAN_COUNT
		case GO_IN_FOR:
			feature.F_type_num = GO_IN_FOR_COUNT
		case RANGE_OVER_CHAN:
			feature.F_type_num = RANGE_OVER_CHAN_COUNT
		case GO_IN_CONSTANT_FOR:
			feature.F_type_num = GO_IN_CONSTANT_FOR_COUNT
		case KNOWN_CHAN_DEPTH:
			feature.F_type_num = KNOWN_CHAN_DEPTH_COUNT
		case UNKNOWN_CHAN_DEPTH:
			feature.F_type_num = UNKNOWN_CHAN_DEPTH_COUNT
		case MAKE_CHAN_IN_FOR:
			feature.F_type_num = MAKE_CHAN_IN_FOR_COUNT
		case MAKE_CHAN_IN_CONSTANT_FOR:
			feature.F_type_num = MAKE_CHAN_IN_CONSTANT_FOR_COUNT
		case ARRAY_OF_CHANNELS:
			feature.F_type_num = ARRAY_OF_CHANNELS_COUNT
		case CONSTANT_CHAN_ARRAY:
			feature.F_type_num = CONSTANT_CHAN_ARRAY_COUNT
		case CHAN_SLICE:
			feature.F_type_num = CHAN_SLICE_COUNT
		case CHAN_MAP:
			feature.F_type_num = CHAN_MAP_COUNT
		case CLOSE_CHAN:
			feature.F_type_num = CLOSE_CHAN_COUNT
		case SELECT:
			feature.F_type_num = SELECT_COUNT
		case DEFAULT_SELECT:
			feature.F_type_num = DEFAULT_SELECT_COUNT
		case ASSIGN_CHAN_IN_FOR:
			feature.F_type_num = ASSIGN_CHAN_IN_FOR_COUNT
		case CHAN_OF_CHANS:
			feature.F_type_num = CHAN_OF_CHANS_COUNT
		case SEND_CHAN:
			feature.F_type_num = SEND_CHAN_COUNT
		case RECEIVE_CHAN:
			feature.F_type_num = RECEIVE_CHAN_COUNT
		case PARAM_CHAN:
			feature.F_type_num = PARAM_CHAN_COUNT
		case WAITGROUP:
			feature.F_type_num = WAITGROUP_COUNT
		case KNOWN_ADD:
			feature.F_type_num = KNOWN_ADD_COUNT
		case UNKNOWN_ADD:
			feature.F_type_num = UNKNOWN_ADD_COUNT
		case MUTEX:
			feature.F_type_num = MUTEX_COUNT
		case LOCK:
			feature.F_type_num = LOCK_COUNT
		case UNLOCK:
			feature.F_type_num = UNLOCK_COUNT
		}

		counter.Features = append(counter.Features, feature)
	}
}
