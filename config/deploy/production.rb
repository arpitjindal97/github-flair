

  set :ssh_options, {
	  keys: %w(secrets/doorlock-root-ssh.key),
    forward_agent: false,
#    auth_methods: %w(password)
  }
