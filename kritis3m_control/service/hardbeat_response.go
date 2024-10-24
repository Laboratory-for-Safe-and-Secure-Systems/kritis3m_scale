package service

type HardbeatInstruction int

const (
	HB_SHUTDOWN          HardbeatInstruction = -1
	HB_RESTART           HardbeatInstruction = -2
	HB_SLEEPMODE         HardbeatInstruction = -3
	HB_NOCONFIGAVAILABLE HardbeatInstruction = -4 //node can keep his old config or go to sleep?
	HB_NOTHING           HardbeatInstruction = 0
	HB_REQUESTPOLICIES   HardbeatInstruction = 1 // request new configuration policies/config from the distribution server
	HB_POSTSYSTEMSTATUS  HardbeatInstruction = 2
	HB_SETDEBUGLEVEL     HardbeatInstruction = 3
	HB_CHANGEHBINTERVAL  HardbeatInstruction = 4
)

type LOGMODE int

const (
	SERVERLOG LOGMODE = 0
	STDOUT    LOGMODE = 1
	FILELOG   LOGMODE = 2
)

type HardbeatResponse struct {
	HBInterval    uint64              `json:"hb_interval,omitempty"`
	DebugLevel    uint                `json:"debug_level,omitempty"`
	LOGLEVEL      LOGMODE             `json:"log_mode,omitempty"`
	HBInstruction HardbeatInstruction `json:"instruction,omitempty"`
}
