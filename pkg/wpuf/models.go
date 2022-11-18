package wpuf

type LoginAttempt struct {
	Success  bool
	Password string
	Error    bool
}

type Settings struct {
	Url       string
	Wordlist  string
	Username  string
	Proxy     string
	Timeout   int
	Threads   int
	MaxErrors int
	Enumerate bool
	Error     bool
}
