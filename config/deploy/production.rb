

  set :ssh_options, {
	  keys: %w(secrets/server-2-ssh.key),
    forward_agent: false,
#    auth_methods: %w(password)
  }
