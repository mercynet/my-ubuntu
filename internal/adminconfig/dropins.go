// Package adminconfig defines the fixed administrative files allowed by the
// catalog. Detection and the privileged helper use the same expected contents.
package adminconfig

type Dropin struct {
	Path, Template string
	Apply          [][]string
}

var Dropins = map[string]Dropin{
	"inotify": {"/etc/sysctl.d/60-meu-ubuntu-inotify.conf",
		"# Gerido por Meu Ubuntu: watchers para IDEs/bundlers em projetos grandes\nfs.inotify.max_user_watches = %s\n",
		nil},
	"journald": {"/etc/systemd/journald.conf.d/60-meu-ubuntu.conf",
		"# Gerido por Meu Ubuntu: limite de retenção do journal\n[Journal]\nSystemMaxUse=%s\n",
		[][]string{{"systemctl", "restart", "systemd-journald"}}},
}
