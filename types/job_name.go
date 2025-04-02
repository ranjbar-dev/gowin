package types

type JobName string

const (
	JobNamePing      JobName = "ping"
	JobCopyClipboard JobName = "copy-clipboard"
	JobLockScreen    JobName = "lock-screen"

	JobMoveMouse       JobName = "move-mouse"
	JobMouseLeftClick  JobName = "mouse-left-click"
	JobMouseRightClick JobName = "mouse-right-click"
)

func JobNameFromString(name string) JobName {

	return JobName(name)
}
