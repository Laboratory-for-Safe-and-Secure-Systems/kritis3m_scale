package service

type HeartbeatInstruction int

const (
	HB_SHUTDOWN          HeartbeatInstruction = -1
	HB_RESTART           HeartbeatInstruction = -2
	HB_SLEEPMODE         HeartbeatInstruction = -3
	HB_NOCONFIGAVAILABLE HeartbeatInstruction = -4 //node can keep his old config or go to sleep?
	HB_NOTHING           HeartbeatInstruction = 0
	HB_REQUESTPOLICIES   HeartbeatInstruction = 1 // request new configuration policies/config from the distribution server
	HB_POSTSYSTEMSTATUS  HeartbeatInstruction = 2
	HB_SETDEBUGLEVEL     HeartbeatInstruction = 3
	HB_CHANGEHBINTERVAL  HeartbeatInstruction = 4
)

type LOGMODE int

const (
	SERVERLOG LOGMODE = 0
	STDOUT    LOGMODE = 1
	FILELOG   LOGMODE = 2
)

type HeartbeatResponse struct {
	HBInterval    uint64               `json:"hb_interval,omitempty"`
	DebugLevel    uint                 `json:"debug_level,omitempty"`
	LOGLEVEL      LOGMODE              `json:"log_mode,omitempty"`
	HBInstruction HeartbeatInstruction `json:"instruction,omitempty"`
}
