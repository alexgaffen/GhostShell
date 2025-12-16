package main

// A comprehensive fake filesystem map for a standard Ubuntu server
// Format: "Directory Path" -> ["List", "of", "Files"]
var fakeFS = map[string][]string{
	"/": {
		"bin", "boot", "dev", "etc", "home", "lib", "lib32", "lib64",
		"libx32", "media", "mnt", "opt", "proc", "root", "run", "sbin",
		"srv", "sys", "tmp", "usr", "var",
	},

	// --- BINARIES (Common commands so 'ls /bin' looks real) ---
	"/bin": {
		"bash", "cat", "chmod", "chown", "cp", "dash", "date", "dd", "df",
		"dir", "dmesg", "echo", "egrep", "false", "fgrep", "grep", "gunzip",
		"gzip", "hostname", "ip", "kill", "ln", "ls", "lsblk", "mkdir",
		"mknod", "mktemp", "more", "mount", "mv", "nano", "nc", "netstat",
		"ping", "ps", "pwd", "rm", "rmdir", "sed", "sh", "sleep", "ss",
		"stty", "su", "sync", "tar", "touch", "true", "umount", "uname",
		"vi", "zcat",
	},
	"/sbin": {
		"fdisk", "fsck", "getty", "halt", "ifconfig", "init", "iptables",
		"mkfs", "mkswap", "reboot", "route", "shutdown", "sysctl",
	},
	"/usr/bin": {
		"apt", "apt-get", "awk", "base64", "curl", "diff", "du", "find",
		"gcc", "git", "groups", "head", "htop", "id", "less", "logger",
		"make", "man", "mysql", "openssl", "perl", "php", "python3", "scp",
		"ssh", "ssh-keygen", "sudo", "tail", "tee", "top", "tr", "uptime",
		"vim", "wc", "wget", "whoami", "zip", "unzip",
	},

	// --- BOOT ---
	"/boot": {
		"config-5.4.0-150-generic", "grub", "initrd.img", "vmlinuz",
		"System.map-5.4.0-150-generic",
	},

	// --- CONFIGURATION (/etc) ---
	"/etc": {
		"alternatives", "apt", "bash.bashrc", "bind", "ca-certificates",
		"cron.d", "cron.daily", "cron.hourly", "crontab", "default",
		"environment", "fstab", "group", "group-", "hostname", "hosts",
		"init.d", "issue", "ld.so.conf", "localtime", "login.defs",
		"logrotate.d", "lsb-release", "machine-id", "modules", "motd",
		"mtab", "mysql", "netplan", "network", "nginx", "nsswitch.conf",
		"os-release", "pam.d", "passwd", "passwd-", "profile", "profile.d",
		"resolv.conf", "security", "services", "shadow", "shadow-", "skel",
		"ssh", "ssl", "sudoers", "sudoers.d", "sysctl.conf", "systemd",
		"timezone", "ucf.conf", "udev", "ufw", "update-motd.d", "vim",
	},
	"/etc/ssh": {
		"moduli", "ssh_config", "sshd_config", "ssh_host_ecdsa_key",
		"ssh_host_ecdsa_key.pub", "ssh_host_ed25519_key",
		"ssh_host_ed25519_key.pub", "ssh_host_rsa_key",
		"ssh_host_rsa_key.pub",
	},
	"/etc/apt": {
		"apt.conf.d", "auth.conf.d", "preferences.d", "sources.list",
		"sources.list.d", "trusted.gpg.d",
	},
	"/etc/network": {
		"if-down.d", "if-post-down.d", "if-pre-up.d", "if-up.d",
		"interfaces",
	},
	"/etc/nginx": {
		"conf.d", "fastcgi_params", "mime.types", "nginx.conf",
		"proxy_params", "scgi_params", "sites-available", "sites-enabled",
		"uwsgi_params",
	},

	// --- USERS (/home & /root) ---
	"/root": {
		".bash_history", ".bashrc", ".cache", ".config", ".local",
		".profile", ".ssh", ".viminfo", "snap",
	},
	"/root/.ssh": {
		"authorized_keys", "known_hosts",
	},
	"/home": {
		"ubuntu",
	},
	"/home/ubuntu": {
		".bash_history", ".bash_logout", ".bashrc", ".cache", ".config",
		".local", ".profile", ".ssh", ".sudo_as_admin_successful",
	},
	"/home/ubuntu/.ssh": {
		"authorized_keys", "id_rsa", "id_rsa.pub", "known_hosts",
	},

	// --- LIBRARIES (Simplified) ---
	"/lib":        {"systemd", "udev", "modules", "firmware"},
	"/lib/systemd": {"system", "user"},
	"/usr/lib":    {"python3.8", "openssh", "sudo", "apt"},

	// --- VARIABLE DATA (/var) ---
	"/var": {
		"backups", "cache", "crash", "lib", "local", "lock", "log", "mail",
		"opt", "run", "spool", "tmp", "www",
	},
	"/var/log": {
		"alternatives.log", "apt", "auth.log", "bootstrap.log", "btmp",
		"dist-upgrade", "dmesg", "dpkg.log", "faillog", "journal",
		"kern.log", "lastlog", "nginx", "syslog", "udev", "wtmp",
	},
	"/var/www": {
		"html",
	},
	"/var/www/html": {
		"index.nginx-debian.html", "index.php", "wp-config.php",
	},
	"/var/backups": {
		"apt.extended_states.0", "dpkg.status.0", "gshadow.bak",
		"shadow.bak",
	},

	// --- DEVICES & PROC (Empty but present) ---
	"/dev":  {"console", "core", "fd", "full", "null", "ptmx", "pts", "random", "stderr", "stdin", "stdout", "tty", "urandom", "zero", "sda", "sda1"},
	"/proc": {"1", "cpuinfo", "meminfo", "mounts", "net", "sys", "uptime", "version"},
	"/sys":  {"block", "bus", "class", "dev", "devices", "firmware", "fs", "kernel", "module", "power"},
	"/tmp":  {".X11-unix", ".ICE-unix", "systemd-private-xyz"},
	"/run":  {"lock", "pid", "sshd.pid", "systemd", "udev", "user"},
}